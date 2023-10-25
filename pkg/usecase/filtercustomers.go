package usecase

import (
	"context"
	"errors"
	"sort"

	"github.com/shopspring/decimal"

	"git.codesubmit.io/sfox/golang-party-invite-ivsjhn/pkg/domain"
	"git.codesubmit.io/sfox/golang-party-invite-ivsjhn/pkg/infrastructure/logger"
)

type FilterCustomers struct {
	log logger.Logger
}

func NewFilterCustomers(log logger.Logger) *FilterCustomers {
	return &FilterCustomers{
		log: log,
	}
}

func (f *FilterCustomers) ByNearLocation(
	ctx context.Context,
	customers domain.Customers,
	baseLocation *domain.Coordinate,
	nearDistanceFilter decimal.Decimal,
	orderBy domain.OrderBy,
) (domain.Customers, error) {
	var (
		log               = f.log.FromContext(ctx)
		nearCustomersByID = make(map[int]domain.Customer) // using a map to remove duplicated customers
	)

	// loop to calculate the distance for every customer and filter them
	for _, customer := range customers {
		select {
		case <-ctx.Done():
			return nil, errors.New("context done while calculating customers' distances")
		default:
		}

		difference := baseLocation.Difference(customer.Location)

		const precision = 3
		log.Infof("Distance calculation, customer-id=%d distance=%s", customer.ID, difference.StringFixed(precision))

		if difference.GreaterThan(nearDistanceFilter) {
			continue // filter out customer
		}

		nearCustomersByID[customer.ID] = customer
	}

	var result = customersMapValues(nearCustomersByID)

	switch orderBy {
	case domain.OrderByCustomerID:
		sort.Slice(result, func(i, j int) bool {
			return result[i].ID < result[j].ID
		})

	default:
		return nil, errors.New("unexpected order by")
	}

	return result, nil
}

func customersMapValues(customersMap map[int]domain.Customer) domain.Customers {
	var r = make([]domain.Customer, 0)

	for _, distance := range customersMap {
		r = append(r, distance)
	}

	return r
}
