package transaction

// DepositCommand requests to deposit an amount to an account from an ATM
type DepositCommand struct {
	AccountTo string  `json:"to"`     // The account which receives the money
	Amount    float64 `json:"amount"` // The amount
}

// WithdrawCommand requests to withdraw an amount from an account via an ATM
type WithdrawCommand struct {
	AccountFrom string  `json:"from"`   // The account which sends the money
	Amount      float64 `json:"amount"` // The amount
}

// TransferCommand requests to transfer an amount between two accounts
type TransferCommand struct {
	AccountFrom string  `json:"from"`   // The account which sends the money
	AccountTo   string  `json:"to"`     // The account which receives the money
	Amount      float64 `json:"amount"` // The amount
}

// Command represents a command
type Command interface{}
