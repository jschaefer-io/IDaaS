package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jschaefer-io/IDaaS/db"
	"github.com/jschaefer-io/IDaaS/model"
	"github.com/jschaefer-io/IDaaS/reponse"
	"github.com/jschaefer-io/IDaaS/resource"
	"net/http"
)

type Test struct {
	Route    string
	Resource resource.Resource
}

func (t *Test) execute(c *gin.Context, fun func(*gin.Context, interface{})) {
	id, err := resource.GetParam(c, "id", "\\d+")
	if err != nil {
		reponse.NewError(http.StatusUnprocessableEntity, "given id is not numeric").Apply(c)
		return
	}
	var identity model.Identity
	if err = db.Get().Find(&identity, id).Error; err != nil {
		reponse.NewError(http.StatusUnprocessableEntity, err.Error()).Apply(c)
		return
	}
	fun(c, identity)
}

func (t *Test) apply(r *gin.Engine) {
	r.GET(t.Route, t.Resource.Index)
	r.POST(t.Route, t.Resource.Create)
	r.GET(t.Route+"/:id", func(c *gin.Context) { t.execute(c, t.Resource.Show) })
	r.PUT(t.Route+"/:id", func(c *gin.Context) { t.execute(c, t.Resource.Update) })
	r.DELETE(t.Route+"/:id", func(c *gin.Context) { t.execute(c, t.Resource.Delete) })
}


type A struct {
	Email string
}

func main() {

	// https://github.com/dgrijalva/jwt-go
	// https://golang.org/pkg/net/smtp/


	db.Get().AutoMigrate(&model.Identity{})

	r := gin.Default()

	t := Test{
		Route:    "/",
		Resource: new(resource.Identity),
	}

	t.apply(r)

	//l:=resource.NewLorem(resource.BaseResource{}, model.Identity{})
	//l.Add(r)
	//resource.Add(r, resource.BaseResource{})
	//r.GET("/", action.IdentityIndex)
	//r.POST("/", action.IdentityCreate)
	//res := new(resource.Identity)
	//r.GET("/:id", res.Show)
	//r.PUT("/:id", action.IdentityUpdate)
	//r.DELETE("/:id", action.IdentityDelete)
	r.Run()
}
