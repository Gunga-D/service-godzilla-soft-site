package order

type PaidOrder struct {
	ID        string `db:"id"`
	Email     string `db:"email"`
	CodeValue string `db:"code_value"`
	ItemSlip  string `db:"item_slip"`
	Amount    int64  `db:"amount"`
}

type UserOrder struct {
	ID        string  `db:"id"`
	CodeValue string  `db:"code_value"`
	ItemSlip  *string `db:"item_slip"`
	ItemName  *string `db:"item_name"`
	Amount    int64   `db:"amount"`
	Status    string  `db:"status"`
}
