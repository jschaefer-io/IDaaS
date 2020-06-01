package action

import (
	"github.com/gorilla/context"
	"github.com/jschaefer-io/IDaaS/crypto"
	"github.com/jschaefer-io/IDaaS/db"
	"github.com/jschaefer-io/IDaaS/model"
	"github.com/jschaefer-io/IDaaS/reponse"
	"net/http"
)

// Returns data for the current logged in user
func AuthMe(w http.ResponseWriter, r *http.Request) {
	identity := context.Get(r, "identity").(model.Identity)
	reponse.NewResponse(http.StatusOK, identity).Apply(w)
}

// Handles an identity login
func AuthLogin(w http.ResponseWriter, r *http.Request) {

	// Map json data to identity form
	var json model.Identity
	if err := model.BindJson(&json, r.Body); err != nil {
		model.ValidationError(err).Apply(w)
		return
	}

	// Try to find the identity and compare passwords
	identity := model.Identity{}
	err := db.Get().Where("email = ?", json.Email).Find(&identity).Error
	pwd := crypto.Password(identity.Password)
	if err != nil || !pwd.Compare(json.Password) {
		reponse.NewError(http.StatusUnauthorized, "Permission denied").Apply(w)
		return
	}

	// Issue and return a jwt
	// with the user id in the free
	// claim data
	jwt, err := crypto.NewJWT(identity.Token, map[string]interface{}{
		"user-id": identity.ID,
	})

	if err != nil {
		reponse.NewError(http.StatusInternalServerError, "An error occurred creating the auth token").Apply(w)
		return
	}

	reponse.NewResponse(http.StatusOK, map[string]interface{}{
		"error":   false,
		"token":   jwt,
		"message": "Authentication successful",
	}).Apply(w)
}
