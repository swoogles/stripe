package stripe

import (
	"fmt"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/product"
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
	fmt.Println("New Product Response")
	fmt.Println(prod)
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
