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
	"github.com/denizzengin/upwork-cs/internal/storage"
	"github.com/denizzengin/upwork-cs/pkg/logger"
)

func Migrate(s storage.Store) {
	if err := s.Users().Migrate(); err != nil {
		logger.Log.Error("failed to migrate users")
	}
}
