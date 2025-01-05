package extendedtypes

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

const maxPrecisionOnDiv = 10

type SafeDecimal struct {
	amount        int64
	decimalPlaces int
}

var SafeDecimalValueZero = SafeDecimal{
	amount:        0,
	decimalPlaces: 0,
}

func SafeDecimalFromString(value string) (SafeDecimal, error) {
	parts := strings.Split(value, ".")

	var valueToConvert string
	var decimalPlaces int

	if len(parts) == 1 {
		decimalPlaces = 0
		valueToConvert = parts[0]
	} else if len(parts) != 2 {
		return SafeDecimalValueZero, fmt.Errorf("String is not valid DecimalValue format")
	} else {
		decimalPlaces = len(parts[1])
		valueToConvert = parts[0] + parts[1]
	}

	numericWhole, err := strconv.ParseInt(valueToConvert, 10, 64)
	if err != nil {
		return SafeDecimalValueZero, err
	}

	return SafeDecimal{
		amount:        int64(numericWhole),
		decimalPlaces: decimalPlaces,
	}, nil
}

func intPow(x, y int) int64 {
	return int64(math.Pow(float64(x), float64(y)))
}

func intAbs(x int64) int64 {
	if x > 0 {
		return x
	}

	return -x
}

func (d SafeDecimal) String() string {
	return fmt.Sprintf("%d.%0*d", d.amount/intPow(10, d.decimalPlaces), d.decimalPlaces, intAbs(d.amount%intPow(10, d.decimalPlaces)))
}

func (d SafeDecimal) Format(decimalPlaces int) string {
	if decimalPlaces > d.decimalPlaces {
		d.amount *= intPow(10, decimalPlaces)
		d.decimalPlaces = decimalPlaces
	}

	if decimalPlaces < d.decimalPlaces {
		factor := intPow(10, d.decimalPlaces-decimalPlaces)

		rest := (d.amount / (factor / 10)) % 10

		d.amount /= factor
		d.decimalPlaces = decimalPlaces

		if rest > 4 {
			d.amount += 1
		}
	}

	if decimalPlaces > 0 {
		return fmt.Sprintf("%d.%0*d", d.amount/intPow(10, d.decimalPlaces), d.decimalPlaces, intAbs(d.amount%intPow(10, d.decimalPlaces)))
	}

	return fmt.Sprintf("%d", d.amount/intPow(10, d.decimalPlaces))
}

func (d SafeDecimal) FormatRoundingDown(decimalPlaces int) string {
	if decimalPlaces > d.decimalPlaces {
		d.amount *= intPow(10, decimalPlaces)
		d.decimalPlaces = decimalPlaces
	}

	if decimalPlaces < d.decimalPlaces {
		factor := intPow(10, d.decimalPlaces-decimalPlaces)
		d.amount /= factor
		d.decimalPlaces = decimalPlaces
	}

	if decimalPlaces == 0 {
		return fmt.Sprintf("%d", d.amount/intPow(10, d.decimalPlaces))
	}

	return fmt.Sprintf("%d.%0*d", d.amount/intPow(10, d.decimalPlaces), d.decimalPlaces, intAbs(d.amount%intPow(10, d.decimalPlaces)))
}

func (d1 SafeDecimal) Add(d2 SafeDecimal) SafeDecimal {
	if d1.decimalPlaces < d2.decimalPlaces {
		return SafeDecimal{
			amount:        d2.amount + intPow(10, d2.decimalPlaces-d1.decimalPlaces)*d1.amount,
			decimalPlaces: d2.decimalPlaces,
		}
	}

	return SafeDecimal{
		amount:        d1.amount + intPow(10, d1.decimalPlaces-d2.decimalPlaces)*d2.amount,
		decimalPlaces: d1.decimalPlaces,
	}
}

func (d1 SafeDecimal) Subtract(d2 SafeDecimal) SafeDecimal {
	var v1, v2 int64
	var decimalPlaces int
	if d1.decimalPlaces < d2.decimalPlaces {
		v1 = intPow(10, d2.decimalPlaces-d1.decimalPlaces) * d1.amount
		v2 = d2.amount
		decimalPlaces = d2.decimalPlaces
	} else {
		v1 = d1.amount
		v2 = intPow(10, d1.decimalPlaces-d2.decimalPlaces) * d2.amount
		decimalPlaces = d1.decimalPlaces
	}

	if v1-v2 == 0 {
		decimalPlaces = 0
	}

	return SafeDecimal{
		amount:        v1 - v2,
		decimalPlaces: decimalPlaces,
	}
}

func (d *SafeDecimal) removeRightZeros() {
	s := strconv.FormatInt(d.amount, 10)
	t := strings.TrimRight(s, "0")
	a := len(s) - len(t)

	d.amount /= intPow(10, a)

	if d.amount != 0 {
		d.decimalPlaces -= a
		return
	}

	d.decimalPlaces = 0
}

func (d1 SafeDecimal) Multiply(d2 SafeDecimal) SafeDecimal {
	amount := d1.amount * d2.amount
	r := SafeDecimal{
		amount:        amount,
		decimalPlaces: d1.decimalPlaces + d2.decimalPlaces,
	}
	r.removeRightZeros()

	return r
}

func (d1 SafeDecimal) Divide(d2 SafeDecimal) SafeDecimal {
	amount := d1.amount / d2.amount
	rest := d1.amount % d2.amount
	decimalPlaces := d1.decimalPlaces - d2.decimalPlaces

	for rest > 0 && decimalPlaces < maxPrecisionOnDiv {
		rest *= 10
		amount *= 10
		decimalPlaces += 1

		amount += rest / d2.amount
		rest %= d2.amount
	}

	r := SafeDecimal{
		amount:        amount,
		decimalPlaces: decimalPlaces,
	}
	r.removeRightZeros()

	return r
}
