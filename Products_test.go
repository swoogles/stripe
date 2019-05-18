package stripe

import (
	"fmt"
	"testing"
)

func TestGetAllProducts(t *testing.T) {

	fmt.Println(JsonSerialize(GetAllProductsWithUnsafeType("STRIPE_SECRET_KEY", "good")))
}
