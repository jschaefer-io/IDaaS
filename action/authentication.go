package action

import (
	"fmt"
	"github.com/gorilla/context"
	"github.com/jschaefer-io/IDaaS/crypto"
	"github.com/jschaefer-io/IDaaS/db"
	"github.com/jschaefer-io/IDaaS/model"
	"github.com/jschaefer-io/IDaaS/reponse"
	"math"
	"net/http"
	"time"
)

// Returns data for the current logged in user
func AuthMe(w http.ResponseWriter, r *http.Request) {
	identity := context.Get(r, "identity").(model.Identity)
	reponse.NewResponse(http.StatusOK, identity).Apply(w)
}

// Handles the token renewing
// only allows renewing of the token 60 seconds
// before expiration
func AuthRenew(w http.ResponseWriter, r *http.Request) {
	token, err := crypto.ExtractJWT(r)
	if err != nil {
		reponse.NewError(http.StatusUnprocessableEntity, "Token missing from Authorisation header").Apply(w)
		return
	}
	identity, claims, t, err := model.Identity{}.CheckJwt(token)

	if err != nil || !t.Valid {
		reponse.NewError(http.StatusBadRequest, "given token is invalid").Apply(w)
		return
	}

	// only allow renewing the token if 60seconds or
	// less are remaining in the exp date
	exp := time.Unix(int64(claims["exp"].(float64)), 0)
	deltaTime := exp.Sub(time.Now()).Seconds()
	if deltaTime > 60 {
		renewTime := int64(math.Round(deltaTime)) - 60
		reponse.NewError(http.StatusBadRequest, fmt.Sprintf("token can only be renewed in %d seconds", renewTime)).Apply(w)
		return
	}

	// Renew Token
	fmt.Println(identity)
	newExp, jwt, err := crypto.NewJWT(identity.Token, map[string]interface{}{
		"user-id": identity.ID,
	})

	if err != nil {
		reponse.NewError(http.StatusInternalServerError, "An error occurred creating the auth token").Apply(w)
		return
	}

	reponse.NewResponse(http.StatusOK, map[string]interface{}{
		"error":   false,
		"token":   jwt,
		"exp":     newExp,
		"message": "Token renewed successful",
	}).Apply(w)
}

// Handles an identity login
func AuthLogin(w http.ResponseWriter, r *http.Request) {

	// Map json data to identity form
	var json model.IdentityLogin
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

	// Return an error if the identity is currently unconfirmed
	if !identity.Confirmed {
		reponse.NewError(http.StatusConflict, "Identity unconfirmed").Apply(w)
		return
	}

	// Issue and return a jwt
	// with the user id in the free
	// claim data
	exp, jwt, err := crypto.NewJWT(identity.Token, map[string]interface{}{
		"user-id": identity.ID,
	})

	if err != nil {
		reponse.NewError(http.StatusInternalServerError, "An error occurred creating the auth token").Apply(w)
		return
	}

	reponse.NewResponse(http.StatusOK, map[string]interface{}{
		"error":   false,
		"token":   jwt,
		"exp":     exp,
		"message": "Authentication successful",
	}).Apply(w)
}
