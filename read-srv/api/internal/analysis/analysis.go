package analysis

import (
	"math"
	"sort"
	"time"
)

const (
	earthRadiusKm = 6371.0088
	StopRadiusKm  = 0.05
	MaxSpeedKmh   = 200.0
)

type Point struct {
	T   time.Time
	Lat float64
	Lon float64
}

type Trip struct {
	Start       time.Time
	End         time.Time
	StartLat    float64
	StartLon    float64
	EndLat      float64
	EndLon      float64
	DistanceKm  float64
	DurationMin float64
	AvgKmh      float64
	MaxKmh      float64
	Points      int
}

// Haversine 2点間距離(km)
func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
	rad := math.Pi / 180
	dLat := (lat2 - lat1) * rad
	dLon := (lon2 - lon1) * rad
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*rad)*math.Cos(lat2*rad)*math.Sin(dLon/2)*math.Sin(dLon/2)
	return 2 * earthRadiusKm * math.Asin(math.Min(1, math.Sqrt(a)))
}

// Segment 区間距離と速度
func Segment(a, b Point) (km, kmh float64, ok bool) {
	dt := b.T.Sub(a.T).Hours()
	if dt <= 0 {
		return 0, 0, false
	}
	km = Haversine(a.Lat, a.Lon, b.Lat, b.Lon)
	kmh = km / dt
	if kmh > MaxSpeedKmh {
		return 0, 0, false
	}
	return km, kmh, true
}

// SegmentTrips トリップ分割
func SegmentTrips(points []Point, gap time.Duration) []Trip {
	if len(points) < 2 {
		return nil
	}
	pts := append([]Point(nil), points...)
	sort.SliceStable(pts, func(i, j int) bool { return pts[i].T.Before(pts[j].T) })

	var segs [][]Point
	start := 0
	for i := 1; i < len(pts); i++ {
		if pts[i].T.Sub(pts[i-1].T) > gap {
			segs = append(segs, pts[start:i])
			start = i
		}
	}
	segs = append(segs, pts[start:])

	var trips []Trip
	for _, s := range segs {
		for _, part := range splitAtStops(s, gap) {
			if len(part) >= 2 {
				trips = append(trips, buildTrip(part))
			}
		}
	}
	return trips
}

func splitAtStops(s []Point, gap time.Duration) [][]Point {
	var out [][]Point
	start, i := 0, 0
	for i < len(s) {
		j := i
		for j+1 < len(s) && Haversine(s[i].Lat, s[i].Lon, s[j+1].Lat, s[j+1].Lon) < StopRadiusKm {
			j++
		}
		if j > i && s[j].T.Sub(s[i].T) >= gap {
			out = append(out, s[start:i+1])
			start = j
			i = j
			continue
		}
		i++
	}
	return append(out, s[start:])
}

func buildTrip(p []Point) Trip {
	t := Trip{
		Start: p[0].T, End: p[len(p)-1].T,
		StartLat: p[0].Lat, StartLon: p[0].Lon,
		EndLat: p[len(p)-1].Lat, EndLon: p[len(p)-1].Lon,
		Points: len(p),
	}
	for i := 1; i < len(p); i++ {
		km, kmh, ok := Segment(p[i-1], p[i])
		if !ok {
			continue
		}
		t.DistanceKm += km
		if kmh > t.MaxKmh {
			t.MaxKmh = kmh
		}
	}
	t.DurationMin = t.End.Sub(t.Start).Minutes()
	if h := t.End.Sub(t.Start).Hours(); h > 0 {
		t.AvgKmh = t.DistanceKm / h
	}
	return t
}
