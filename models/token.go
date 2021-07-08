package models

import "time"

// Token represents a token used to access external resources
// The field CreatedAt is used by Mongo to control token TTL
type Token struct {
	ClientID    string    `json:"clientid,omitempty"`
	IssuerBank  string    `json:"issuerbank,omitempty"`
	AccessToken string    `json:"accesstoken,omitempty"`
	CreatedAt   time.Time `json:"createdat"`
}

//Token Creates a Token instance
func NewToken(clientID, issuerbank, token string) Token {
	return Token{
		ClientID:    clientID,
		IssuerBank:  issuerbank,
		AccessToken: token,
		CreatedAt:   time.Now(),
	}
}
