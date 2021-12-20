// Command ncscr - Starry Night screensaver from Norton Commander from 1980s.
// Serves no real purpose other than nostalgic feelings
package main

import "testing"

func Test_maxStars(t *testing.T) {
	type args struct {
		maxX int
		maxY int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"80x25", args{80, 25}, 100},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maxStars(tt.args.maxX, tt.args.maxY); got != tt.want {
				t.Errorf("maxStars() = %v, want %v", got, tt.want)
			}
		})
	}
}
