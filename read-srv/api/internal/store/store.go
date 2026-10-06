package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"readsrv/internal/analysis"
)

var ErrNotFound = errors.New("not found")

type Latest struct {
	Lat, Lon   float64
	RecordedAt time.Time
}

type Device struct {
	ID     string
	Latest Latest
}

type Store struct {
	db    *sql.DB
	label string
}

func New(db *sql.DB, label string) *Store { return &Store{db: db, label: label} }

func (s *Store) Ping(ctx context.Context) error { return s.db.PingContext(ctx) }

// lat / lon が数値のレコードのみ
const (
	cols = `recorded_at, CAST(JSON_EXTRACT(payload, '$.lat') AS DECIMAL(9, 6)), CAST(JSON_EXTRACT(payload, '$.lon') AS DECIMAL(9, 6))`
	ok   = `JSON_TYPE(JSON_EXTRACT(payload, '$.lat')) IN ('INTEGER', 'DOUBLE', 'DECIMAL')
	  AND JSON_TYPE(JSON_EXTRACT(payload, '$.lon')) IN ('INTEGER', 'DOUBLE', 'DECIMAL')`
)

func (s *Store) latest(ctx context.Context, id string) (*Latest, error) {
	var l Latest
	err := s.db.QueryRowContext(ctx,
		`SELECT `+cols+` FROM records
		 WHERE device_id = ? AND label = ? AND `+ok+`
		 ORDER BY recorded_at DESC LIMIT 1`, id, s.label).Scan(&l.RecordedAt, &l.Lat, &l.Lon)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func (s *Store) deviceIDs(ctx context.Context, from, to *time.Time) ([]string, error) {
	q := `SELECT DISTINCT device_id FROM records WHERE label = ?`
	args := []any{s.label}
	if from != nil {
		q += ` AND recorded_at >= ? AND recorded_at < ?`
		args = append(args, from.UTC(), to.UTC())
	}
	rows, err := s.db.QueryContext(ctx, q+` ORDER BY device_id`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func (s *Store) ListDevices(ctx context.Context) ([]Device, error) {
	ids, err := s.deviceIDs(ctx, nil, nil)
	if err != nil {
		return nil, err
	}
	out := []Device{}
	for _, id := range ids {
		l, err := s.latest(ctx, id)
		if err != nil {
			return nil, err
		}
		if l != nil {
			out = append(out, Device{ID: id, Latest: *l})
		}
	}
	return out, nil
}

func (s *Store) GetDevice(ctx context.Context, id string) (*Device, error) {
	l, err := s.latest(ctx, id)
	if err != nil {
		return nil, err
	}
	if l == nil {
		return nil, ErrNotFound
	}
	return &Device{ID: id, Latest: *l}, nil
}

// 指定期間にデータのあるデバイス
func (s *Store) DevicesInRange(ctx context.Context, from, to time.Time) ([]string, error) {
	return s.deviceIDs(ctx, &from, &to)
}

// 期間は [from, to)
func (s *Store) Track(ctx context.Context, id string, from, to time.Time) ([]analysis.Point, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT `+cols+` FROM records
		 WHERE device_id = ? AND label = ? AND recorded_at >= ? AND recorded_at < ? AND `+ok+`
		 ORDER BY recorded_at`, id, s.label, from.UTC(), to.UTC())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var pts []analysis.Point
	for rows.Next() {
		var p analysis.Point
		if err := rows.Scan(&p.T, &p.Lat, &p.Lon); err != nil {
			return nil, err
		}
		pts = append(pts, p)
	}
	return pts, rows.Err()
}
