package bb

import (
	"testing"

	"github.com/mundipagg/boleto-api/test"
)

func TestShouldCalculateAgencyDigitFromBb(t *testing.T) {
	test.ExpectTrue(bbAgencyDigitCalculator("0137") == "6", t)

	test.ExpectTrue(bbAgencyDigitCalculator("3418") == "5", t)

	test.ExpectTrue(bbAgencyDigitCalculator("3324") == "3", t)

	test.ExpectTrue(bbAgencyDigitCalculator("5797") == "5", t)
}

func TestShouldCalculateAccountDigitFromBb(t *testing.T) {
	test.ExpectTrue(bbAccountDigitCalculator("", "00006685") == "0", t)

	test.ExpectTrue(bbAccountDigitCalculator("", "00025619") == "6", t)

	test.ExpectTrue(bbAccountDigitCalculator("", "00006842") == "X", t)

	test.ExpectTrue(bbAccountDigitCalculator("", "00000787") == "0", t)
}
