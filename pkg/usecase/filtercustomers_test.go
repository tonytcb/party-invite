package usecase

import (
	"bytes"
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"git.codesubmit.io/sfox/golang-party-invite-ivsjhn/pkg/domain"
	"git.codesubmit.io/sfox/golang-party-invite-ivsjhn/pkg/infrastructure/logger"
)

func TestFilterCustomers_ByNearLocation(t *testing.T) {
	t.Parallel()

	saoPaulo, err := domain.NewCoordinate("-23.533773", "-46.625290")
	if err != nil {
		t.Fatal("failed to build coordinate")
	}
	montreal, err := domain.NewCoordinate("45.508888", "-73.561668")
	if err != nil {
		t.Fatal("failed to build coordinate")
	}
	paris, err := domain.NewCoordinate("2.2945", "48.8584")
	if err != nil {
		t.Fatal("failed to build coordinate")
	}
	rioDeJaneiro, err := domain.NewCoordinate("-22.908333", "-43.196388")
	if err != nil {
		t.Fatal("failed to build coordinate")
	}
	curitiba, err := domain.NewCoordinate("-25.441105", "-49.276855")
	if err != nil {
		t.Fatal("failed to build coordinate")
	}

	var (
		customer1 = domain.NewCustomer(1, "User name 1", domain.DublinLocation)
		customer2 = domain.NewCustomer(2, "User name 2", saoPaulo)
		customer3 = domain.NewCustomer(3, "User name 3", montreal)
		customer4 = domain.NewCustomer(4, "User name 4", paris)
		customer5 = domain.NewCustomer(5, "User name 5", rioDeJaneiro)
		customer6 = domain.NewCustomer(6, "User name 6", rioDeJaneiro)
		customer7 = domain.NewCustomer(7, "User name 7", curitiba)
	)

	var log = logger.NewLogger(&bytes.Buffer{})

	// log = logger.NewLogger(os.Stderr)

	type args struct {
		ctx                context.Context
		customers          domain.Customers
		baseLocation       *domain.Coordinate
		nearDistanceFilter decimal.Decimal
		orderBy            domain.OrderBy
	}
	tests := []struct {
		name    string
		args    args
		want    domain.Customers
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "should return an empty customers list when input is empty",
			args: args{
				ctx:                context.Background(),
				customers:          []domain.Customer{},
				baseLocation:       domain.DublinLocation,
				nearDistanceFilter: decimal.NewFromInt32(100),
				orderBy:            domain.OrderByCustomerID,
			},
			want:    []domain.Customer{},
			wantErr: assert.NoError,
		},
		{
			name: "should return a list with one customer near to Dublin",
			args: args{
				ctx:                context.Background(),
				customers:          []domain.Customer{customer1, customer2, customer3, customer4},
				baseLocation:       domain.DublinLocation,
				nearDistanceFilter: decimal.NewFromInt32(100),
				orderBy:            domain.OrderByCustomerID,
			},
			want:    []domain.Customer{customer1},
			wantErr: assert.NoError,
		},
		{
			name: "should return a list with all brazilian customers near to 500km from Sao Paulo",
			args: args{
				ctx:                context.Background(),
				customers:          []domain.Customer{customer7, customer6, customer1, customer3, customer5},
				baseLocation:       saoPaulo,
				nearDistanceFilter: decimal.NewFromInt32(500),
				orderBy:            domain.OrderByCustomerID,
			},
			want:    []domain.Customer{customer5, customer6, customer7},
			wantErr: assert.NoError,
		},
		{
			name: "should return filter a list with 5.000 customers",
			args: args{
				ctx:                context.Background(),
				customers:          generateCustomersList([]domain.Customer{customer7, customer6, customer1, customer3, customer5}, 1_000),
				baseLocation:       saoPaulo,
				nearDistanceFilter: decimal.NewFromInt32(500),
				orderBy:            domain.OrderByCustomerID,
			},
			want:    []domain.Customer{customer5, customer6, customer7},
			wantErr: assert.NoError,
		},
		{
			name: "should error on orderBy parameter",
			args: args{
				ctx:                context.Background(),
				customers:          []domain.Customer{customer1},
				baseLocation:       domain.DublinLocation,
				nearDistanceFilter: decimal.NewFromInt32(1000),
				orderBy:            55,
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "unexpected order by")
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			f := NewFilterCustomers(log)

			got, err := f.ByNearLocation(
				tt.args.ctx,
				tt.args.customers,
				tt.args.baseLocation,
				tt.args.nearDistanceFilter,
				tt.args.orderBy,
			)

			tt.wantErr(t, err)

			assert.EqualValues(t, tt.want, got)
		})
	}
}

func generateCustomersList(baseList domain.Customers, N int) domain.Customers {
	var result = append([]domain.Customer{}, baseList...)

	for i := 0; i < N; i++ {
		for _, customer := range baseList {
			location := customer.Location
			newLocation, _ := domain.NewCoordinate(
				location.Latitude.StringFixed(distancePrecision),
				location.Longitude.Add(decimal.NewFromInt32(int32(randNumber(100, 200)))).StringFixed(distancePrecision),
			)
			result = append(result, customer.WithLocation(newLocation))
		}
	}

	return result
}

func randNumber(min, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return min + r.Intn(max-min+1)
}
