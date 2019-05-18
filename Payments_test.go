package stripe

import (
	"fmt"
	"testing"
)

func TestGetAllTestProducts(*testing.T) {
	allGoods := GetAllProductsWithUnsafeType("TEST_STRIPE_SECRET_KEY",
		"good",
	)
	fmt.Println(string(JsonSerialize(allGoods)))
}

func TestGetAllTestServices(*testing.T) {
	allServices := GetAllProductsWithUnsafeType("TEST_STRIPE_SECRET_KEY",
		"service",
	)
	fmt.Println(string(JsonSerialize(allServices)))
}

func TestGetAllTestPlans(*testing.T) {
	allGoods := GetAllPlans("TEST_STRIPE_SECRET_KEY")
	fmt.Println(string(JsonSerialize(allGoods)))
}

func TestCreateOrder(t *testing.T) {
	CreateOrder("sku_Euy4ksCdJWGUnZ", "TEST_STRIPE_SECRET_KEY", "cus_F0Xqt5aj0CMm50")
}
