package resource

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jschaefer-io/IDaaS/action"
	"github.com/jschaefer-io/IDaaS/db"
	"github.com/jschaefer-io/IDaaS/reponse"
	"net/http"
	"regexp"
	"strconv"
)

type Resource struct {
	route string
	set   action.Set
}

func NewResource(route string, set action.Set) Resource {
	return Resource{route, set}
}

func (r Resource) getParam(value string, regex string) (string, error) {
	if !regexp.MustCompile(regex).Match([]byte(value)) {
		return value, errors.New("parameter does not match")
	}
	return value, nil
}

func (r Resource) findId(item interface{}, id string) error {
	return db.Get().Find(item, id).Error
}

func (r Resource) execute(c *gin.Context, fun func(*gin.Context, int)) {
	idString, err := r.getParam(c.Param("id"), "\\d+")
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

func (r Resource) Apply(e *gin.Engine) {
	e.GET(r.route, r.set.Index)
	e.POST(r.route, r.set.Create)
	e.GET(r.route+"/:id", func(c *gin.Context) { r.execute(c, r.set.Show) })
	e.PUT(r.route+"/:id", func(c *gin.Context) { r.execute(c, r.set.Update) })
	e.DELETE(r.route+"/:id", func(c *gin.Context) { r.execute(c, r.set.Delete) })
}
