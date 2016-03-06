package models

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"

	"github.com/eirka/eirka-libs/db"
	e "github.com/eirka/eirka-libs/errors"
)

func TestAddTagValidateInput(t *testing.T) {

	var err error

	tag := AddTagModel{
		Ib:    1,
		Tag:   0,
		Image: 1,
	}

	err = tag.ValidateInput()
	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err, e.ErrInvalidParam, "Error should match")
	}

}

func TestAddTagIsValid(t *testing.T) {

	tag := AddTagModel{
		Ib:    0,
		Tag:   1,
		Image: 1,
	}

	assert.False(t, tag.IsValid(), "Should be false")

}

func TestAddTagStatus(t *testing.T) {

	var err error

	mock, err := db.NewTestDb()
	assert.NoError(t, err, "An error was not expected")

	statusrows := sqlmock.NewRows([]string{"count"}).AddRow(1)
	mock.ExpectQuery(`SELECT count\(1\) FROM images`).WillReturnRows(statusrows)

	duperows := sqlmock.NewRows([]string{"count"}).AddRow(0)
	mock.ExpectQuery(`select count\(1\) from tagmap`).WillReturnRows(duperows)

	tag := AddTagModel{
		Ib:    1,
		Tag:   1,
		Image: 1,
	}

	err = tag.Status()
	assert.NoError(t, err, "An error was not expected")

}

func TestAddTagStatusNotFound(t *testing.T) {

	var err error

	mock, err := db.NewTestDb()
	assert.NoError(t, err, "An error was not expected")

	statusrows := sqlmock.NewRows([]string{"count"}).AddRow(0)
	mock.ExpectQuery(`SELECT count\(1\) FROM images`).WillReturnRows(statusrows)

	duperows := sqlmock.NewRows([]string{"count"}).AddRow(0)
	mock.ExpectQuery(`select count\(1\) from tagmap`).WillReturnRows(duperows)

	tag := AddTagModel{
		Ib:    1,
		Tag:   1,
		Image: 1,
	}

	err = tag.Status()
	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err, e.ErrNotFound, "Error should match")
	}

}

func TestAddTagStatusDuplicate(t *testing.T) {

	var err error

	mock, err := db.NewTestDb()
	assert.NoError(t, err, "An error was not expected")

	statusrows := sqlmock.NewRows([]string{"count"}).AddRow(1)
	mock.ExpectQuery(`SELECT count\(1\) FROM images`).WillReturnRows(statusrows)

	duperows := sqlmock.NewRows([]string{"count"}).AddRow(1)
	mock.ExpectQuery(`select count\(1\) from tagmap`).WillReturnRows(duperows)

	tag := AddTagModel{
		Ib:    1,
		Tag:   1,
		Image: 1,
	}

	err = tag.Status()
	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err, e.ErrDuplicateTag, "Error should match")
	}

}

func TestAddTagPost(t *testing.T) {

	var err error

	mock, err := db.NewTestDb()
	assert.NoError(t, err, "An error was not expected")

	mock.ExpectExec("INSERT into tagmap").
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	tag := AddTagModel{
		Ib:    1,
		Tag:   1,
		Image: 1,
	}

	err = tag.Post()
	assert.NoError(t, err, "An error was not expected")

}

func TestAddTagPostInvalid(t *testing.T) {

	var err error

	mock, err := db.NewTestDb()
	assert.NoError(t, err, "An error was not expected")

	mock.ExpectExec("INSERT into tagmap").
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	tag := AddTagModel{
		Ib:    0,
		Tag:   1,
		Image: 1,
	}

	err = tag.Post()
	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err, errors.New("AddTagModel is not valid"), "Error should match")
	}

}