package controllers

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eirka/eirka-libs/config"
	"github.com/eirka/eirka-libs/db"
	"github.com/eirka/eirka-libs/user"
)

func init() {

	// Database connection settings
	dbase := db.Database{

		User:           local.Settings.Database.User,
		Password:       local.Settings.Database.Password,
		Proto:          local.Settings.Database.Proto,
		Host:           local.Settings.Database.Host,
		Database:       local.Settings.Database.Database,
		MaxIdle:        local.Settings.Database.MaxIdle,
		MaxConnections: local.Settings.Database.MaxConnections,
	}

	// Set up DB connection
	dbase.NewDb()

	// Get limits and stuff from database
	config.GetDatabaseSettings()

	user.Secret = "secret"
}

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func performJwtFormRequest(r http.Handler, method, path, token string, body io.Reader) *httptest.ResponseRecorder {

	req, _ := http.NewRequest(method, path, &b)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestEmailController(t *testing.T) {

	var err error

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.Use(user.Auth(true))

	router.POST("/email", c.EmailController)

	first := performRequest(router, "POST", "/email")

	assert.Equal(t, first.Code, 401, "HTTP request code should match")

	u := user.DefaultUser()
	u.SetId(2)
	u.SetAuthenticated()

	u.hash, err = user.HashPassword("testpassword")
	if assert.NoError(t, err, "An error was not expected") {
		assert.NotNil(t, user.hash, "password should be returned")
	}

	token, err := u.CreateToken()
	if assert.NoError(t, err, "An error was not expected") {
		assert.NotEmpty(t, badtoken, "token should be returned")
	}

	var b bytes.Buffer

	mw := multipart.NewWriter(&b)
	mw.WriteField("ib", 1)
	mw.WriteField("email", "test@cool.com")
	mw.Close()

	second := performJwtFormRequest(router, "POST", "/email", token, b)

}
