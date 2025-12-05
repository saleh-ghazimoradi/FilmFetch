package domain

import "time"

const (
	ScopeActivation = "activation"
)

type Token struct {
	Plaintext string    `json:"plaintext"`
	Hash      []byte    `json:"hash"`
	UserId    int64     `json:"user_id"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"scope"`
}
