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

package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/golang-jwt/jwt"
	"time"
)

type User struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `gorm:"<-:create" json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsDeleted bool      `json:"is_deleted"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
}

type Role struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `gorm:"<-:create" json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsDeleted bool      `json:"is_deleted"`
	RoleName  string    `json:"role_name"`
}

type UserRole struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `gorm:"<-:create" json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsDeleted bool      `json:"is_deleted"`
	UserId    uint      `json:"user_id"`
	RoleId    uint      `json:"role_id"`
}

type RoleRight struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `gorm:"<-:create" json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsDeleted bool      `json:"is_deleted"`
	Pattern   string    `json:"pattern"`
	RoleId    uint      `json:"role_id"`
}

type UserDTO struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	Roles    []Role `json:"roles,omitempty"`
}

func (u UserDTO) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Name, validation.Required),
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.Required, validation.Length(6, 50)),
	)
}

func ToUser(userDTO *UserDTO) *User {
	return &User{
		ID:       userDTO.ID,
		Name:     userDTO.Name,
		Email:    userDTO.Email,
		Password: userDTO.Password,
	}
}

func ToUserDTO(user *User, includePassword ...bool) *UserDTO {
	userDto := &UserDTO{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
	if len(includePassword) != 0 && includePassword[0] {
		userDto.Password = user.Password
	}

	return userDto
}

func ToUserDTOs(users []User) []*UserDTO {
	userDtos := make([]*UserDTO, 0)

	for _, user := range users {
		userDtos = append(userDtos, &UserDTO{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		})
	}

	return userDtos
}

func FillEmptyFields(userDto *UserDTO, dbRowUser *UserDTO) *UserDTO {
	if userDto.Email == "" {
		userDto.Email = dbRowUser.Email
	}

	if userDto.Name == "" {
		userDto.Name = dbRowUser.Name
	}

	if userDto.Password == "" {
		userDto.Password = dbRowUser.Password
	}

	return userDto
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type TokenResponse struct {
	AccessToken string `json:"accesstoken"`
	ExpireAt    int64  `json:"expireAt"`
}
