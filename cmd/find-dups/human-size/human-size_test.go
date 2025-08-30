package humansize

import "testing"

func TestSizeToString(t *testing.T) {
	type args struct {
		size int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"bytes", args{512}, "512 b"},
		{"kilobytes", args{21512}, "22 Kb"},
		{"megabytes", args{14_221_512}, "14 Mb"},
		{"megabytes", args{14_235_221_512}, "14.2 Gb"},
		{"megabytes", args{1_404_235_221_512}, "1.4 Tb"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SizeToString(tt.args.size); got != tt.want {
				t.Errorf("SizeToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
