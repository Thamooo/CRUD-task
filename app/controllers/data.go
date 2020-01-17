package controllers

import (
	"aahframe.work"

	"test-task/app/models"
)

// AppController struct application controller
type DataController struct {
	*aah.Context
}

// Index method is application's home page.
func (c *DataController) Index() {
	data := aah.Data{
		"Greet": models.Greet{
			Message: "Welcome to aah framework - Web Application",
		},
	}

	c.Reply().Ok().HTML(data)
}
