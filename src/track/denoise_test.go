package track

import (
	"github.com/linger1216/go-gis/model/hub"
	"reflect"
	"testing"
)

func TestDenoise__part(t *testing.T) {
	d := NewDenoise()
	tests := []struct {
		name string
		args []hub.TrackPointer
		want int
	}{
		{
			args: []hub.TrackPointer{
				hub.NewTrackPoint("1", 1, 131, 31),
				hub.NewTrackPoint("1", 2, 131, 31),
				hub.NewTrackPoint("1", 3, 131, 31),
				hub.NewTrackPoint("1", 4, 131, 31),
				hub.NewTrackPoint("1", 5, 131, 31),
				hub.NewTrackPoint("1", 6, 131, 31),
				hub.NewTrackPoint("1", 8, 131, 31),
				hub.NewTrackPoint("1", 10, 131, 31),
				hub.NewTrackPoint("1", 11, 131, 31),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := d._part(tt.args...); !reflect.DeepEqual(len(got), tt.want) {
				t.Errorf("_part() = %v, want %v", len(got), tt.want)
			}
		})
	}
}
