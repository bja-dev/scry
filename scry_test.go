package scry

import (
	"testing"
	"slices"
)

func TestScanFile(t *testing.T) {
	tests := []struct {
		path string
		want []string
	}{
		{
			path: "user-list.txt",
			want: []string{"locking_latu", "maltesemario", "dfgrdds"},
		},
		{
			path: "test-scanfile.txt",
			want: []string{"this is a test","hello world","123","^ that was empty", "but it shouldn't make it to the list", "!!!!"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			lines, _ := scanFile(tt.path)
			if slices.Compare(tt.want, lines) != 0 {
				t.Errorf("Expected %v - Got %v", tt.want, lines)
			}
		})
	}
}
