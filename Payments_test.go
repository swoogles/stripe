package stripe

import (
	"testing"
)

func TestGetAllTestProducts(*testing.T) {
	GetAllProductsWithUnsafeType("TEST_STRIPE_SECRET_KEY", "good")
}
