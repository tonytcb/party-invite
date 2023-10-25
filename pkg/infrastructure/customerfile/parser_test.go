package customerfile

import (
	"context"
	"git.codesubmit.io/sfox/golang-party-invite-ivsjhn/pkg/domain"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestCustomersFileParser_Parse(t *testing.T) {
	t.Parallel()

	canceledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	type args struct {
		ctx         context.Context
		fileContent string
	}
	tests := []struct {
		name    string
		args    args
		want    domain.Customers
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "should parse a valid file with 3 customers",
			args: args{
				ctx: context.Background(),
				fileContent: `{"latitude": "54.1225", "user_id": 27, "name": "Enid Gallagher", "longitude": "-8.143333"}

{"latitude": "53.1229599", "user_id": 6, "name": "Theresa Enright", "longitude": "-6.2705202"}
{"latitude": "52.2559432", "user_id": 9, "name": "Jack Dempsey", "longitude": "-7.1048927"}`,
			},
			want: []domain.Customer{
				{
					ID:   27,
					Name: "Enid Gallagher",
					Location: &domain.Coordinate{
						Latitude:  decimal.RequireFromString("54.1225"),
						Longitude: decimal.RequireFromString("-8.143333"),
					},
				},
				{
					ID:   6,
					Name: "Theresa Enright",
					Location: &domain.Coordinate{
						Latitude:  decimal.RequireFromString("53.1229599"),
						Longitude: decimal.RequireFromString("-6.2705202"),
					},
				},
				{
					ID:   9,
					Name: "Jack Dempsey",
					Location: &domain.Coordinate{
						Latitude:  decimal.RequireFromString("52.2559432"),
						Longitude: decimal.RequireFromString("-7.1048927"),
					},
				},
			},
			wantErr: assert.NoError,
		},

		{
			name: "should error on invalid latitude on 2nd line",
			args: args{
				ctx: context.Background(),
				fileContent: `{"latitude": "54.1225", "user_id": 27, "name": "Enid Gallagher", "longitude": "-8.143333"}
{"latitude": "invalid number", "user_id": 9, "name": "Jack Dempsey", "longitude": "-7.1048927"}`,
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "error to parse customers' location")
			},
		},
		{
			name: "should error on invalid json on 1st line",
			args: args{
				ctx:         context.Background(),
				fileContent: `invalid json line`,
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "error to parse line=1")
			},
		},
		{
			name: "should error on context done",
			args: args{
				ctx:         canceledCtx,
				fileContent: `{"latitude": "52.2559432", "user_id": 9, "name": "Jack Dempsey", "longitude": "-7.1048927"}`,
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "context done while parsing file")
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.args.fileContent)

			c := CustomersFileParser{}

			got, err := c.Parse(tt.args.ctx, reader)

			tt.wantErr(t, err)
			if err != nil {
				return
			}

			assert.EqualValues(t, tt.want, got)
		})
	}
}
