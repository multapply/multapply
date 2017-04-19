package models

import (
	"testing"

	"github.com/multapply/multapply/util/userRoles"
	"github.com/stretchr/testify/assert"
)

var n = &NewUser{
	Username:        "  testerino  ",
	FirstName:       " Bob   ",
	LastName:        "\tJenkins ",
	Email:           "\ntest@testemail.com\t",
	Password:        "\t\thunter2 ",
	ConfirmPassword: "     hunter2\t",
	PasswordHash:    "asdf",
}

func TestTrim(t *testing.T) {
	n.Trim()

	assert.Equal(t, "testerino", n.Username, "Username should be trimmed")
	assert.Equal(t, "Bob", n.FirstName, "First name should be trimmed")
	assert.Equal(t, "Jenkins", n.LastName, "Last name should be trimmed")
	assert.Equal(t, "test@testemail.com", n.Email, "Email should be trimmed")
	assert.Equal(t, "hunter2", n.Password, "Password should be trimmed")
	assert.Equal(t, "hunter2", n.ConfirmPassword, "Password confirmation should be trimmed")
	assert.Equal(t, "asdf", n.PasswordHash, "Password hash should be unchanged")
}

func TestCreateUser(t *testing.T) {
	u := CreateUser(n)

	assert.Empty(t, u.UserID, "UserID should be default value (0)")
	assert.Equal(t, n.Username, u.Username, "Username should be same as NewUser's")
	assert.Equal(t, n.FirstName, u.FirstName, "First name should be same as NewUser's")
	assert.Equal(t, n.LastName, u.LastName, "Last name should be same as NewUser's")
	assert.Equal(t, n.Email, u.Email, "Email should be same as NewUser's")
	assert.Equal(t, n.PasswordHash, u.PasswordHash, "Password hash should be same as NewUser's")
	assert.Equal(t, userRoles.BasicUser, u.Roles, "Roles should be equal to \"USER\"")
}
