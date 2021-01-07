package track

import (
	"reflect"
	"testing"
)

func TestDrift__DriftByStandardDeviation(t *testing.T) {
	tests := []struct {
		name string
		args []TrackPointer
		want int
	}{

		{
			args: []TrackPointer{
				NewTrackPoint("1", 1, 0, 0),
				NewTrackPoint("1", 1, 1, 1),
				NewTrackPoint("1", 1, 2, 2),
				NewTrackPoint("1", 1, 10, 2),
				NewTrackPoint("1", 1, 4, 4),
				NewTrackPoint("1", 1, 4, 5),
				NewTrackPoint("1", 1, 5, 6),
				NewTrackPoint("1", 1, 7, 7),
				NewTrackPoint("1", 1, 6, 8),
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
