package models

const (
	Percentual = "percentual"
	Nominal    = "nominal"
)

type Fees struct {
	Fine     Fine     `json:"fine,omitempty"`
	Interest Interest `json:"interest,omitempty"`
}
