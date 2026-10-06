package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"math"
	"net/http"
	"time"

	"readsrv/internal/analysis"
	"readsrv/internal/store"
)

const maxRange = 7 * 24 * time.Hour

type Config struct {
	OfflineThreshold time.Duration
	TripGap          time.Duration
}

type Server struct {
	st  *store.Store
	cfg Config
	now func() time.Time
}

func New(st *store.Store, cfg Config) *Server {
	return &Server{st: st, cfg: cfg, now: time.Now}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", s.healthz)
	mux.HandleFunc("GET /api/devices", s.listDevices)
	mux.HandleFunc("GET /api/devices/{device_id}", s.getDevice)
	mux.HandleFunc("GET /api/devices/{device_id}/track", s.track)
	mux.HandleFunc("GET /api/devices/{device_id}/trips", s.trips)
	mux.HandleFunc("GET /api/reports/daily", s.dailyReport)
	return mux
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("encode: %v", err)
	}
}

func writeErr(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": msg})
}

func (s *Server) internal(w http.ResponseWriter, err error) {
	log.Printf("internal error: %v", err)
	writeErr(w, http.StatusInternalServerError, "internal error")
}

func ts(t time.Time) string { return t.UTC().Format(time.RFC3339) }

func round(v float64, n int) float64 {
	p := math.Pow10(n)
	return math.Round(v*p) / p
}

func (s *Server) parseRange(r *http.Request) (from, to time.Time, err error) {
	q := r.URL.Query()
	to = s.now().UTC()
	if v := q.Get("to"); v != "" {
		if to, err = time.Parse(time.RFC3339, v); err != nil {
			return from, to, errors.New("to must be RFC3339")
		}
	}
	from = to.Add(-24 * time.Hour)
	if v := q.Get("from"); v != "" {
		if from, err = time.Parse(time.RFC3339, v); err != nil {
			return from, to, errors.New("from must be RFC3339")
		}
	}
	if !from.Before(to) {
		return from, to, errors.New("from must be before to")
	}
	if to.Sub(from) > maxRange {
		return from, to, errors.New("range must be 7 days or less")
	}
	return from.UTC(), to.UTC(), nil
}

