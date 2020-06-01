package action

import (
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
func (b Identity) Index(w http.ResponseWriter, r *http.Request) {
	var ids []model.Identity
	db.Get().Find(&ids)
	reponse.NewResponse(200, ids).Apply(w)
}

// Main creat route
func (b Identity) Create(w http.ResponseWriter, r *http.Request) {
	var json model.Identity
	if err := model.BindJson(&json, r.Body); err != nil {
		model.ValidationError(err).Apply(w)
		return
	}

	id := model.NewIdentity(json.Email, json.Password)
	_ = db.Get().Create(&id)

	//fmt.Println(res.Error.Error())
	reponse.NewResponse(200, id).Apply(w)
}

// Main show single route
func (b Identity) Show(w http.ResponseWriter, r *http.Request, id int) {
	identity, err := model.Identity{}.Find(id)
	if err != nil {
		reponse.NewError(http.StatusUnprocessableEntity, err.Error()).Apply(w)
		return
	}
	reponse.NewResponse(http.StatusOK, identity).Apply(w)
}

// Main update route
func (b Identity) Update(w http.ResponseWriter, r *http.Request, id int) {
	identity, err := model.Identity{}.Find(id)
	if err != nil {
		reponse.NewError(http.StatusUnprocessableEntity, err.Error()).Apply(w)
		return
	}

	var json model.PasswordForm
	if err := model.BindJson(&json, r.Body); err != nil {
		model.ValidationError(err).Apply(w)
		return
	}

	pwd := crypto.NewPassword(json.Password)
	identity.Password = pwd.String()
	db.Get().Save(&identity)
	reponse.NewResponse(http.StatusOK, identity).Apply(w)
}

// Main delete route
func (b Identity) Delete(w http.ResponseWriter, r *http.Request, id int) {
	identity, err := model.Identity{}.Find(id)
	if err != nil {
		reponse.NewError(http.StatusUnprocessableEntity, err.Error()).Apply(w)
		return
	}

	db.Get().Delete(&identity)
	reponse.NewResponse(http.StatusOK, identity).Apply(w)
}
