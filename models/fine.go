package models

import "strings"

const (
	MaxFineRate                = 0.010
	MinDaysToStartChargingFine = 0
)

//Fine Representa as informações sobre Juros
type Fine struct {
	Type string  `json:"type,omitempty"`
	Days int     `json:"days,omitempty"`
	Rate float64 `json:"rate,omitempty"`
}

//IsValid Valida as regras de negócio sobre Juros
func (f *Fine) IsValid(amountInCents int) bool {
	if f.Days < MinDaysToStartChargingFine {
		return false
	}

	if !f.isValidRate(amountInCents) {
		return false
	}

	return true
}

func (f *Fine) isValidRate(amountInCents int) bool {
	pcond := f.isValidPercentualRate()
	ncond := f.isValidNominalRate(amountInCents)

	return (!pcond && ncond) || (pcond && !ncond)
}

func (f *Fine) isValidPercentualRate() bool {
	if strings.ToLower(f.Type) != Percentual {
		return false
	}

	if f.Rate < 0 || f.Rate > MaxFineRate {
		return false
	}

	return true
}

func (f *Fine) isValidNominalRate(amountInCents int) bool {
	if strings.ToLower(f.Type) != Nominal {
		return false
	}

	if amountInCents <= 0 {
		return false
	}

	rate := float64(int(f.Rate)) / float64(amountInCents)

	return rate <= MaxFineRate
}
