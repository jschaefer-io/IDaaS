package action

import (
	"github.com/gin-gonic/gin"
	"github.com/jschaefer-io/IDaaS/crypto"
	"github.com/jschaefer-io/IDaaS/db"
	"github.com/jschaefer-io/IDaaS/model"
	"github.com/jschaefer-io/IDaaS/reponse"
	"net/http"
)

type UpdatePasswordForm struct {
	Password string `json:"password"`
}

type AddIdentityForm struct {
	UpdatePasswordForm
	Email string `json:"email" binding:"required,email"`
}

type Identity struct {
	BaseActionSet
}

func (i Identity) requireIdentity(id int, c *gin.Context) (model.Identity, error) {
	identity := model.Identity{}
	if err := db.Get().Find(&identity, id).Error; err != nil {
		reponse.NewError(http.StatusUnprocessableEntity, err).Apply(c)
		return identity, err
	}
	return identity, nil
}

func (i Identity) Index(c *gin.Context) {
	var ids []model.Identity
	db.Get().Find(&ids)
	c.JSON(http.StatusOK, ids)
}

func (i Identity) Create(c *gin.Context) {
	var json AddIdentityForm
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := model.NewIdentity(json.Email, json.Password)
	_ = db.Get().Create(&id)

	//fmt.Println(res.Error.Error())
	c.JSON(http.StatusOK, id)
}

func (i Identity) Show(c *gin.Context, id int) {
	identity, err := i.requireIdentity(id, c)
	if err != nil {
		return
	}
	reponse.NewResponse(http.StatusOK, identity).Apply(c)
}

func (i Identity) Update(c *gin.Context, id int) {
	identity, err := i.requireIdentity(id, c)
	if err != nil {
		return
	}

	var json UpdatePasswordForm
	if err := c.ShouldBindJSON(&json); err != nil {
		reponse.NewError(http.StatusUnprocessableEntity, err.Error())
		return
	}

	identity.Password = crypto.NewPassword(json.Password)
	db.Get().Save(&identity)
	reponse.NewResponse(http.StatusOK, identity).Apply(c)
}

func (i Identity) Delete(c *gin.Context, id int) {
	identity, err := i.requireIdentity(id, c)
	if err != nil {
		return
	}
	db.Get().Delete(&identity)
	reponse.NewResponse(http.StatusOK, identity).Apply(c)
}
