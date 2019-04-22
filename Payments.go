package stripe

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/order"
	"github.com/stripe/stripe-go/plan"
	"github.com/stripe/stripe-go/product"
	"github.com/stripe/stripe-go/sku"
	"github.com/stripe/stripe-go/sub"
	"log"
	"os"
)

func createTestPaymentFunction(testKey string) func(string, int64) string {
	return func(stripePaymentToken string, amount int64) string {
		stripe.Key = os.Getenv(testKey)

		params := &stripe.ChargeParams{
			Amount:      stripe.Int64(amount),
			Currency:    stripe.String(string(stripe.CurrencyUSD)),
			Description: stripe.String("Example charge"),
		}
		params.SetSource(stripePaymentToken)
		ch, _ := charge.New(params)
		fmt.Println(ch)

		return "faux payment response"
	}

}

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

func createPlan(key string, productId string) string {
	stripe.Key = os.Getenv(key)
	params := &stripe.PlanParams{
		ProductID: stripe.String(productId),
		Nickname:  stripe.String("Gym Membership USD"),
		Interval:  stripe.String(string(stripe.PlanIntervalMonth)),
		Currency:  stripe.String("usd"),
		Amount:    stripe.Int64(10000),
	}
	p, _ := plan.New(params)
	fmt.Println("New Plan:")
	fmt.Println(p)
	return p.ID
}

//func createSource(testKey string) string {
//	&stripe.SourceParams{}
//	source.New()
//}

//func CreateSource(stripeSecretKey string) {
//	&stripe.SourceParams{}
//}

func CreateCustomer(key string, sourceToken string, email string) string {
	stripe.Key = os.Getenv(key)

	params := &stripe.CustomerParams{
		Email: stripe.String("jenny.rosen@example.com"),
	}
	params.SetSource(sourceToken)
	cus, errors := customer.New(params)
	if errors != nil {
		fmt.Println(errors)
	}
	fmt.Println("New Customer: ")
	fmt.Println(cus)
	return cus.ID
}

func createSubscription(testkey string) func(string, string) string {
	return func(planId string, customerId string) string {

		items := []*stripe.SubscriptionItemsParams{
			{
				Plan: stripe.String(planId),
			},
		}
		params := &stripe.SubscriptionParams{
			Customer: stripe.String(customerId),
			Items:    items,
		}
		newSubscription, _ := sub.New(params)
		fmt.Println("New Subscription: ")
		fmt.Println(newSubscription)
		return newSubscription.ID
	}
}

func CreateTestSubscription(planId string, customerId string) string {
	createSubscriptionFor := createSubscription("TEST_STRIPE_SECRET_KEY")
	return createSubscriptionFor(planId, customerId)
}

func CreateTestCustomer(sourceToken string) string {
	return CreateCustomer("TEST_STRIPE_SECRET_KEY", sourceToken, "")
}

func CreateTestPlan(productId string) string {
	return createPlan("TEST_STRIPE_SECRET_KEY", productId)
}

func CreateTestProduct() {
	CreateProduct("TEST_STRIPE_SECRET_KEY")
}

func ExecuteTestStripePaymentWithAmount(stripePaymentToken string, amount int64) string {
	return createTestPaymentFunction("TEST_STRIPE_SECRET_KEY")(stripePaymentToken, amount)
}

func ExecuteLiveStripePaymentWithAmount(stripePaymentToken string, amount int64) string {
	return createTestPaymentFunction("LIVE_STRIPE_SECRET_KEY")(stripePaymentToken, amount)
}

//func GetAllTestProducts() []Product {
//	return GetAllProducts("TEST_STRIPE_SECRET_KEY")
//}

type Sku struct {
	ID          string
	Price       int64
	Description string
	Membership  string
}
type Product struct {
	Name string
	Skus []Sku
}

func createSkuFrom(stripeSku *stripe.SKU) Sku {

	var membership string // TODO Can it be immutable?
	if value, ok := stripeSku.Attributes["membership"]; ok {
		membership = value
	} else {
		membership = "Anyone"
	}
	return Sku{
		Price:       stripeSku.Price,
		ID:          stripeSku.ID,
		Description: stripeSku.Description,
		Membership:  membership,
	}
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

func ExecuteStripePaymentWithAmount(stripePaymentToken string, amount int64, stripeSecretKeyVariable string) string {
	stripe.Key = os.Getenv(stripeSecretKeyVariable)

	params := &stripe.ChargeParams{
		Amount:      stripe.Int64(amount),
		Currency:    stripe.String(string(stripe.CurrencyUSD)),
		Description: stripe.String("Example charge"),
	}
	params.SetSource(stripePaymentToken)
	ch, _ := charge.New(params)
	fmt.Println(ch)

	return "payment:" + string(amount)
}

func CreateOrder(sku string, stripeSecretKeyVariable string, customerId string) string {
	stripe.Key = os.Getenv(stripeSecretKeyVariable)

	params := &stripe.OrderParams{
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		Email:    stripe.String("jenny.rosen@example.com"),
		Items: []*stripe.OrderItemParams{
			{
				Type:     stripe.String(string(stripe.OrderItemTypeSKU)),
				Parent:   &sku,
				Quantity: stripe.Int64(1),
			},
		},
	}
	ord, _ := order.New(params)
	fmt.Println(JsonSerialize(ord))
	orderPayParams := stripe.OrderPayParams{
		Customer: &customerId,
	}
	pay, e := order.Pay(ord.ID, &orderPayParams)
	fmt.Println(pay)
	if e != nil {
		log.Println(e)
	}
	return pay.ID

}

func JsonSerialize(class interface{}) string {
	out, _ := json.Marshal(class)
	var prettyJSON bytes.Buffer
	prettyPrintError := json.Indent(&prettyJSON, out, "", "\t")
	if prettyPrintError != nil {
		log.Println("JSON parse error: ", prettyPrintError)
	}
	return string(prettyJSON.Bytes())
}
