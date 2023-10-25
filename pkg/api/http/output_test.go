package http

import (
	"context"
	"git.codesubmit.io/sfox/golang-party-invite-ivsjhn/pkg/domain"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

func Test_customersToJSONOutput(t *testing.T) {
	t.Parallel()

	type args struct {
		input domain.Customers
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "should return an empty json when customers is empty",
			args: args{
				input: []domain.Customer{},
			},
			want:    []byte(`[]`),
			wantErr: assert.NoError,
		},
		{
			name: "should return a valid json containing 2 customers",
			args: args{
				input: []domain.Customer{
					{
						ID:       100,
						Name:     "Tony Tester",
						Location: domain.DublinLocation,
					},
					{
						ID:       200,
						Name:     "Jon Doe",
						Location: domain.DublinLocation,
					},
				},
			},
			want:    []byte(`[{"id":100,"name":"Tony Tester"},{"id":200,"name":"Jon Doe"}]`),
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			got, err := customersToJSONOutput(tt.args.input)

			tt.wantErr(t, err)
			if err != nil {
				return
			}

			assert.EqualValues(t, string(tt.want), string(got))
		})
	}
}

func Test_errToHttpErrorCode(t *testing.T) {
	t.Parallel()

	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "should return StatusUnprocessableEntity http status code",
			args: args{
				err: domain.NewErrInvalidArgument("root error", "some details"),
			},
			want: http.StatusUnprocessableEntity,
		},
		{
			name: "should return StatusGatewayTimeout http status code",
			args: args{
				err: context.Canceled,
			},
			want: http.StatusGatewayTimeout,
		},
		{
			name: "should return StatusInternalServerError http status code for unknown errors",
			args: args{
				err: io.ErrNoProgress,
			},
			want: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			got := errToStatusCode(tt.args.err)

			assert.Equalf(t, tt.want, got, "errToHttpErrorCode(%v)", tt.args.err)
		})
	}
}
