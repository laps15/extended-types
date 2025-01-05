package extendedtypes_test

import (
	"testing"

	et "github.com/laps15/extended-types"
)

func secureSafeDecimalFromString(v string) et.SafeDecimal {
	m, _ := et.SafeDecimalFromString(v)

	return m
}

type invalidStringTest struct {
	v string
}

var invalidStringTests = []invalidStringTest{
	{"1.1.1"},
	{"invalid number"},
}

func TestInvalidSafeDecimalFromString(t *testing.T) {
	var err error
	for _, test := range invalidStringTests {
		_, err = et.SafeDecimalFromString(test.v)
		if err == nil {
			t.Errorf("Processing %q, should have resulted in error.", test.v)
		}
	}
}

type validFromStringTestCase struct {
	v, expected string
}

var validFromStringTests = []validFromStringTestCase{
	{"12.345", "12.345"},
	{".0123", "0.0123"},
}

func TestSafeDecimalFromString(t *testing.T) {
	for _, test := range validFromStringTests {
		got, err := et.SafeDecimalFromString(test.v)
		if err != nil {
			t.Errorf("Processing %q, resulted in an unexpected error: %q", test.v, err)
		}

		if got.String() != test.expected {
			t.Errorf("got %q, wanted %q", got, test.expected)
		}
	}
}

type formatTestCase struct {
	v             et.SafeDecimal
	decimalPlaces int
	expected      string
}

var formatTests = []formatTestCase{
	{secureSafeDecimalFromString("1003658.53658536585"), 0, "1003659"},
	{secureSafeDecimalFromString("1003658.53658536585"), 1, "1003658.5"},
	{secureSafeDecimalFromString("1003658.53658536585"), 2, "1003658.54"},
	{secureSafeDecimalFromString("1003658.53658536585"), 3, "1003658.537"},
	{secureSafeDecimalFromString("1003658.53658536585"), 4, "1003658.5366"},
	{secureSafeDecimalFromString("1003658.53658536585"), 5, "1003658.53659"},
	{secureSafeDecimalFromString("1003658.53658536585"), 6, "1003658.536585"},
	{secureSafeDecimalFromString("1003658.53658536585"), 7, "1003658.5365854"},
	{secureSafeDecimalFromString("1003658.53658536585"), 8, "1003658.53658537"},
	{secureSafeDecimalFromString("1003658.53658536585"), 9, "1003658.536585366"},
	{secureSafeDecimalFromString("1003658.53658536585"), 10, "1003658.5365853659"},
	{secureSafeDecimalFromString("1003658.53658536585"), 11, "1003658.53658536585"},
	{secureSafeDecimalFromString("1.6"), 0, "2"},
	{secureSafeDecimalFromString("1"), 5, "1.00000"},
}

func TestFormat(t *testing.T) {
	for _, test := range formatTests {
		got := test.v.Format(test.decimalPlaces)

		if got != test.expected {
			t.Errorf("got %q, wanted %q", got, test.expected)
		}
	}
}

var formatRoundingDownTests = []formatTestCase{
	{secureSafeDecimalFromString("1003658.53658536585"), 0, "1003658"},
	{secureSafeDecimalFromString("1003658.53658536585"), 1, "1003658.5"},
	{secureSafeDecimalFromString("1003658.53658536585"), 2, "1003658.53"},
	{secureSafeDecimalFromString("1003658.53658536585"), 3, "1003658.536"},
	{secureSafeDecimalFromString("1003658.53658536585"), 4, "1003658.5365"},
	{secureSafeDecimalFromString("1003658.53658536585"), 5, "1003658.53658"},
	{secureSafeDecimalFromString("1003658.53658536585"), 6, "1003658.536585"},
	{secureSafeDecimalFromString("1003658.53658536585"), 7, "1003658.5365853"},
	{secureSafeDecimalFromString("1003658.53658536585"), 8, "1003658.53658536"},
	{secureSafeDecimalFromString("1003658.53658536585"), 9, "1003658.536585365"},
	{secureSafeDecimalFromString("1003658.53658536585"), 10, "1003658.5365853658"},
	{secureSafeDecimalFromString("1003658.53658536585"), 11, "1003658.53658536585"},
	{secureSafeDecimalFromString("1.6"), 0, "1"},
	{secureSafeDecimalFromString("1"), 5, "1.00000"},
}

func TestFormatRoundingDown(t *testing.T) {
	for _, test := range formatRoundingDownTests {
		got := test.v.FormatRoundingDown(test.decimalPlaces)

		if got != test.expected {
			t.Errorf("got %q, wanted %q", got, test.expected)
		}
	}
}

type toStringTest struct {
	v        et.SafeDecimal
	expected string
}

