package models

import (
	//"errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"

	"github.com/eirka/eirka-libs/db"
	e "github.com/eirka/eirka-libs/errors"
)

func TestEmailIsValid(t *testing.T) {

	email := EmailModel{
		Uid:   0,
		Name:  "test",
		Email: "cool@test.com",
	}

	assert.False(t, email.IsValid(), "Should be false")

}

func TestEmailValidate(t *testing.T) {

	var err error

	mock, err := db.NewTestDb()
	assert.NoError(t, err, "An error was not expected")

	rows := sqlmock.NewRows([]string{"name", "email"}).AddRow("test", "old@test.com")
	mock.ExpectQuery(`SELECT user_name,user_email FROM users WHERE user_id`).WillReturnRows(rows)

	email := EmailModel{
		Uid:   1,
		Email: "cool@test.com",
	}

	err = email.Validate()
	if assert.NoError(t, err, "An error was not expected") {
		assert.Equal(t, email.Name, "test", "Should match")
		assert.Equal(t, email.CurrentEmail, "old@test.com", "Should match")
	}
}

func TestEmailValidateBadEmails(t *testing.T) {

	var err error

	first := EmailModel{
		Email: "notanemail",
	}

	err = first.Validate()
	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err, e.ErrInvalidEmail, "Error should match")
	}

	second := EmailModel{
		Email: "not@anemail",
	}

	err = second.Validate()
	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err, e.ErrInvalidEmail, "Error should match")
	}
}
