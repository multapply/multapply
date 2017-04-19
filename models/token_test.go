package models

import (
	"errors"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/multapply/multapply/util/userRoles"
	"github.com/stretchr/testify/assert"
)

var u = &User{
	UserID:       42,
	Username:     "testerino",
	FirstName:    "Larry",
	LastName:     "Skywalker",
	Email:        "star@wars.com",
	PasswordHash: "lololol",
	Roles:        userRoles.Admin,
}

func TestGetAccessToken(t *testing.T) {
	_, err := GetAccessToken(u)

	assert.Nil(t, err, "There should be no error creating the access token")
}

func TestAccessTokenCorrect(t *testing.T) {
	accessToken, err := GetAccessToken(u)
	assert.Nil(t, err, "There should be no error creating the access token")

	parsedToken, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Token signature invalid")
		}
		return []byte("ayylmao"), nil
	})
	assert.Nil(t, err, "There should be no error parsing the access token")

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	assert.True(t, ok, "There should be no error parsing the access token's claims")
	assert.Equal(t, "multapply.io", claims["iss"], "Claims' issuer should be correct")
	assert.Equal(t, float64(1800), claims["exp"].(float64)-claims["iat"].(float64),
		"Token's expiration should be 30 minutes exactly after issue date")
	assert.Equal(t, userRoles.Admin, claims["roles"], "Claims' roles should be correct")
	assert.Equal(t, float64(42), claims["uid"], "Claims' uid (user id) should be correct")
}

// TODO: Test RefreshToken and other DB-reliant funcs
