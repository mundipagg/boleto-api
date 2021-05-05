package models

import "strings"

const (
	MaxInterestRate                = 0.020
	MinDaysToStartChargingInterest = 0
)

//Interest Representa as informações sobre Multa
type Interest struct {
	Type string  `json:"type,omitempty"`
	Days int     `json:"days,omitempty"`
	Rate float64 `json:"rate,omitempty"`
}

//IsValid Valida as regras de negócio sobre Juros
func (i *Interest) IsValid(amountInCents int) bool {
	if i.Days < MinDaysToStartChargingInterest {
		return false
	}

	if !i.isValidRate(amountInCents) {
		return false
	}

	return true
}

func (i *Interest) isValidRate(amountInCents int) bool {
	pcond := i.isValidPercentualRate()
	ncond := i.isValidNominalRate(amountInCents)

	return (!pcond && ncond) || (pcond && !ncond)
}

func (i *Interest) isValidPercentualRate() bool {
	if strings.ToLower(i.Type) != Percentual {
		return false
	}

	if i.Rate < 0 || i.Rate > MaxInterestRate {
		return false
	}

	return true
}

func (i *Interest) isValidNominalRate(amountInCents int) bool {
	if strings.ToLower(i.Type) != Nominal {
		return false
	}

	if amountInCents <= 0 {
		return false
	}

	rate := float64(int(i.Rate)) / float64(amountInCents)

	return rate <= MaxInterestRate
}
