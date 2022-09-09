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
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/denizzengin/upwork-cs/internal/storage/user"
	"github.com/denizzengin/upwork-cs/pkg/config"
)

type Database struct {
	db    *gorm.DB
	users UserRepository
}

func GetDBConnection(cfg config.DatabaseConfiguration) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", cfg.Host, cfg.Username, cfg.Password, cfg.Name, cfg.Port)
	db, err = gorm.Open(postgres.Open(dsn))
	if err != nil {
		return nil, fmt.Errorf("could not open postgresql connection: %v", err)
	}

	return db, err
}

// New opens a database according to configuration.
func New(db *gorm.DB) *Database {
	return &Database{
		db:    db,
		users: user.NewRepository(db),
	}
}

func (db *Database) Users() UserRepository {
	return db.users
}

func (db *Database) Ping() error {
	sqlDB, err := db.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
