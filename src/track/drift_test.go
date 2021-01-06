package track

import (
	"github.com/linger1216/go-gis/model/hub"
	"reflect"
	"testing"
)

func TestDrift__DriftByStandardDeviation(t *testing.T) {
	tests := []struct {
		name string
		args []hub.TrackPointer
		want int
	}{

		{
			args: []hub.TrackPointer{
				hub.NewTrackPoint("1", 1, 0, 0),
				hub.NewTrackPoint("1", 1, 1, 1),
				hub.NewTrackPoint("1", 1, 2, 2),
				hub.NewTrackPoint("1", 1, 10, 2),
				hub.NewTrackPoint("1", 1, 4, 4),
				hub.NewTrackPoint("1", 1, 4, 5),
				hub.NewTrackPoint("1", 1, 5, 6),
				hub.NewTrackPoint("1", 1, 7, 7),
				hub.NewTrackPoint("1", 1, 6, 8),
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Drift{}
			got := d._DriftByStandardDeviation(nil, tt.args...)
			if !reflect.DeepEqual(len(got), tt.want) {
				t.Errorf("_segmentByInterval() = %v, want %v", len(got), tt.want)
			}
		})
	}
}
