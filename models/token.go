package models

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
)

// TODO: Move these to somewhere secure
const (
	accessTokenSecret  = "ayylmao"
	refreshTokenSecret = "lmaoayy"
)

// AccessTokenClaims - Custom Claims type jwt access tokens
type AccessTokenClaims struct {
	Issuer    string `json:"iss,omitempty"`
	ExpiresAt int64  `json:"exp,omitempty"`
	IssuedAt  int64  `json:"iat,omitempty"`
	Roles     string `json:"roles,omitempty"`
	UserID    int    `json:"uid,omitempty"`
}

// RefreshTokenClaims - Custom Claims type for jwt refresh tokens
type RefreshTokenClaims struct {
	Issuer    string `json:"iss,omitempty"`
	ExpiresAt int64  `json:"exp,omitempty"`
	IssuedAt  int64  `json:"iat,omitempty"`
	TokenID   int    `json:"tid,omitempty"`
}

// Valid - AccessTokenClaims needs a Valid() method to be a valid jwt.Claims
func (c AccessTokenClaims) Valid() error {
	if c.Issuer != "multapply.io" {
		return errors.New("Invalid issuer")
	} else if c.ExpiresAt > time.Now().Unix() {
		return errors.New("Token expired")
	} else if (c.ExpiresAt - c.IssuedAt) != 1800 { // TODO: Put this in constants
		return errors.New("Invalid issue date")
	}
	return nil
}

// Valid - RefreshTokenClaims needs a Valid() method to be a valid jwt.Claims
func (c RefreshTokenClaims) Valid() error {
	if c.Issuer != "multapply.io" {
		return errors.New("Invalid issuer")
	} else if c.ExpiresAt > time.Now().Unix() {
		return errors.New("Token expired")
	} else if (c.ExpiresAt - c.IssuedAt) != 31536000 { // TODO: Put this in constants
		return errors.New("Invalid issue date")
	}
	return nil
}

// GetAccessToken - Attempts to create an access token for the user
// TODO: Maybe just pass in roles string instead of *User
func GetAccessToken(u *User) (string, error) {
	// expiry time
	// TODO: Define constant in pkg/constants to define this expiry time
	now := time.Now()
	// Set claims for the token
	claims := new(AccessTokenClaims)
	claims.Issuer = "multapply.io"
	claims.ExpiresAt = now.Add(time.Minute * 30).Unix()
	claims.IssuedAt = now.Unix()
	claims.Roles = u.Roles
	claims.UserID = u.UserID

	// create and sign the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(accessTokenSecret))

	return signedToken, err
}

// GetRefreshToken - Attempts to create a refresh token for the user
// TODO: Maybe just pass in user's username instead of entire *User
func GetRefreshToken(db *sqlx.DB, u *User) (string, error) {
	// First attempt to insert new token into refresh_tokens table
	tokenID := 0
	err := db.QueryRow(`INSERT INTO refresh_tokens (username) VALUES ($1) RETURNING token_id`, u.Username).Scan(&tokenID)
	if err != nil {
		return "", err
	}

	// expiry time
	// TODO: Use constant defined in pkg/constants to define this expiry time
	now := time.Now()

	// Set claims for the token
	claims := new(RefreshTokenClaims)
	claims.ExpiresAt = now.Add(time.Minute * 525600).Unix()
	claims.Issuer = "multapply.io"
	claims.IssuedAt = now.Unix()
	claims.TokenID = tokenID

	// create and sign the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	var signedToken string
	signedToken, err = token.SignedString([]byte(refreshTokenSecret))

	return signedToken, err
}

// GetUserByTokenID - Get User associated with given token_id
func GetUserByTokenID(db *sqlx.DB, tid int) (*User, error) {
	var username string
	err := db.Get(&username, "SELECT username FROM refresh_tokens WHERE token_id=$1 LIMIT 1", tid)
	if err != nil {
		return nil, err
	}

	u, err := GetUserByUsername(db, username)
	return u, err
}
