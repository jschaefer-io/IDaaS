package resource

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
	BaseResource
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

func (i Identity) Show(c *gin.Context, m interface{}) {
	identity := m.(model.Identity)
	reponse.NewResponse(http.StatusOK, identity).Apply(c)
}

func (i Identity) Update(c *gin.Context, m interface{}) {
	identity := m.(model.Identity)

	var json UpdatePasswordForm
	if err := c.ShouldBindJSON(&json); err != nil {
		reponse.NewError(http.StatusUnprocessableEntity, err.Error())
		return
	}

	identity.Password = crypto.NewPassword(json.Password)
	db.Get().Save(&identity)
	reponse.NewResponse(http.StatusOK, identity).Apply(c)
}

func (i Identity) Delete(c *gin.Context, m interface{}) {
	identity := m.(model.Identity)
	db.Get().Delete(&identity)
	reponse.NewResponse(http.StatusOK, identity).Apply(c)
}
