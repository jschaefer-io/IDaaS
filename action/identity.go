package action

import (
	"github.com/gin-gonic/gin"
	"github.com/jschaefer-io/IDaaS/crypto"
	"github.com/jschaefer-io/IDaaS/db"
	"github.com/jschaefer-io/IDaaS/model"
	"github.com/jschaefer-io/IDaaS/reponse"
	"net/http"
)

// Identity Action Set
type Identity struct {
	BaseActionSet
}

// Main index route
func (i Identity) Index(c *gin.Context) {
	var ids []model.Identity
	db.Get().Find(&ids)
	c.JSON(http.StatusOK, ids)
}

// Main creat route
func (i Identity) Create(c *gin.Context) {
	var json model.Identity
	if err := c.ShouldBindJSON(&json); err != nil {
		model.ValidationError(err).Apply(c)
		return
	}

	id := model.NewIdentity(json.Email, json.Password)
	_ = db.Get().Create(&id)

	//fmt.Println(res.Error.Error())
	c.JSON(http.StatusOK, id)
}

// Main show single route
func (i Identity) Show(c *gin.Context, id int) {
	identity, err := model.Identity{}.Find(id)
	if err != nil {
		reponse.NewError(http.StatusUnprocessableEntity, err).Apply(c)
		return
	}
	reponse.NewResponse(http.StatusOK, identity).Apply(c)
}

// Main update route
func (i Identity) Update(c *gin.Context, id int) {
	identity, err := model.Identity{}.Find(id)
	if err != nil {
		reponse.NewError(http.StatusUnprocessableEntity, err).Apply(c)
		return
	}

	var json model.PasswordForm
	if err := c.ShouldBindJSON(&json); err != nil {
		model.ValidationError(err).Apply(c)
		return
	}

	pwd := crypto.NewPassword(json.Password)
	identity.Password = pwd.String()
	db.Get().Save(&identity)
	reponse.NewResponse(http.StatusOK, identity).Apply(c)
}

// Main delete route
func (i Identity) Delete(c *gin.Context, id int) {
	identity, err := model.Identity{}.Find(id)
	if err != nil {
		reponse.NewError(http.StatusUnprocessableEntity, err).Apply(c)
		return
	}

	db.Get().Delete(&identity)
	reponse.NewResponse(http.StatusOK, identity).Apply(c)
}