func (s *Server) healthz(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	if err := s.st.Ping(ctx); err != nil {
		writeErr(w, http.StatusServiceUnavailable, "db unreachable")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

type position struct {
	Lat        float64 `json:"lat"`
	Lon        float64 `json:"lon"`
	RecordedAt string  `json:"recorded_at"`
}

type deviceJSON struct {
	DeviceID     string   `json:"device_id"`
	Status       string   `json:"status"`
	LastPosition position `json:"last_position"`
}

func (s *Server) deviceView(d store.Device) deviceJSON {
	status := "offline"
	if s.now().Sub(d.Latest.RecordedAt) <= s.cfg.OfflineThreshold {
		status = "online"
	}
	return deviceJSON{
		DeviceID: d.ID, Status: status,
		LastPosition: position{Lat: d.Latest.Lat, Lon: d.Latest.Lon, RecordedAt: ts(d.Latest.RecordedAt)},
	}
}

func (s *Server) listDevices(w http.ResponseWriter, r *http.Request) {
	ds, err := s.st.ListDevices(r.Context())
	if err != nil {
		s.internal(w, err)
		return
	}
	out := make([]deviceJSON, 0, len(ds))
	for _, d := range ds {
		out = append(out, s.deviceView(d))
	}
	writeJSON(w, http.StatusOK, map[string]any{"devices": out})
}

func (s *Server) deviceOr404(w http.ResponseWriter, r *http.Request) *store.Device {
	d, err := s.st.GetDevice(r.Context(), r.PathValue("device_id"))
	if errors.Is(err, store.ErrNotFound) {
		writeErr(w, http.StatusNotFound, "device not found")
		return nil
	}
	if err != nil {
		s.internal(w, err)
		return nil
	}
	return d
}

func (s *Server) getDevice(w http.ResponseWriter, r *http.Request) {
	if d := s.deviceOr404(w, r); d != nil {
		writeJSON(w, http.StatusOK, s.deviceView(*d))
	}
}

func (s *Server) track(w http.ResponseWriter, r *http.Request) {
	v := s.deviceOr404(w, r)
	if v == nil {
		return
	}
	from, to, err := s.parseRange(r)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	pts, err := s.st.Track(r.Context(), v.ID, from, to)
	if err != nil {
		s.internal(w, err)
		return
	}
	type pt struct {
		Time string  `json:"time"`
		Lat  float64 `json:"lat"`
		Lon  float64 `json:"lon"`
	}
	out := make([]pt, 0, len(pts))
	for _, p := range pts {
		out = append(out, pt{ts(p.T), p.Lat, p.Lon})
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"device_id": v.ID, "from": ts(from), "to": ts(to), "count": len(out), "points": out,
	})
}

type coord struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type tripJSON struct {
	StartTime   string  `json:"start_time"`
	EndTime     string  `json:"end_time"`
	Start       coord   `json:"start"`
	End         coord   `json:"end"`
	DistanceKm  float64 `json:"distance_km"`
	DurationMin float64 `json:"duration_min"`
	AvgSpeedKmh float64 `json:"avg_speed_kmh"`
	MaxSpeedKmh float64 `json:"max_speed_kmh"`
	PointCount  int     `json:"point_count"`
}

func tripView(t analysis.Trip) tripJSON {
	return tripJSON{
		StartTime: ts(t.Start), EndTime: ts(t.End),
		Start: coord{t.StartLat, t.StartLon}, End: coord{t.EndLat, t.EndLon},
		DistanceKm: round(t.DistanceKm, 2), DurationMin: round(t.DurationMin, 1),
		AvgSpeedKmh: round(t.AvgKmh, 1), MaxSpeedKmh: round(t.MaxKmh, 1), PointCount: t.Points,
	}
}

func (s *Server) trips(w http.ResponseWriter, r *http.Request) {
	v := s.deviceOr404(w, r)
	if v == nil {
		return
	}
	from, to, err := s.parseRange(r)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	pts, err := s.st.Track(r.Context(), v.ID, from, to)
	if err != nil {
		s.internal(w, err)
		return
	}
	out := []tripJSON{}
	for _, t := range analysis.SegmentTrips(pts, s.cfg.TripGap) {
		out = append(out, tripView(t))
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"device_id": v.ID, "from": ts(from), "to": ts(to), "count": len(out), "trips": out,
	})
}

type dailyRow struct {
	DeviceID       string  `json:"device_id"`
	DistanceKm     float64 `json:"distance_km"`
	MovingMinutes  float64 `json:"moving_minutes"`
	TripCount      int     `json:"trip_count"`
	FirstDeparture *string `json:"first_departure"`
	LastArrival    *string `json:"last_arrival"`
}

func (s *Server) dailyReport(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	tz := q.Get("tz")
	if tz == "" {
		tz = "Asia/Tokyo"
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "invalid tz")
		return
	}
	var day time.Time
	if d := q.Get("date"); d != "" {
		if day, err = time.ParseInLocation("2006-01-02", d, loc); err != nil {
			writeErr(w, http.StatusBadRequest, "date must be YYYY-MM-DD")
			return
		}
	} else {
		n := s.now().In(loc)
		day = time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, loc)
	}
	from := day
	to := time.Date(day.Year(), day.Month(), day.Day()+1, 0, 0, 0, 0, loc)

	ids, err := s.st.DevicesInRange(r.Context(), from, to)
	if err != nil {
		s.internal(w, err)
		return
	}
	rows := make([]dailyRow, 0, len(ids))
	for _, id := range ids {
		pts, err := s.st.Track(r.Context(), id, from, to)
		if err != nil {
			s.internal(w, err)
			return
		}
		row := dailyRow{DeviceID: id}
		trips := analysis.SegmentTrips(pts, s.cfg.TripGap)
		for i, t := range trips {
			row.DistanceKm += t.DistanceKm
			row.MovingMinutes += t.DurationMin
			row.TripCount++
			if i == 0 {
				f := ts(t.Start)
				row.FirstDeparture = &f
			}
			l := ts(t.End)
			row.LastArrival = &l
		}
		row.DistanceKm = round(row.DistanceKm, 2)
		row.MovingMinutes = round(row.MovingMinutes, 1)
		rows = append(rows, row)
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"date": day.Format("2006-01-02"), "tz": tz, "devices": rows,
	})
}
