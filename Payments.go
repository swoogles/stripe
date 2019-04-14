package stripe

import (
	"fmt"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/product"
	"github.com/stripe/stripe-go/plan"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/sub"
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

func createPlan(key string, productId) string {
	stripe.Key = os.Getenv(key)
	params := &stripe.PlanParams{
		ProductID: stripe.String(productId),
		Nickname:  stripe.String("Gym Membership USD"),
		Interval:  stripe.String(string(stripe.PlanIntervalMonth)),
		Currency:  stripe.String("usd"),
		Amount:    stripe.Int64(10000),
	}
	p, err := plan.New(params)
	fmt.Println("New Plan:")
	fmt.Println(p)
	fmt.Println("New Plan error:" + err.Error())
	return p.ID
}

func createCustomer(testKey string) string {
	stripe.Key = os.Getenv(testKey)

	params := &stripe.CustomerParams{
		Email: stripe.String("jenny.rosen@example.com"),
	}
	params.SetSource("src_18eYalAHEMiOZZp1l9ZTjSU0")
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

func CreateTestCustomer() string {
	return createCustomer("TEST_STRIPE_SECRET_KEY")
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
