package models

import (
	"testing"

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

// TODO: Test RefreshToken and other DB-reliant funcs
