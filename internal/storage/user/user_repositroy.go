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

package user

import (
	"github.com/denizzengin/upwork-cs/model"
	"github.com/denizzengin/upwork-cs/pkg/logger"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (p *Repository) RoleByName(roleName string) (*model.UserRole, error) {

	role := new(model.UserRole)
	err := p.db.Where(`name = ?  and is_deleted = ?`, roleName, false).First(&role).Error
	if err != nil {
		//logger.Log.Error("Error during create user userRole...")
		return nil, err
	}

	return role, nil
}

func (p *Repository) CheckRoleRight(roleId uint, pattern string) bool {

	roleRight := new(model.RoleRight)
	err := p.db.Where(`role_id = ? and pattern = ? and is_deleted = ?`, roleId, pattern, false).First(&roleRight).Error
	if err != nil {
		//logger.Log.Error("Error during create user userRole...")
		return false
	}

	if roleRight == nil || roleRight.ID <= 0 {
		return false
	}
	return true
}

func (p *Repository) Create(user *model.User) (*model.User, error) {

	userRole := &model.UserRole{}
	err := p.db.Transaction(func(tx *gorm.DB) error {
		err := p.db.Create(&user).Error
		if err != nil {
			logger.Log.Error("Error during create user process...")
			return err
		}

		userRole.UserId = user.ID
		userRole.RoleId = 2 // todo read from db
		err = p.db.Create(&userRole).Error
		if err != nil {
			logger.Log.Error("Error during create user process...")
			return err
		}

		return nil
	})
	return user, err
}

func (p *Repository) Update(user *model.User) (*model.User, error) {
	err := p.db.Save(&user).Error
	if err != nil {
		logger.Log.Error("Error during update user process...")
		return nil, err
	}

	return user, nil
}

func (p *Repository) Delete(user *model.User) error {
	err := p.db.Save(user).Error
	if err != nil {
		logger.Log.Error("Error during delete user process...")
	}
	return err
}

func (p *Repository) FindByID(id uint) (*model.User, error) {
	user := new(model.User)
	err := p.db.Where(`id = ? and is_deleted = ?`, id, false).First(&user).Error
	if err != nil {
		logger.Log.Error("Error during finding user process...")
		return nil, err
	}
	return user, err
}

func (p *Repository) FindByName(userName string) (*model.User, error) {
	user := new(model.User)
	err := p.db.Where(`name = ? and is_deleted = ?`, userName, false).First(&user).Error
	if err != nil {
		logger.Log.Error("Error during finding user process...")
		return nil, err
	}
	return user, err
}

func (p *Repository) All() ([]model.User, error) {
	users := []model.User{}
	err := p.db.Where(`is_deleted = ?`, false).Find(&users).Error
	return users, err
}

func (p *Repository) FindByEmail(email string) (*model.User, error) {
	user := new(model.User)
	err := p.db.Where(`email = ? and is_deleted = ?`, email, false).First(&user).Error
	return user, err
}

func (p *Repository) CreateRole(role *model.Role) (*model.Role, error) {
	err := p.db.Create(&role).Error
	if err != nil {
		logger.Log.Error("Error during create user Role...")
		return nil, err
	}

	return role, nil
}

func (p *Repository) CreateUserRole(userRole *model.UserRole) (*model.UserRole, error) {
	err := p.db.Create(&userRole).Error
	if err != nil {
		logger.Log.Error("Error during create user userRole...")
		return nil, err
	}

	return userRole, nil
}

func (p *Repository) DeleteUserRole(userRole *model.UserRole) error {
	err := p.db.Save(userRole).Error
	if err != nil {
		logger.Log.Error("Error during delete user process...")
	}
	return err
}

func (p *Repository) FindUserRoles(userId uint) (*[]model.UserRole, error) {
	userRoles := new([]model.UserRole)
	err := p.db.Where(`user_id = ? and is_deleted = ?`, userId, false).Find(&userRoles).Error
	if err != nil {
		logger.Log.Error("Error during create user userRole...")
		return nil, err
	}

	return userRoles, nil
}

func (p *Repository) Migrate() error {
	err := p.db.AutoMigrate(&model.User{}, &model.Role{}, &model.UserRole{}, &model.RoleRight{})
	return err
}
