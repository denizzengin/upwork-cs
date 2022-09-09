// Copyright 2022 TCDZENGIN
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/denizzengin/upwork-cs/internal/bll"
	"github.com/denizzengin/upwork-cs/internal/storage"
	"github.com/denizzengin/upwork-cs/model"
	"github.com/gorilla/mux"
)

func CreateUser(s storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userDTO := new(model.UserDTO)
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&userDTO); err != nil {
			CreateResponse(w, http.StatusBadRequest, &Response{Error: "Bad request"})
			return
		}

		err := userDTO.Validate()
		if err != nil {
			CreateResponse(w, http.StatusBadRequest, &Response{Error: "Bad request"})
			return
		}
		defer r.Body.Close()

		user, err := bll.FindUserByEmail(s, userDTO)
		if err != nil {
			CreateResponse(w, http.StatusInternalServerError, &Response{Error: "server error"})
			return
		}

		if user != nil && user.ID > 0 {
			CreateResponse(w, http.StatusForbidden, &Response{Error: "User with that email already exists"})
			return
		}

		newUser, err := bll.CreateUser(s, userDTO)
		if err != nil {
			CreateResponse(w, http.StatusInternalServerError, &Response{Error: "server error"})
			return
		}
		CreateResponse(w, http.StatusOK, model.ToUserDTO(newUser))
	}
}

func Login(s storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userDTO := new(model.UserDTO)
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&userDTO); err != nil {
			CreateResponse(w, http.StatusBadRequest, &Response{Error: "Bad request"})
			return
		}

		err := userDTO.Validate()
		if err != nil {
			CreateResponse(w, http.StatusBadRequest, &Response{Error: "Check request"})
			return
		}
		defer r.Body.Close()

		token, err := bll.Login(s, userDTO)
		if err != nil {
			CreateResponse(w, http.StatusUnauthorized, &Response{Error: err.Error()})
			return
		}
		CreateResponse(w, http.StatusOK, token)
	}
}

func FindAll(s storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer r.Body.Close()

		users, err := s.Users().All()
		if err != nil {
			CreateResponse(w, http.StatusNotFound, &Response{Error: "User with that id does not exist"})
			return
		}

		CreateResponse(w, http.StatusOK, model.ToUserDTOs(users))
	}
}

func UpdateUser(s storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			CreateResponse(w, http.StatusBadRequest, &Response{Error: "Bad request"})
			return
		}

		var userDTO model.UserDTO
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&userDTO); err != nil {
			CreateResponse(w, http.StatusBadRequest, &Response{Error: "Bad request"})
			return
		}

		defer r.Body.Close()
		userDTO.ID = uint(id)
		record, err := bll.FindUserByID(s, &userDTO, true)
		if err != nil {
			CreateResponse(w, http.StatusNotFound, &Response{Error: "User with that id does not exist"})
			return
		}

		userDTO = *model.FillEmptyFields(&userDTO, record)
		updatedUser, err := bll.UpdateUser(s, &userDTO)
		if err != nil {
			CreateResponse(w, http.StatusInternalServerError, &Response{Error: "server error"})
			return
		}

		CreateResponse(w, http.StatusOK, model.ToUserDTO(updatedUser))
	}
}

func FindUserByName(s storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer r.Body.Close()
		ctx := r.Context()
		u, ok := ctx.Value("user").(*model.Claims)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		user, err := bll.FindUserByName(s, &model.UserDTO{Name: u.Username})
		if err != nil {
			CreateResponse(w, http.StatusNotFound, &Response{Error: "User with that id does not exist"})
			return
		}

		CreateResponse(w, http.StatusOK, user)
	}
}

func FindUserByID(s storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			CreateResponse(w, http.StatusBadRequest, &Response{Error: "Bad request"})
			return
		}

		defer r.Body.Close()

		user, err := bll.FindUserByID(s, &model.UserDTO{ID: uint(id)})
		if err != nil {
			CreateResponse(w, http.StatusNotFound, &Response{Error: "User with that id does not exist"})
			return
		}

		CreateResponse(w, http.StatusOK, user)
	}
}
