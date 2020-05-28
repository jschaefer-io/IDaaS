package resource

import (
	"errors"
	"github.com/gin-gonic/gin"
	"regexp"
)

func GetParam(c *gin.Context, param string, regex string) (string, error) {
	value := c.Param(param)
	if !regexp.MustCompile(regex).Match([]byte(value)) {
		return value, errors.New("parameter does not match")
	}
	return value, nil
}


//
//
//type Lorem struct {
//	resource Resource
//	model    interface{}
//}
//
//func NewLorem(r Resource, m interface{}) Lorem {
//	return Lorem{r, m}
//}
//
//func (l *Lorem) Add(engine *gin.Engine) {
//	engine.GET("/", func(context *gin.Context) {
//		l.resource.Index(context)
//	})
//
//	engine.POST("/", func(context *gin.Context) {
//		l.resource.Create(context)
//	})
//
//	engine.GET("/:id", func(context *gin.Context) {
//		l.dispatchWithParam("id", context, l.resource.Show)
//	})
//
//	engine.DELETE("/:id", func(context *gin.Context) {
//		l.dispatchWithParam("id", context, l.resource.Delete)
//	})
//
//	engine.PUT("/:id", func(context *gin.Context) {
//		l.dispatchWithParam("id", context, l.resource.Update)
//	})
//}
//
//func (l *Lorem) dispatchWithParam(param string, context *gin.Context, callback func(c *gin.Context, model interface{})) {
//	m, c := l.checkParam(context, "id")
//	if c != nil {
//		c.Apply(context)
//
//		return
//	}
//	callback(context, m)
//}
//
//func (l Lorem) checkParam(c *gin.Context, param string) (interface{}, reponse.Response) {
//	id := c.Param(param)
//	if !regexp.MustCompile("\\d+").Match([]byte(id)) {
//		return nil, reponse.NewError(http.StatusUnprocessableEntity, "given id is not numeric")
//	}
//	a := l.model
//	if err := db.Get().Find(&a, id).Error; err != nil {
//		return nil, reponse.NewError(http.StatusNotFound, err.Error())
//	}
//	return a, nil
//}
