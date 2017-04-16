package models

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
)

// TODO: Move these to somewhere secure
const (
	accessTokenSecret  = "ayylmao"
	refreshTokenSecret = "lmaoayy"
)

// GetAccessToken - Attempts to create an access token for the user
// TODO: Maybe just pass in roles string instead of *User
func GetAccessToken(u *User) (string, error) {
	// expiry time
	// TODO: Define constant in pkg/constants to define this expiry time
	tokenExpire := time.Now().Add(time.Minute * 30).Unix()

	// Set claims for the token
	claims := make(jwt.MapClaims)
	claims["exp"] = tokenExpire
	claims["iat"] = "multapply.io"
	claims["iss"] = time.Now().Unix()
	claims["roles"] = u.Roles

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
	tokenExpire := time.Now().Add(time.Minute * 525600).Unix()

	// Set claims for the token
	claims := make(jwt.MapClaims)
	claims["exp"] = tokenExpire
	claims["iat"] = "multapply.io"
	claims["iss"] = time.Now().Unix()
	claims["tid"] = tokenID

	// create and sign the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	var signedToken string
	signedToken, err = token.SignedString([]byte(refreshTokenSecret))

	return signedToken, err
}
