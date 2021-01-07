package track

import (
	"fmt"
	"reflect"
	"testing"
)

func p(x [][]TrackPointer) {
	for _, arr := range x {
		for _, v := range arr {
			fmt.Printf("%d ", v.Timestamp())
		}
		fmt.Println()
	}
}

func TestSegment__segmentByInterval(t *testing.T) {
	tests := []struct {
		name string
		args []TrackPointer
		want int
	}{

		{
			args: []TrackPointer{
				NewTrackPoint("1", 1, 131, 31),
				NewTrackPoint("1", 2, 131, 31),
				NewTrackPoint("1", 3, 131, 31),
				NewTrackPoint("1", 11, 131, 31),
				NewTrackPoint("1", 13, 131, 31),
				NewTrackPoint("1", 15, 131, 31),
				NewTrackPoint("1", 21, 131, 31),
				NewTrackPoint("1", 25, 131, 31),
				NewTrackPoint("1", 30, 131, 31),
			},
			want: 3,
		},

		//
		//{
		//	args: []hub.TrackPointer{
		//		hub.NewTrackPoint("1", 1, 131, 31),
		//		hub.NewTrackPoint("1", 2, 131, 31),
		//		hub.NewTrackPoint("1", 3, 131, 31),
		//		hub.NewTrackPoint("1", 4, 131, 31),
		//		hub.NewTrackPoint("1", 5, 131, 31),
		//		hub.NewTrackPoint("1", 11, 131, 31),
		//		hub.NewTrackPoint("1", 12, 131, 31),
		//		hub.NewTrackPoint("1", 13, 131, 31),
		//		hub.NewTrackPoint("1", 14, 131, 31),
		//	},
		//},
		//
		//
		//{
		//	args: []hub.TrackPointer{
		//		hub.NewTrackPoint("1", 1, 131, 31),
		//		hub.NewTrackPoint("1", 2, 131, 31),
		//		hub.NewTrackPoint("1", 3, 131, 31),
		//		hub.NewTrackPoint("1", 4, 131, 31),
		//		hub.NewTrackPoint("1", 5, 131, 31),
		//		hub.NewTrackPoint("1", 6, 131, 31),
		//		hub.NewTrackPoint("1", 7, 131, 31),
		//		hub.NewTrackPoint("1", 8, 131, 31),
		//		hub.NewTrackPoint("1", 9, 131, 31),
		//	},
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Segment{}
			got := d.SegmentByInterval(nil, tt.args...)
			p(got)
			if !reflect.DeepEqual(len(got), tt.want) {
				t.Errorf("_segmentByInterval() = %v, want %v", len(got), tt.want)
			}
		})
	}
}
