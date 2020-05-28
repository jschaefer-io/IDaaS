package action

import (
	"github.com/gin-gonic/gin"
	"github.com/jschaefer-io/IDaaS/crypto"
	"github.com/jschaefer-io/IDaaS/db"
	"github.com/jschaefer-io/IDaaS/model"
	"github.com/jschaefer-io/IDaaS/reponse"
	"net/http"
	"regexp"
)

type UpdatePasswordForm struct {
	Password string `json:"password"`
}

type AddIdentityForm struct {
	UpdatePasswordForm
	Email string `json:"email" binding:"required,email"`
}

func IdentityIndex(c *gin.Context) {
	var ids []model.Identity
	db.Get().Find(&ids)
	c.JSON(http.StatusOK, ids)
}

func IdentityCreate(c *gin.Context) {
	var json AddIdentityForm
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := model.NewIdentity(json.Email, json.Password)
	db.Get().Create(id)
	c.JSON(http.StatusOK, id)
}

func IdentityShow(c *gin.Context) {
	id := c.Param("id")
	if !regexp.MustCompile("\\d+").Match([]byte(id)) {
		reponse.NewError(http.StatusUnprocessableEntity, "given id is not numeric").Apply(c)
		return
	}
	var identity model.Identity
	if err := db.Get().Find(&identity, id).Error; err != nil {
		reponse.NewError(http.StatusNotFound, err.Error()).Apply(c)
		return
	}
	reponse.NewResponse(http.StatusOK, identity).Apply(c)
}

func IdentityUpdate(c *gin.Context) {
	id := c.Param("id")
	if !regexp.MustCompile("\\d+").Match([]byte(id)) {
		reponse.NewError(http.StatusUnprocessableEntity, "given id is not numeric").Apply(c)
		return
	}
	var identity model.Identity
	if err := db.Get().Find(&identity, id).Error; err != nil {
		reponse.NewError(http.StatusNotFound, err.Error()).Apply(c)
		return
	}

	var json UpdatePasswordForm
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	identity.Password = crypto.NewPassword(json.Password)
	db.Get().Save(identity)
	reponse.NewResponse(http.StatusOK, identity).Apply(c)
}

func IdentityDelete(c *gin.Context) {
	id := c.Param("id")
	if !regexp.MustCompile("\\d+").Match([]byte(id)) {
		reponse.NewError(http.StatusUnprocessableEntity, "given id is not numeric").Apply(c)
		return
	}
	var identity model.Identity
	if err := db.Get().Find(&identity, id).Error; err != nil {
		reponse.NewError(http.StatusNotFound, err.Error()).Apply(c)
		return
	}

	db.Get().Delete(identity)
	reponse.NewResponse(http.StatusOK, identity).Apply(c)
}
