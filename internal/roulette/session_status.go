package roulette

type SessionStatus string

const (
	WaitForPayment SessionStatus = "wait-for-payment"
	ReadyToRoll    SessionStatus = "ready-to-roll"
	OutOfGames     SessionStatus = "out-of-games"
)
