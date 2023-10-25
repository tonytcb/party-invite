package domain

import (
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type Coordinate struct {
	Latitude  decimal.Decimal
	Longitude decimal.Decimal
}

func NewCoordinate(latitude string, longitude string) (*Coordinate, error) {
	latitudeDecimal, err := decimal.NewFromString(latitude)
	if err != nil {
		return nil, NewErrInvalidArgument(err.Error(), "invalid latitude")
	}

	longitudeDecimal, err := decimal.NewFromString(longitude)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid longitude (%s) value as number", longitude)
	}

	return &Coordinate{
		Latitude:  latitudeDecimal,
		Longitude: longitudeDecimal,
	}, nil
}

func (c *Coordinate) Difference(c2 *Coordinate) decimal.Decimal {
	return distance(c, c2)
}

var (
	DublinLocation = &Coordinate{
		Latitude:  decimal.RequireFromString("53.339428"),
		Longitude: decimal.RequireFromString("-6.257664"),
	}
)
