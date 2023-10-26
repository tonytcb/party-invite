package customernotify

import (
	"context"
	"fmt"

	"github.com/tonytcb/party-invite/pkg/domain"
	"github.com/tonytcb/party-invite/pkg/infrastructure/logger"
)

type StdOutNotifier struct {
	log logger.Logger
}

func NewStdOutNotifier(log logger.Logger) *StdOutNotifier {
	return &StdOutNotifier{
		log: log,
	}
}

func (s *StdOutNotifier) Notify(ctx context.Context, customer *domain.Customer) error {
	notification := fmt.Sprintf("[STD OUT NOTIFICATION] customer %s invited", customer.Name)

	fmt.Println(notification) // sends notification

	s.log.FromContext(ctx).Infof("Customer %d successfully notified", customer.ID)

	return nil
}
