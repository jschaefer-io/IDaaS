package action

import (
	"github.com/gin-gonic/gin"
	"github.com/jschaefer-io/IDaaS/crypto"
	"github.com/jschaefer-io/IDaaS/db"
	"github.com/jschaefer-io/IDaaS/model"
	"github.com/jschaefer-io/IDaaS/reponse"
	"net/http"
)

// Returns data for the current logged in user
func AuthMe(c *gin.Context) {
	identity, _ := c.Get("identity")
	reponse.NewResponse(http.StatusOK, identity).Apply(c)
}

// Handles an identity login
func AuthLogin(c *gin.Context) {

	// Map json data to identity form
	var json model.Identity
	if err := c.ShouldBindJSON(&json); err != nil {
		model.ValidationError(err).Apply(c)
		return
	}

	// Try to find the identity and compare passwords
	identity := model.Identity{}
	err := db.Get().Where("email = ?", json.Email).Find(&identity).Error
	pwd := crypto.Password(identity.Password)
	if err != nil || !pwd.Compare(json.Password) {
		reponse.NewError(http.StatusUnauthorized, "Permission denied").Apply(c)
		return
	}

	// Issue and return a jwt
	// with the user id in the free
	// claim data
	jwt, err := crypto.NewJWT(identity.Token, map[string]interface{}{
		"user-id": identity.ID,
	})

	if err != nil {
		reponse.NewError(http.StatusInternalServerError, "An error occurred creating the auth token").Apply(c)
		return
	}

	reponse.NewResponse(http.StatusOK, map[string]interface{}{
		"error":   false,
		"token":   jwt,
		"message": "Authentication successful",
	}).Apply(c)
}
