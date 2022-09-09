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
package storage

import (
	"github.com/denizzengin/upwork-cs/model"
)

type UserRepository interface {
	Create(user *model.User) (*model.User, error)
	Update(user *model.User) (*model.User, error)
	Delete(user *model.User) error
	FindByID(id uint) (*model.User, error)
	FindByName(userName string) (*model.User, error)
	All() ([]model.User, error)
	FindByEmail(email string) (*model.User, error)
	CreateRole(role *model.Role) (*model.Role, error)
	CreateUserRole(userRole *model.UserRole) (*model.UserRole, error)
	DeleteUserRole(userRole *model.UserRole) error
	FindUserRoles(userId uint) (*[]model.UserRole, error)
	CheckRoleRight(roleId uint, pattern string) bool
	RoleByName(roleName string) (*model.UserRole, error)
	Migrate() error
}
