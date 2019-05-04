package stripe

import (
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/product"
	"github.com/stripe/stripe-go/sku"
	"golang.org/x/tools/go/ssa/interp/testdata/src/fmt"
	"golang.org/x/tools/go/ssa/interp/testdata/src/os"
)

type Product struct {
	Name string
	Skus []Sku
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

	params := &stripe.ProductListParams{
		Type: &productTypeString,
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
