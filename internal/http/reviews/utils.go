package reviews

import (
	"math"
	"time"
)

func Round(x float64, prec int) float64 {
	var rounder float64
	pow := math.Pow(10, float64(prec))
	intermed := x * pow
	_, frac := math.Modf(intermed)
	if frac >= 0.5 {
		rounder = math.Ceil(intermed)
	} else {
		rounder = math.Floor(intermed)
	}

	return rounder / pow
}

func RussianMonth(t time.Time) (month string) {
	switch t.Month() {
	case time.January:
		month = "Января"
	case time.February:
		month = "Февраля"
	case time.March:
		month = "Марта"
	case time.April:
		month = "Апреля"
	case time.May:
		month = "Мая"
	case time.June:
		month = "Июня"
	case time.July:
		month = "Июля"
	case time.August:
		month = "Августа"
	case time.September:
		month = "Сентября"
	case time.October:
		month = "Октября"
	case time.November:
		month = "Ноября"
	case time.December:
		month = "Декабря"
	}
	return
}
