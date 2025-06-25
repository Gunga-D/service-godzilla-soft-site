package subscribe

import "time"

var prices = map[string]int64{
	"5minute": 10000,
	"month":   25000,
	"year":    150000,
}

var durations = map[string]time.Duration{
	"5minute": 5 * time.Minute,
	"month":   time.Hour * 24 * 31,
	"year":    time.Hour * 24 * 31 * 12,
}

var durationNames = map[string]string{
	"5minute": "5 минут",
	"month":   "месяц",
	"year":    "год",
}
