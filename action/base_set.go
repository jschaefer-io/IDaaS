package action

import (
	"net/http"
)

// Basic Resource Action Set
type Set interface {
	Index(http.ResponseWriter, *http.Request)
	Create(http.ResponseWriter, *http.Request)
	Show(http.ResponseWriter, *http.Request, interface{})
	Delete(http.ResponseWriter, *http.Request, interface{})
	Update(http.ResponseWriter, *http.Request, interface{})
}

// Base Action set which can be
// embedded to fulfill the Set
// interface with only partial
// action support
type BaseActionSet struct{}

// Default index Route
// results in a 404 error
func (b BaseActionSet) Index(w http.ResponseWriter, r *http.Request) {
	Error404(w, r)
}

// Default index Route
// results in a 404 error
func (b BaseActionSet) Create(w http.ResponseWriter, r *http.Request) {
	Error404(w, r)
}

// Default show Route
// results in a 404 error
func (b BaseActionSet) Show(w http.ResponseWriter, r *http.Request, obj interface{}) {
	Error404(w, r)
}

// Default delete Route
// results in a 404 error
func (b BaseActionSet) Delete(w http.ResponseWriter, r *http.Request, obj interface{}) {
	Error404(w, r)
}

// Default update Route
// results in a 404 error
func (b BaseActionSet) Update(w http.ResponseWriter, r *http.Request, obj interface{}) {
	Error404(w, r)
}
