package usecase

import (
	"context"
	"sort"
	"sync"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/tonytcb/party-invite/pkg/domain"
	"github.com/tonytcb/party-invite/pkg/infrastructure/logger"
)

const distancePrecision = 3

type FilterCustomersNotifier interface {
	Notify(context.Context, *domain.Customer) error
}

type FilterCustomers struct {
	log      logger.Logger
	notifier FilterCustomersNotifier
}

func NewFilterCustomers(log logger.Logger, notifier FilterCustomersNotifier) *FilterCustomers {
	return &FilterCustomers{
		log:      log,
		notifier: notifier,
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

			if err := f.notifier.Notify(ctx, &customer); err != nil {
				log.Infof("Error to notify customer invited id=%d", customer.ID)
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

	if err := f.sort(result, orderBy); err != nil {
		return nil, errors.Wrap(err, "error to sort result")
	}

	return result, nil
}

func (f *FilterCustomers) sort(result domain.Customers, orderBy domain.OrderBy) error {
	switch orderBy {
	case domain.OrderByCustomerID:
		sort.Slice(result, func(i, j int) bool {
			return result[i].ID < result[j].ID
		})

	default:
		return errors.New("unexpected order by")
	}

	return nil
}

func customersMapValues(customersMap map[int]domain.Customer) domain.Customers {
	var r = make([]domain.Customer, 0)

	for _, distance := range customersMap {
		r = append(r, distance)
	}

	return r
}
