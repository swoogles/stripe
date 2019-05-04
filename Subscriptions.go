package stripe

import (
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/sub"
	"golang.org/x/tools/go/ssa/interp/testdata/src/fmt"
	"golang.org/x/tools/go/ssa/interp/testdata/src/os"
)

func CreateSubscription(stripeSecretKey string, planId string, customerId string) string {
	stripe.Key = os.Getenv(stripeSecretKey)
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
