package stripe

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/order"
	"log"
	"os"
)

type Sku struct {
	ID          string
	Price       int64
	Description string
	Membership  string
}

func createPayment(testKey string, stripePaymentToken string, amount int64) string {
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

//func createSource(testKey string) string {
//	&stripe.SourceParams{}
//	source.New()
//}

//func CreateSource(stripeSecretKey string) {
//	&stripe.SourceParams{}
//}

func ExecuteTestStripePaymentWithAmount(stripePaymentToken string, amount int64) string {
	return createPayment("TEST_STRIPE_SECRET_KEY", stripePaymentToken, amount)
}

func ExecuteLiveStripePaymentWithAmount(stripePaymentToken string, amount int64) string {
	return createPayment("LIVE_STRIPE_SECRET_KEY", stripePaymentToken, amount)
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
