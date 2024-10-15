package middleware

import (
	"context"
	"errors"
)

// CustomClaimsExample contains custom data we want from the token.
type UserClaims struct {
	Name         string `json:"name"`
	Username     string `json:"username"`
	ShouldReject bool   `json:"shouldReject,omitempty"`
}

// Validate errors out if `ShouldReject` is true.
func (uc *UserClaims) Validate(c context.Context) error {
	if uc.ShouldReject {
		return errors.New("should reject was set to true")
	}
	return nil
}
