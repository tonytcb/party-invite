package usecase

import (
	"context"
	"errors"
	"sort"
	"sync"

	"github.com/shopspring/decimal"

	"github.com/tonytcb/party-invite/pkg/domain"
	"github.com/tonytcb/party-invite/pkg/infrastructure/logger"
)

const distancePrecision = 3

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
		customersCh       = make(chan domain.Customer)
	)

	log.Infof("Count customers=%d", len(customers))

	wg := &sync.WaitGroup{}

	// loop to calculate the distance for every customer and filter them
	for _, c := range customers {
		wg.Add(1)

		go func(customer domain.Customer) {
			defer wg.Done()

			difference := baseLocation.Difference(customer.Location)

			log.Infof("Distance calculation, customer-id=%d distance=%s", customer.ID, difference.StringFixed(distancePrecision))

			if difference.GreaterThan(nearDistanceFilter) {
				return // filter out customer
			}

			customersCh <- customer
		}(c)
	}

	go func() {
		wg.Wait()
		close(customersCh)
	}()

	for customer := range customersCh {
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
