package stripe

import (
	"fmt"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/plan"
	"os"
	"strconv"
)

type Plan struct {
	Name     string
	Id       string
	Interval string
	Price    string
}

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
			interval = strconv.FormatInt(p.IntervalCount, 10) + " " + string(p.Interval) + "s"
		} else {
			interval = string(p.Interval)
		}
		curProduct := Plan{Name: p.Nickname, Id: p.ID, Interval: interval, Price: FormatPrice(p.Amount)}
		productList = append(productList, curProduct)
	}
	return productList
}
