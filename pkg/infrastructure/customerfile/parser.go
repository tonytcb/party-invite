package customerfile

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"

	"git.codesubmit.io/sfox/golang-party-invite-ivsjhn/pkg/domain"
)

type rawCustomer struct {
	UserID    int    `json:"user_id"`
	Name      string `json:"name"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

type CustomersFileParser struct {
}

func NewCustomersFileParser() *CustomersFileParser {
	return &CustomersFileParser{}
}

// Parse parses a file into a list of customers, reading each line separately.
func (c CustomersFileParser) Parse(ctx context.Context, file io.Reader) (domain.Customers, error) {
	var (
		customers   = make([]domain.Customer, 0)
		fileScanner = bufio.NewScanner(file)
		i           = 1
	)

	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		select {
		case <-ctx.Done():
			return nil, errors.New("context done while parsing file")
		default:
		}

		var lineContents = fileScanner.Bytes()

		if len(lineContents) == 0 {
			continue // empty line
		}

		var line = &rawCustomer{}

		if err := json.Unmarshal(lineContents, &line); err != nil {
			return nil, domain.NewErrInvalidArgument(
				err.Error(),
				fmt.Sprintf("error to parse line=%d, content='%s'", i, string(lineContents)),
			)
		}

		location, err := domain.NewCoordinate(line.Latitude, line.Longitude)
		if err != nil {
			return nil, errors.Wrap(err, "error to parse customers' location")
		}

		customers = append(customers, domain.NewCustomer(line.UserID, line.Name, location))

		i++
	}

	return customers, nil
}
