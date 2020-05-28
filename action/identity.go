package action

import (
	"github.com/gin-gonic/gin"
	"github.com/jschaefer-io/IDaaS/db"
	"github.com/jschaefer-io/IDaaS/model"
	"net/http"
)


type UpdatePasswordForm struct {
	Password string `json:"password"`
}

type AddIdentityForm struct {
	UpdatePasswordForm
	Email     string `json:"email" binding:"required,email"`
}


func AddIdentity(c *gin.Context){
	var json AddIdentityForm
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := model.NewIdentity(json.Email, json.Password)
	db.Get().Create(&id)
	c.JSON(200, id)
}


func GetIdentity(c *gin.Context){
	var ids []model.Identity
	db.Get().Find(&ids)
	c.JSON(200, ids)
}
