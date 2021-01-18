package datatypes

import "testing"

func TestPolygon_GeomFromText(t *testing.T) {
	tests := []struct {
		name string
		p    Polygon
		want string
	}{
		{
			name: "1",
			p:    []Point{{"0","0"},{"0","3"},{"3","3"},{"3","0"}},
			want: "GeomFromText('Polygon((0 0,0 3,3 3,3 0))')",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.GeomFromText(); got != tt.want {
				t.Errorf("GeomFromText() = %v, want %v", got, tt.want)
			}
		})
	}
}
