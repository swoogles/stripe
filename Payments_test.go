package stripe

import (
	"fmt"
	"testing"
)

func TestGetAllTestProducts(*testing.T) {
	allGoods := GetAllProductsWithUnsafeType("TEST_STRIPE_SECRET_KEY", "good")
	fmt.Println(string(JsonSerialize(allGoods)))
}
