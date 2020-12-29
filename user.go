package main

import (
	"crypto/rsa"
	"fmt"
	"github.com/fatih/structs"
	"github.com/go-chi/chi"
	"github.com/jschaefer-io/IDaaS/crypto"
	"github.com/jschaefer-io/IDaaS/model"
	"github.com/jschaefer-io/IDaaS/util"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"sync"
)

func (s *Server) UsersGetAll() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		users := make([]model.User, 0)
		s.db.Find(&users)
		util.MarshalToResponseWriter(users, writer, http.StatusOK)
	}
}

func (s *Server) UsersGetSingle() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		id := chi.URLParam(request, "userId")
		usr := model.User{}
		err := s.db.First(&usr, id)
		if err.Error != nil {
			util.MarshalToResponseWriter(util.NewErrorObject(fmt.Sprintf("user %v not found", id)), writer, http.StatusNotFound)
			return
		}
		util.MarshalToResponseWriter(usr, writer, http.StatusOK)
	}
}

func (s *Server) UsersCreate() http.HandlerFunc {
	type userCreateForm struct {
		Name     string `json:"name" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}
	return func(writer http.ResponseWriter, request *http.Request) {
		data := new(userCreateForm)
		if err := model.BindJson(data, request.Body); err != nil {
			util.MarshalToResponseWriter(util.ValidationError(err), writer, http.StatusUnprocessableEntity)
			return
		}

		// Check if the given email address is unique
		uniqueCheck := make([]model.User, 0)
		s.db.Where("email = ?", data.Email).Find(&uniqueCheck)
		if len(uniqueCheck) > 0 {
			util.MarshalToResponseWriter(util.NewErrorObject(fmt.Sprintf("a user with this the email %s already exists", data.Email)), writer, http.StatusBadRequest)
			return
		}

		user := model.User{
			Name:  data.Name,
			Email: data.Email,
		}
		user.SetPassword(data.Password)

		if err := s.db.Create(&user); err.Error != nil {
			util.MarshalToResponseWriter(util.NewErrorObject(err.Error.Error()), writer, http.StatusInternalServerError)
			return
		}
		util.MarshalToResponseWriter(user, writer, http.StatusCreated)
	}
}

func (s *Server) UserUpdateSingle() http.HandlerFunc {
	type userUpdateForm struct {
		Name     string `json:"name" validate:"omitempty"`
		Email    string `json:"email" validate:"omitempty,email"`
		Password string `json:"password" validate:"omitempty"`
	}
	return func(writer http.ResponseWriter, request *http.Request) {
		// Prepare input data
		data := new(userUpdateForm)
		if err := model.BindJson(data, request.Body); err != nil {
			util.MarshalToResponseWriter(util.ValidationError(err), writer, http.StatusUnprocessableEntity)
			return
		}
		updateData := make(map[string]interface{})
		for key, value := range structs.Map(data) {
			if value != "" {
				if key == "Password" {
					usr := model.User{}
					usr.SetPassword(value.(string))
					value = usr.Password
				}
				updateData[key] = value
			}
		}

		// Find User
		id := chi.URLParam(request, "userId")
		usr := model.User{}
		err := s.db.First(&usr, id)
		if err.Error != nil {
			util.MarshalToResponseWriter(util.NewErrorObject(fmt.Sprintf("user %v not found", id)), writer, http.StatusNotFound)
			return
		}

		// Update User
		if err := s.db.Model(&usr).Updates(updateData); err.Error != nil {
			util.MarshalToResponseWriter(util.NewErrorObject(err.Error.Error()), writer, http.StatusInternalServerError)
			return
		}
		util.MarshalToResponseWriter(usr, writer, http.StatusOK)
	}
}

func (s *Server) UserDeleteSingle() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		id := chi.URLParam(request, "userId")
		del := s.db.Delete(model.User{}, id)
		if del.Error != nil {
			util.MarshalToResponseWriter(util.NewErrorObject(del.Error.Error()), writer, http.StatusInternalServerError)
			return
		}
		if del.RowsAffected == 0 {
			util.MarshalToResponseWriter(util.NewErrorObject(fmt.Sprintf("user %v not found", id)), writer, http.StatusNotFound)
			return
		}
		util.MarshalToResponseWriter(nil, writer, http.StatusNoContent)
	}
}

func (s *Server) UserAuthenticate() http.HandlerFunc {
	type userLoginForm struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	var init sync.Once
	var rsaSecret *rsa.PrivateKey
	var rsaError error
	return func(writer http.ResponseWriter, request *http.Request) {
		init.Do(func() {
			rsaSecret, rsaError = crypto.ReadPrivateRsaKey()
		})
		if rsaError != nil {
			util.MarshalToResponseWriter(util.NewErrorObject("an unknown error occurred"), writer, http.StatusInternalServerError)
			panic(rsaError)
		}

		// check user credentials
		data := new(userLoginForm)
		usr := model.User{}
		if err := model.BindJson(data, request.Body); err != nil {
			util.MarshalToResponseWriter(util.NewErrorObject(err.Error()), writer, http.StatusInternalServerError)
			return
		}
		if err := s.db.Where("email = ?", data.Email).First(&usr); err.Error != nil {
			util.MarshalToResponseWriter(util.NewErrorObject(err.Error.Error()), writer, http.StatusBadRequest)
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(data.Password)); err != nil {
			util.MarshalToResponseWriter(util.NewErrorObject("permission denied"), writer, http.StatusUnauthorized)
			return
		}

		tokenString, err := crypto.NewJwt(rsaSecret, usr.ID)
		if err != nil {
			util.MarshalToResponseWriter(util.NewErrorObject("an unknown error occurred"), writer, http.StatusInternalServerError)
			return
		}

		util.MarshalToResponseWriter(map[string]interface{}{
			"token": tokenString,
		}, writer, http.StatusOK)
	}
}
