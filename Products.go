package stripe

import (
	"fmt"
	"github.com/leekchan/accounting"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/plan"
	"github.com/stripe/stripe-go/product"
	"github.com/stripe/stripe-go/sku"
	"os"
)

type Product struct {
	Name string
	Skus []Sku
}

type Plan struct {
	Name     string
	Id       string
	Interval string
	Price    string
}

func CreateTestProduct() {
	CreateProduct("TEST_STRIPE_SECRET_KEY")
}

//func GetAllTestProducts() []Product {
//	returnGetAllProducts("TEST_STRIPE_SECRET_KEY")
//}

func CreateProduct(testKey string) {
	stripe.Key = os.Getenv(testKey)

	params := &stripe.ProductParams{
		Name: stripe.String("Gym Membership"),
		Type: stripe.String(string(stripe.ProductTypeService)),
	}
	prod, _ := product.New(params)
	//prod.ID

	fmt.Println("New Product Response")
	fmt.Println(prod)
}

func GetAllProductsWithUnsafeType(stripePaymentToken string, productTypeString string) []Product {
	return GetAllProducts(stripePaymentToken, stripe.ProductType(productTypeString))
}

func GetAllProducts(stripePaymentToken string, productType stripe.ProductType) []Product {
	stripe.Key = os.Getenv(stripePaymentToken)
	productTypeString := string(productType)

	shippable := false
	params := &stripe.ProductListParams{
		Type:      &productTypeString,
		Shippable: &shippable,
	}
	i := product.List(params)

	productList := make([]Product, i.Meta().TotalCount)

	for i.Next() {
		p := i.Product()
		skuParams := &stripe.SKUListParams{
			Product: &p.ID,
			// TODO
		}
		skuResponse := sku.List(skuParams)

		curProduct := Product{Name: p.Name}
		for skuResponse.Next() {
			curProduct.Skus = append(curProduct.Skus, createSkuFrom(skuResponse.SKU()))
		}
		productList = append(productList, curProduct)
	}
	return productList
}

func FormatPrice(amount float64) string {
	ac := accounting.Accounting{Symbol: "$", Precision: 2}
	return ac.FormatMoney(amount)
}

func GetAllPlans(stripePaymentToken string) []Plan {
	stripe.Key = os.Getenv(stripePaymentToken)

	params := &stripe.PlanListParams{}
	i := plan.List(params)
	plan.List(params)

	productList := make([]Plan, i.Meta().TotalCount)

	for i.Next() {
		p := i.Plan()
		var interval string
		if p.IntervalCount > 1 {
			interval = string(p.Interval) + "s"
		} else {
			interval = string(p.Interval)
		}
		curProduct := Plan{Name: p.Nickname, Id: p.ID, Interval: interval, Price: FormatPrice(float64(p.Amount) / 100.0)}
		productList = append(productList, curProduct)
	}
	return productList
}
