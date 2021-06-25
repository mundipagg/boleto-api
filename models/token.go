package models

import "time"

// Token represents a token used to access external resources
// The field CreatedAt is used by Mongo to control token TTL
type Token struct {
	Origin      string    `json:"origin,omitempty"`
	AccessToken string    `json:"accesstoken,omitempty"`
	CreatedAt   time.Time `json:"createdat"`
}

//Token Creates a Token instance
func NewToken(origin, token string) Token {
	return Token{
		Origin:      origin,
		AccessToken: token,
		CreatedAt:   time.Now(),
	}
}
