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

package logger

import (
	"github.com/denizzengin/upwork-cs/pkg/common"
	"go.uber.org/zap"
)

var Log *zap.Logger = New()

func New() *zap.Logger {

	var cfg zap.Config
	if common.IsProduction() {
		cfg = zap.NewProductionConfig()
		cfg.Encoding = "json"
		cfg.OutputPaths = []string{"../../upwork-cs.log"}
	} else {
		cfg = zap.NewDevelopmentConfig()
		cfg.Encoding = "console"
	}
	var loggerInstance *zap.Logger
	loggerInstance, _ = cfg.Build()
	return loggerInstance
}
