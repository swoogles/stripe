package stripe

import "github.com/leekchan/accounting"

func FormatPrice(amount int64) string {
	ac := accounting.Accounting{Symbol: "$", Precision: 2}
	return ac.FormatMoney(float64(amount) / 100.0)
}
