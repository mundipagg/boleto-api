package models

//Rules Define regras de pagamento e baixa do título
type Rules struct {
	AcceptDivergentAmount bool `json:"acceptDivergentAmount,omitempty"`
	MaxDaysToPayPastDue   int  `json:"maxDaysToPayPastDue,omitempty"`
}
