package calculate

import (
	"strconv"
)

type Interface interface {
	getAmount() int
}

type AmountStruct struct {
	Amount int
}

func (a AmountStruct) getAmount() int {
	return a.Amount
}

func Calculate(a ...any) (res int) {
	for _, val := range a {
		switch val.(type) {
		case int:
			res += val.(int)
		case string:
			amount, _ := strconv.Atoi(val.(string))
			res += amount
		case Interface:
			res += val.(Interface).getAmount()
		default:
			panic("Type is not supported")
		}
	}

	return
}