var toStringTests = []toStringTest{
	{secureSafeDecimalFromString("12.345"), "12.345"},
	{secureSafeDecimalFromString("12345.010"), "12345.010"},
	{secureSafeDecimalFromString("12345.0123"), "12345.0123"},
	{secureSafeDecimalFromString("-12345.0123"), "-12345.0123"},
}

func TestToString(t *testing.T) {
	for _, test := range toStringTests {
		got := test.v.String()
		if got != test.expected {
			t.Errorf("got %q, wanted %q", got, test.expected)
		}
	}
}

type safeDecimalOperationsTest struct {
	d1, d2, expected et.SafeDecimal
}

var safeDecimalAddTests = []safeDecimalOperationsTest{
	{et.SafeDecimalValueZero, secureSafeDecimalFromString("12.345"), secureSafeDecimalFromString("12.345")},
	{secureSafeDecimalFromString("12.345"), secureSafeDecimalFromString("12.345"), secureSafeDecimalFromString("24.690")},
	{secureSafeDecimalFromString("12.345"), secureSafeDecimalFromString("1"), secureSafeDecimalFromString("13.345")},
	{secureSafeDecimalFromString("12345"), secureSafeDecimalFromString("0.0123"), secureSafeDecimalFromString("12345.0123")},
	{secureSafeDecimalFromString("1"), secureSafeDecimalFromString("-1"), secureSafeDecimalFromString("0")},
}

func TestSafeDecimalAdd(t *testing.T) {
	for _, test := range safeDecimalAddTests {
		got := test.d1.Add(test.d2)
		if got != test.expected {
			t.Errorf("got %q, wanted %q", got, test.expected)
		}
	}
}

var safeDecimalSubtractTests = []safeDecimalOperationsTest{
	{et.SafeDecimalValueZero, secureSafeDecimalFromString("12.345"), secureSafeDecimalFromString("-12.345")},
	{secureSafeDecimalFromString("12.345"), secureSafeDecimalFromString("12.345"), et.SafeDecimalValueZero},
	{secureSafeDecimalFromString("12.345"), secureSafeDecimalFromString("1"), secureSafeDecimalFromString("11.345")},
	{secureSafeDecimalFromString("12345"), secureSafeDecimalFromString("0.0123"), secureSafeDecimalFromString("12344.9877")},
	{secureSafeDecimalFromString("-1"), secureSafeDecimalFromString("-1"), secureSafeDecimalFromString("0")},
}

func TestSafeDecimalSubtract(t *testing.T) {
	for _, test := range safeDecimalSubtractTests {
		got := test.d1.Subtract(test.d2)
		if got != test.expected {
			t.Errorf("got %q, wanted %q", got, test.expected)
		}
	}
}

var safeDecimalMultiplyTests = []safeDecimalOperationsTest{
	{secureSafeDecimalFromString("12.345"), et.SafeDecimalValueZero, et.SafeDecimalValueZero},
	{secureSafeDecimalFromString("12.345"), secureSafeDecimalFromString("0.1"), secureSafeDecimalFromString("1.2345")},
	{secureSafeDecimalFromString("12.345"), secureSafeDecimalFromString("1"), secureSafeDecimalFromString("12.345")},
	{secureSafeDecimalFromString("12.345"), secureSafeDecimalFromString("10"), secureSafeDecimalFromString("123.45")},
	{secureSafeDecimalFromString("12345"), secureSafeDecimalFromString("0.0123"), secureSafeDecimalFromString("151.8435")},
}

func TestSafeDecimalMultiply(t *testing.T) {
	for _, test := range safeDecimalMultiplyTests {
		got := test.d1.Multiply(test.d2)
		if got != test.expected {
			t.Errorf("got %q, wanted %q", got, test.expected)
		}
	}
}

var safeDecimalDivideTests = []safeDecimalOperationsTest{
	{secureSafeDecimalFromString("12.345"), secureSafeDecimalFromString("0.1"), secureSafeDecimalFromString("123.45")},
	{secureSafeDecimalFromString("12.345"), secureSafeDecimalFromString("1"), secureSafeDecimalFromString("12.345")},
	{secureSafeDecimalFromString("12.345"), secureSafeDecimalFromString("10"), secureSafeDecimalFromString("1.2345")},
	{secureSafeDecimalFromString("12345"), secureSafeDecimalFromString("0.0123"), secureSafeDecimalFromString("1003658.5365853658")},
}

func TestSafeDecimalDivide(t *testing.T) {
	for _, test := range safeDecimalDivideTests {
		got := test.d1.Divide(test.d2)
		if got != test.expected {
			t.Errorf("got %q, wanted %q", got, test.expected)
		}
	}
}
