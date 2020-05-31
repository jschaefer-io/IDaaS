package resource

import (
	"github.com/gin-gonic/gin"
	"github.com/jschaefer-io/IDaaS/action"
	"github.com/jschaefer-io/IDaaS/reponse"
	"net/http"
	"strconv"
)

// Route Resource type
type Resource struct {
	route string
	set   action.Set
}

// Creates a new Resource from a route and an action set
func NewResource(route string, set action.Set) Resource {
	return Resource{route, set}
}

// Executes the parameter validation on the show delete and update routes
func (r Resource) execute(c *gin.Context, fun func(*gin.Context, int)) {
	idString, err := GetParam(c.Param("id"), "\\d+")
	if err != nil {
		reponse.NewError(http.StatusUnprocessableEntity, "given id is not numeric").Apply(c)
		return
	}

	id, err := strconv.Atoi(idString)
	if err != nil {
		reponse.NewError(http.StatusNotAcceptable, "given id is not a valid int").Apply(c)
	}

	// Call route handler with the id
	fun(c, id)
}


// Applies the resource routes to the
// given routing engine
func (r Resource) Apply(e *gin.Engine) {
	e.GET(r.route, r.set.Index)
	e.POST(r.route, r.set.Create)
	e.GET(r.route+"/:id", func(c *gin.Context) { r.execute(c, r.set.Show) })
	e.PUT(r.route+"/:id", func(c *gin.Context) { r.execute(c, r.set.Update) })
	e.DELETE(r.route+"/:id", func(c *gin.Context) { r.execute(c, r.set.Delete) })
}
