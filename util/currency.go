package util

const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
	INR = "INR"
)

// IsCurrencySupported returns true if currency is supported
func IsCurrencySupported(currency string) bool {
	switch currency {
	case USD, CAD, EUR, INR:
		return true
	default:
		return false
	}
}
