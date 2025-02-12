package order

type PaidOrder struct {
	ID        string `db:"id"`
	Email     string `db:"email"`
	CodeValue string `db:"code_value"`
	ItemSlip  string `db:"item_slip"`
}
