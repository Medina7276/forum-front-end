package dto

import model "forum/model"

type DefaultClaims struct {
	// unix time
	IssuedAt  int64 `json:"iat"`
	ExpiresAt int64 `json:"exp"`
}

type UserWithClaims struct {
	User   model.User
	Claims DefaultClaims
}
