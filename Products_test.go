package stripe

import (
	"fmt"
	"testing"
)

func TestGetAllProducts(t *testing.T) {

	fmt.Println(JsonSerialize(GetActiveProductsWithUnsafeType("STRIPE_SECRET_KEY", "good")))
}
