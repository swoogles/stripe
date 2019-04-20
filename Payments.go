package stripe

import (
	"fmt"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/product"
	"github.com/stripe/stripe-go/plan"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/sub"
	//"github.com/stripe/stripe-go/source"
	"os"
	"encoding/json"
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

func createCustomer(key string, sourceToken string) string {
	stripe.Key = os.Getenv(key)

	params := &stripe.CustomerParams{
		Email: stripe.String("jenny.rosen@example.com"),
	}
	params.SetSource(sourceToken)
	cus, _ := customer.New(params)
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
	return createCustomer("TEST_STRIPE_SECRET_KEY", sourceToken)
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

func GetAllTestProducts() []Product {
	return getAllProducts("TEST_STRIPE_SECRET_KEY")
}

type Product struct {
	Name string
}

func getAllProducts(stripePaymentToken string) []Product {
	stripe.Key = os.Getenv(stripePaymentToken)

	var productList []Product
	params := &stripe.ProductListParams{}
	params.Filters.AddFilter("limit", "", "3")
	i := product.List(params)

	for i.Next() {
		p := i.Product()
		curProduct := Product{Name: p.Name}
		productList = append(productList, curProduct)
		out, _ := json.Marshal(curProduct)
		fmt.Println("curProduct: " + string(out))
	}
	fmt.Println("products: " + string(len(productList)))
	return productList
}

/* Slice fuckery
func getAllProducts(stripePaymentToken string) []Product {
	stripe.Key = os.Getenv(stripePaymentToken)

	params := &stripe.ProductListParams{}
	params.Filters.AddFilter("limit", "", "3")
	i := product.List(params)
	products := i[:0]

	for i, n := range i {
		//for i.Next() {
		p := i.Product()
		curProduct := Product{Name: p.Name}
		out, _ := json.Marshal(curProduct)
		fmt.Println("curProduct: " + string(out))
		products = append(products, curProduct)
	}
	fmt.Println("Number of products: " + string(len(products)))
	return products
}
*/
