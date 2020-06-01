package resource

import (
	"github.com/gorilla/mux"
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
func (r Resource) execute(writer http.ResponseWriter, request *http.Request, fun func(http.ResponseWriter, *http.Request, int)) {
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		reponse.NewError(http.StatusUnprocessableEntity, "given id is not numeric").Apply(writer)
		return
	}

	// Call route handler with the id
	fun(writer, request, id)
}

// Applies the resource routes to the given router
func (r Resource) Apply(e *mux.Router) {
	e.HandleFunc(r.route, r.set.Index).Methods("GET")
	e.HandleFunc(r.route, r.set.Create).Methods("POST")

	idSuffix := "/{id:[0-9]+}"
	handleIdRoutes := func(handler func(http.ResponseWriter, *http.Request, int)) func(http.ResponseWriter, *http.Request) {
		return func(w http.ResponseWriter, req *http.Request) {
			r.execute(w, req, handler)
		}
	}

	e.HandleFunc(r.route+idSuffix, handleIdRoutes(r.set.Show)).Methods("GET")
	e.HandleFunc(r.route+idSuffix, handleIdRoutes(r.set.Update)).Methods("PUT", "PATCH")
	e.HandleFunc(r.route+idSuffix, handleIdRoutes(r.set.Delete)).Methods("DELETE")
}
