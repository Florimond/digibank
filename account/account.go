package account

// Account represents state of an account
type Account struct {
	ID       string  `json:"id"`       // The ID of the account
	Customer string  `json:"customer"` // The customer owning the account
	Version  uint    `json:"version"`  // The version of the amount
	Amount   float64 `json:"amount"`   // The amount on the account
}
