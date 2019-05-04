package stripe

import (
	"errors"
	"fmt"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
	"os"
)

type Customer struct {
	ID        string
	firstName string
	lastName  string
}

func FindCustomer(key string, email string) (*Customer, error) {
	stripe.Key = os.Getenv(key)
	list := customer.List(
		&stripe.CustomerListParams{
			Email: &email,
		})
	var customer *stripe.Customer
	count := 0
	for list.Next() {
		count = count + 1
		customer = list.Customer()
	}
	if count > 1 {
		return nil, errors.New("Multiple accounts with this email: " + email)
	} else if count == 0 {
		return nil, errors.New("No account with this email: " + email)
		// Return empty Option
	} else {
		// TODO Retrieve Subscriptions for Customers?
		//  I *think* that's helpful to let customers (And Joyce) know if they are Active members
		return &Customer{
			customer.ID,
			customer.Name,
			"SPLIT THIS UP INTO FIRST/LAST",
		}, nil
	}

}

/*
func CreateSource(key string, token string) {
	sourceParams := &stripe.SourceParams{
		Token: &token,
	}
	source.New(
		sourceParams
		)
}
*/

func CreateCustomer(key string, sourceToken string, email string, name string) string {
	stripe.Key = os.Getenv(key)

	params := &stripe.CustomerParams{
		Email: &email,
		Name:  &name,
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

func CreateTestCustomer(sourceToken string) string {
	return CreateCustomer("TEST_STRIPE_SECRET_KEY", sourceToken, "", "HARDCODED_NAME")
}
