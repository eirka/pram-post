package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"

	e "github.com/techjanitor/pram-post/errors"
	"github.com/techjanitor/pram-post/models"
	u "github.com/techjanitor/pram-post/utils"
)

// Input from change email form
type emailForm struct {
	Ib    uint   `json:"ib" binding:"required"`
	Email string `json:"email" binding:"required"`
}

// EmailController handles initial registration
func EmailController(c *gin.Context) {
	var err error
	var ef emailForm

	// get userdata from session middleware
	userdata := c.MustGet("userdata").(u.User)

	err = c.Bind(&ef)
	if err != nil {
		c.JSON(e.ErrorMessage(e.ErrInvalidParam))
		c.Error(err)
		return
	}

	// Set parameters to EmailModel
	m := models.EmailModel{
		Uid:   userdata.Id,
		Email: ef.Email,
	}

	// Validate input
	err = m.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_message": err.Error()})
		c.Error(err)
		return
	}

	// update password
	err = m.Update()
	if err != nil {
		c.JSON(e.ErrorMessage(e.ErrInternalError))
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success_message": u.AuditEmailUpdate})

	audit := u.Audit{
		User:   userdata.Id,
		Ib:     ef.Ib,
		Ip:     c.ClientIP(),
		Action: u.AuditEmailUpdate,
		Info:   userdata.Name,
	}

	err = audit.Submit()
	if err != nil {
		c.Error(err)
	}

	return

}