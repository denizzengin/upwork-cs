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

package bll

import (
	b64 "encoding/base64"
	"errors"
	"time"

	"github.com/denizzengin/upwork-cs/internal/storage"
	"github.com/denizzengin/upwork-cs/model"
	"github.com/denizzengin/upwork-cs/pkg/logger"
	"github.com/golang-jwt/jwt"
)

const (
	ErrAlreadyExist = "already exists record"
	ErrNotExist     = "record doesn't exist"
)

func createToken(userDTO *model.UserDTO) (*model.TokenResponse, error) {

	expirationTime := time.Now().Add(5 * time.Hour)
	claims := &model.Claims{
		Username: userDTO.Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte("my_secret_key"))
	if err != nil {
		return nil, err
	}
	encodedToken := b64.StdEncoding.EncodeToString([]byte(tokenStr))
	tokenResponse := &model.TokenResponse{AccessToken: encodedToken, ExpireAt: claims.StandardClaims.ExpiresAt}
	return tokenResponse, nil

}

func Login(s storage.Store, userDTO *model.UserDTO) (*model.TokenResponse, error) {

	userRecord, _ := s.Users().FindByEmail(userDTO.Email)

	if userRecord != nil && userRecord.ID <= 0 {
		logger.Log.Error("record not found exist")
		return nil, errors.New("user not found")
	}

	if userRecord.Password != userDTO.Password {
		return nil, errors.New("wrong email or password")
	}

	token, err := createToken(userDTO)

	return token, err
}

func CreateUser(s storage.Store, userDTO *model.UserDTO) (*model.User, error) {
	var err error

	record, _ := s.Users().FindByEmail(userDTO.Email)

	if record != nil && record.ID > 0 {
		logger.Log.Error("record already exist")
		return nil, errors.New(ErrAlreadyExist)
	}

	createdUser, err := s.Users().Create(model.ToUser(userDTO))
	if err != nil {
		logger.Log.Error("Error while creating user")
		return nil, err
	}

	return createdUser, nil
}

func UpdateUser(s storage.Store, userDTO *model.UserDTO) (*model.User, error) {
	var err error

	record, _ := s.Users().FindByID(userDTO.ID)

	if record == nil {
		logger.Log.Error(ErrNotExist)
		return nil, errors.New(ErrNotExist)
	}

	updatedUser, err := s.Users().Update(model.ToUser(userDTO))
	if err != nil {
		logger.Log.Error("Error while updating user")
		return nil, err
	}

	oldRoles, _ := s.Users().FindUserRoles(updatedUser.ID)

	for _, v := range *oldRoles {
		v.IsDeleted = true
		s.Users().DeleteUserRole(&v)
	}

	if userDTO.Roles != nil && len(userDTO.Roles) > 0 {
		for _, role := range userDTO.Roles {
			recordRole, _ := s.Users().RoleByName(role.RoleName)
			s.Users().CreateUserRole(&model.UserRole{UserId: updatedUser.ID, RoleId: recordRole.ID})
		}
	}

	return updatedUser, nil
}

func DeleteUser(s storage.Store, userDTO *model.UserDTO) error {

	//todo validation comes here

	record, _ := s.Users().FindByID(userDTO.ID)

	if record == nil {
		logger.Log.Error(ErrNotExist)
		return errors.New(ErrNotExist)
	}

	record.IsDeleted = true
	err := s.Users().Delete(record)
	if err != nil {
		logger.Log.Error("Error while updating user")
		return err
	}

	return nil
}

func FindUserByID(s storage.Store, userDTO *model.UserDTO, includePassword ...bool) (*model.UserDTO, error) {
	//todo validation comes here

	record, _ := s.Users().FindByID(userDTO.ID)

	if record == nil {
		logger.Log.Error(ErrNotExist)
		return nil, errors.New(ErrNotExist)
	}

	if len(includePassword) > 0 && includePassword[0] {
		return model.ToUserDTO(record, true), nil
	}
	return model.ToUserDTO(record), nil
}

func FindUserByName(s storage.Store, userDTO *model.UserDTO) (*model.UserDTO, error) {

	record, _ := s.Users().FindByName(userDTO.Name)

	if record == nil {
		logger.Log.Error(ErrNotExist)
		return nil, errors.New(ErrNotExist)
	}

	return model.ToUserDTO(record), nil
}

func FindUserByEmail(s storage.Store, userDTO *model.UserDTO) (*model.UserDTO, error) {
	//todo validation comes here

	record, _ := s.Users().FindByEmail(userDTO.Email)

	if record == nil {
		logger.Log.Error(ErrNotExist)
		return nil, errors.New(ErrNotExist)
	}

	return model.ToUserDTO(record), nil
}

func All(s storage.Store) ([]*model.UserDTO, error) {
	//todo validation comes here

	records, err := s.Users().All()

	if err != nil {
		logger.Log.Error(ErrNotExist)
		return nil, errors.New(ErrNotExist)
	}

	return model.ToUserDTOs(records), nil
}
