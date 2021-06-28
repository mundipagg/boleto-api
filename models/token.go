package models

import "time"

// Token represents a token used to access external resources
// The field CreatedAt is used by Mongo to control token TTL
type Token struct {
	ClientID    string    `json:"clientid,omitempty"`
	Origin      string    `json:"origin,omitempty"`
	AccessToken string    `json:"accesstoken,omitempty"`
	CreatedAt   time.Time `json:"createdat"`
}

//Token Creates a Token instance
func NewToken(clientID, origin, token string) Token {
	return Token{
		ClientID:    clientID,
		Origin:      origin,
		AccessToken: token,
		CreatedAt:   time.Now(),
	}
}
