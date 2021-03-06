package stripe

import (
	"fmt"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/product"
	"github.com/stripe/stripe-go/sku"
	"os"
)

type Product struct {
	Name        string
	Id          string
	Price       string
	Description string
	//Membership  string
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

func GetActiveProductsWithUnsafeType(stripeSecretKey string, productTypeString string) []Product {
	return GetActiveProducts(stripeSecretKey, stripe.ProductType(productTypeString))
}

func GetActiveProducts(stripeSecretKey string, productType stripe.ProductType) []Product {
	fmt.Println("retrieved key: " + os.Getenv(stripeSecretKey))
	stripe.Key = os.Getenv(stripeSecretKey)
	productTypeString := string(productType)

	active := true
	shippable := false
	params := &stripe.ProductListParams{
		Type:      &productTypeString,
		Shippable: &shippable,
		Active:    &active,
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
		if skuResponse.Next() {
			fmt.Println("In condition")
			curProduct.Id = skuResponse.SKU().ID
			curProduct.Price = FormatPrice(skuResponse.SKU().Price)
			curProduct.Description = p.Description
		}
		productList = append(productList, curProduct)
	}
	return productList
}
