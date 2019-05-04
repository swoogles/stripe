package stripe

import (
	"fmt"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/plan"
	"os"
)

func CreateTestPlan(productId string) string {
	return createPlan("TEST_STRIPE_SECRET_KEY", productId)
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
