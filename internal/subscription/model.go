package subscription

type UserSubscription struct {
	Status    string `db:"status"`
	ExpiredAt int64  `db:"expired_at"`
}

type PaidSubscription struct {
	UserID      int64  `db:"user_id"`
	CreatedAt   int64  `db:"created_at"`
	ExpiredAt   int64  `db:"expired_at"`
	RebillID    string `db:"rebill_id"`
	NeedProlong bool   `db:"need_prolong"`
}

type SubscriptionProduct struct {
	Login    string `db:"login"`
	Password string `db:"password"`
}
