package domain

import (
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
)

// To validate the outputs was used the https://latlongdata.com/distance-calculator/ website.
func TestCoordinate_Difference(t *testing.T) {
	t.Parallel()

	const decimalPlaces = 3

	var (
		winnipeg = &Coordinate{
			Latitude:  decimal.RequireFromString("49.895077"),
			Longitude: decimal.RequireFromString("-97.138451"),
		}
		regina = &Coordinate{
			Latitude:  decimal.RequireFromString("50.445210"),
			Longitude: decimal.RequireFromString("-104.618896"),
		}
		saoPaulo = &Coordinate{
			Latitude:  decimal.RequireFromString("-23.533773"),
			Longitude: decimal.RequireFromString("-46.625290"),
		}
		paris = &Coordinate{
			Latitude:  decimal.RequireFromString("2.2945"),
			Longitude: decimal.RequireFromString("48.8584"),
		}
	)

	type fields struct {
		firstLocation *Coordinate
	}
	type args struct {
		secondLocation *Coordinate
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "should calculate the distance between Winnipeg and Regina",
			fields: fields{
				firstLocation: winnipeg,
			},
			args: args{
				secondLocation: regina,
			},
			want: "536.036",
		},
		{
			name: "should calculate the distance between Sao Paulo and Paris",
			fields: fields{
				firstLocation: saoPaulo,
			},
			args: args{
				secondLocation: paris,
			},
			want: "10668.315",
		},
		{
			name: "should calculate the distance between Sao Paulo and Winnipeg",
			fields: fields{
				firstLocation: saoPaulo,
			},
			args: args{
				secondLocation: winnipeg,
			},
			want: "9560.150",
		},
		{
			name: "should calculate the distance between Dublin and Sao Paulo",
			fields: fields{
				firstLocation: DublinLocation,
			},
			args: args{
				secondLocation: saoPaulo,
			},
			want: "9390.052",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			c := tt.fields.firstLocation

			if got := c.Difference(tt.args.secondLocation); got.StringFixed(decimalPlaces) != tt.want {
				t.Errorf("Difference() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCoordinate(t *testing.T) {
	t.Parallel()

	type args struct {
		latitude  string
		longitude string
	}
	tests := []struct {
		name    string
		args    args
		want    *Coordinate
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "should build a valid coordinate struct",
			args: args{
				latitude:  "-23.533773",
				longitude: "-46.625290",
			},
			want: &Coordinate{
				Latitude:  decimal.RequireFromString("-23.533773"),
				Longitude: decimal.RequireFromString("-46.625290"),
			},
			wantErr: assert.NoError,
		},
		{
			name: "should error on invalid latitude value",
			args: args{
				latitude:  "x",
				longitude: "-46.625290",
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "invalid latitude") && errors.As(err, &ErrInvalidArgument{})
			},
		},
		{
			name: "should error on invalid longitude value",
			args: args{
				latitude:  "-46.625290",
				longitude: "y",
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "invalid longitude") && errors.Is(err, &ErrInvalidArgument{})
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCoordinate(tt.args.latitude, tt.args.longitude)

			tt.wantErr(t, err)

			assert.EqualValues(t, tt.want, got)
		})
	}
}
