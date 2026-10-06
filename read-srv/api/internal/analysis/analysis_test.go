package analysis

import (
	"math"
	"testing"
	"time"
)

var t0 = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

func TestHaversine(t *testing.T) {
	tests := []struct {
		name                   string
		lat1, lon1, lat2, lon2 float64
		want, tol              float64
	}{
		{"同一点", 35.0, 139.0, 35.0, 139.0, 0, 1e-9},
		{"緯度1度", 0, 0, 1, 0, 111.19, 0.05},
		{"東京駅-新宿駅", 35.681236, 139.767125, 35.690921, 139.700258, 6.1, 0.2},
		{"赤道で経度1度", 0, 0, 0, 1, 111.19, 0.05},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Haversine(tt.lat1, tt.lon1, tt.lat2, tt.lon2)
			if math.Abs(got-tt.want) > tt.tol {
				t.Errorf("got %f want %f±%f", got, tt.want, tt.tol)
			}
		})
	}
}

func TestSegment(t *testing.T) {
	a := Point{t0, 35, 139}
	tests := []struct {
		name   string
		b      Point
		wantOK bool
	}{
		{"通常走行 約36km/h", Point{t0.Add(10 * time.Second), 35.001, 139}, true},
		{"時間差ゼロ", Point{t0, 35.001, 139}, false},
		{"時間が逆転", Point{t0.Add(-time.Second), 35.001, 139}, false},
		{"ノイズ 200km/h 超", Point{t0.Add(10 * time.Second), 35.1, 139}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, ok := Segment(a, tt.b)
			if ok != tt.wantOK {
				t.Errorf("ok=%v want %v", ok, tt.wantOK)
			}
		})
	}
}

func line(start time.Duration, lat0 float64, n int, stepDeg float64, interval time.Duration) []Point {
	var p []Point
	for i := 0; i < n; i++ {
		p = append(p, Point{t0.Add(start + time.Duration(i)*interval), lat0 + float64(i)*stepDeg, 139})
	}
	return p
}

func TestSegmentTrips(t *testing.T) {
	gap := 10 * time.Minute
	moving := line(0, 35.0, 31, 0.001, 20*time.Second)
	tests := []struct {
		name      string
		points    func() []Point
		wantTrips int
		check     func(t *testing.T, trips []Trip)
	}{
		{"空", func() []Point { return nil }, 0, nil},
		{"1点のみ", func() []Point { return line(0, 35, 1, 0, time.Second) }, 0, nil},
		{"単一トリップ", func() []Point { return moving }, 1, func(t *testing.T, tr []Trip) {
			if math.Abs(tr[0].DistanceKm-3.336) > 0.05 {
				t.Errorf("distance %f", tr[0].DistanceKm)
			}
			if math.Abs(tr[0].DurationMin-10) > 1e-9 {
				t.Errorf("duration %f", tr[0].DurationMin)
			}
			if math.Abs(tr[0].AvgKmh-20.0) > 0.5 || math.Abs(tr[0].MaxKmh-20.0) > 0.5 {
				t.Errorf("speed avg=%f max=%f", tr[0].AvgKmh, tr[0].MaxKmh)
			}
		}},
		{"入力順が逆でも同じ", func() []Point {
			r := append([]Point(nil), moving...)
			for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
				r[i], r[j] = r[j], r[i]
			}
			return r
		}, 1, nil},
		{"時間差が gap 超で分割", func() []Point {
			return append(line(0, 35, 10, 0.001, 20*time.Second), line(time.Hour, 35.1, 10, 0.001, 20*time.Second)...)
		}, 2, nil},
		{"時間差がちょうど gap なら分割しない", func() []Point {
			a := line(0, 35, 5, 0.001, 20*time.Second)
			b := line(80*time.Second+gap, 35.005, 5, 0.001, 20*time.Second)
			return append(a, b...)
		}, 1, nil},
		{"短い停車は分割しない", func() []Point {
			a := line(0, 35, 10, 0.001, 20*time.Second)
			stop := line(200*time.Second, 35.009, 18, 0, 20*time.Second)
			b := line(560*time.Second, 35.009, 10, 0.001, 20*time.Second)
			return append(append(a, stop...), b...)
		}, 1, nil},
		{"gap 以上の停車で分割", func() []Point {
			a := line(0, 35, 10, 0.001, 20*time.Second)
			stop := line(200*time.Second, 35.009, 61, 0, 20*time.Second)
			b := line(1420*time.Second, 35.009, 10, 0.001, 20*time.Second)
			return append(append(a, stop...), b...)
		}, 2, func(t *testing.T, tr []Trip) {
			if tr[0].EndLat != 35.009 || tr[1].StartLat != 35.009 {
				t.Errorf("境界座標: %f %f", tr[0].EndLat, tr[1].StartLat)
			}
			if !tr[1].Start.After(tr[0].End) {
				t.Errorf("後続トリップの開始は先行の終了より後")
			}
		}},
		{"ノイズ点は距離・速度から除外", func() []Point {
			p := line(0, 35, 11, 0.001, 20*time.Second)
			p[5].Lat += 1.0
			return p
		}, 1, func(t *testing.T, tr []Trip) {
			if tr[0].MaxKmh > 25 {
				t.Errorf("max %f", tr[0].MaxKmh)
			}
			if tr[0].DistanceKm > 3.0 {
				t.Errorf("distance %f", tr[0].DistanceKm)
			}
		}},
		{"停車のみは 0 トリップ", func() []Point { return line(0, 35, 61, 0, 20*time.Second) }, 0, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trips := SegmentTrips(tt.points(), gap)
			if len(trips) != tt.wantTrips {
				t.Fatalf("trips=%d want %d: %+v", len(trips), tt.wantTrips, trips)
			}
			if tt.check != nil {
				tt.check(t, trips)
			}
		})
	}
}
