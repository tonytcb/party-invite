package domain

// To calculate the distance between two points, we are using the Haversine formula described
// here https://en.wikipedia.org/wiki/Haversine_formula, and here https://en.wikipedia.org/wiki/Great-circle_distance
// you can read more about the Great-circle Distance (shortest distance).
//
// Instead of making all math using float numbers, we are using the shopspring/decimal package,
// avoiding floating point number precision issues.

import (
	"math"

	"github.com/shopspring/decimal"
)

const (
	earthRadiusInKm = 6371
)

var (
	pi = decimal.NewFromFloat(math.Pi)
)

// distance calculates the distance, in kilometers, between two coordinates applying the Haversine formula.
func distance(c1 *Coordinate, c2 *Coordinate) decimal.Decimal {
	var (
		earthRadius     = decimal.NewFromInt32(earthRadiusInKm)
		oneDecimalValue = decimal.NewFromInt32(1) // nolint:golint,gomnd // it's only a definition of the one as decimal type
		twoDecimalValue = decimal.NewFromInt32(2) // nolint:golint,gomnd // it's only a definition of the two as decimal type

		latitude1  = toRadians(c1.Latitude)
		longitude1 = toRadians(c1.Longitude)
		latitude2  = toRadians(c2.Latitude)
		longitude2 = toRadians(c2.Longitude)

		latitudeDiff  = latitude2.Add(latitude1.Neg())
		longitudeDiff = longitude2.Add(longitude1.Neg())

		a1 = latitudeDiff.Div(twoDecimalValue).Sin().Pow(twoDecimalValue)
		a2 = latitude1.Cos().Mul(latitude2.Cos())
		a3 = longitudeDiff.Div(twoDecimalValue).Sin().Pow(twoDecimalValue)
		a  = a2.Mul(a3).Add(a1)

		arcTangent = decimalAtan2(
			decimalSqrt(a),
			decimalSqrt(oneDecimalValue.Add(a.Neg())),
		).Mul(twoDecimalValue)
	)

	return arcTangent.Mul(earthRadius)
}

// toRadians converts a degree value in radians.
func toRadians(degreeValue decimal.Decimal) decimal.Decimal {
	const angle = 180
	return degreeValue.Mul(pi).Div(decimal.NewFromFloat(angle))
}

// decimalAtan2 makes use of the native math.atan2 (arc tangent) converting its result to Decimal type.
func decimalAtan2(v1 decimal.Decimal, v2 decimal.Decimal) decimal.Decimal {
	fSqrtV1, _ := v1.Float64()
	fSqrtV2, _ := v2.Float64()
	return decimal.NewFromFloat(math.Atan2(fSqrtV1, fSqrtV2))
}

// decimalAtan2 makes use of the native math.Sqrt (square root) converting its result to Decimal type.
func decimalSqrt(v decimal.Decimal) decimal.Decimal {
	fValue, _ := v.Float64() // ignoring loss precision here
	return decimal.NewFromFloat(math.Sqrt(fValue))
}
