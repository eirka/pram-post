package models

import (
	"database/sql"

	"github.com/techjanitor/pram-libs/config"
	"github.com/techjanitor/pram-libs/db"
	e "github.com/techjanitor/pram-libs/errors"
	"github.com/techjanitor/pram-libs/validate"
)

// loginmodel contains user name and password
type LoginModel struct {
	Name     string
	Password string
	Id       uint
	Hash     []byte
}

// Validate will check the provided name and password
func (r *LoginModel) Validate() (err error) {

	// Validate name input
	name := validate.Validate{Input: r.Name, Max: config.Settings.Limits.NameMaxLength, Min: config.Settings.Limits.NameMinLength}
	if name.IsEmpty() {
		return e.ErrNameEmpty
	} else if name.MinLength() {
		return e.ErrNameShort
	} else if name.MaxLength() {
		return e.ErrNameLong
	} else if !name.IsUsername() {
		return e.ErrNameAlphaNum
	}

	// Validate password input
	password := validate.Validate{Input: r.Password, Max: config.Settings.Limits.NameMaxLength, Min: config.Settings.Limits.NameMinLength}
	if password.IsEmpty() {
		return e.ErrPasswordEmpty
	} else if password.MinLength() {
		return e.ErrPasswordShort
	} else if password.MaxLength() {
		return e.ErrPasswordLong
	}

	return

}

// query user info from database
func (r *LoginModel) Query() (err error) {

	// Get Database handle
	dbase, err := db.GetDb()
	if err != nil {
		return e.ErrInternalError
	}

	// get hashed password from database
	err = dbase.QueryRow("select user_id, user_password from users where user_name = ?", r.Name).Scan(&r.Id, &r.Hash)
	if err == sql.ErrNoRows {
		return e.ErrInvalidUser
	} else if err != nil {
		return e.ErrInternalError
	}

	return

}
