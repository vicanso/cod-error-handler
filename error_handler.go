// Copyright 2018 tree xie
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

package errorhandler

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/vicanso/cod"
	"github.com/vicanso/hes"
)

type (
	// Config error handler config
	Config struct {
		Skipper cod.Skipper
	}
)

const (
	errErrorHandlerCategory = "cod-error-handler"
)

// NewDefault create a default error handler
func NewDefault() cod.Handler {
	return New(Config{})
}

// New create a error handler
func New(config Config) cod.Handler {
	skipper := config.Skipper
	if skipper == nil {
		skipper = cod.DefaultSkipper
	}
	return func(c *cod.Context) error {
		if skipper(c) {
			return c.Next()
		}
		err := c.Next()
		// 如果没有出错，直接返回
		if err == nil {
			return nil
		}
		he, ok := err.(*hes.Error)
		if !ok {
			he = &hes.Error{
				StatusCode: http.StatusInternalServerError,
				Message:    err.Error(),
				Category:   errErrorHandlerCategory,
				Err:        err,
			}
		}
		// 出错响应，直接设置为no cache
		c.NoCache()
		c.StatusCode = he.StatusCode
		// 默认以json的形式返回
		buf, e := json.Marshal(he)
		if e != nil {
			return e
		}
		c.BodyBuffer = bytes.NewBuffer(buf)
		c.SetHeader(cod.HeaderContentType, cod.MIMEApplicationJSON)
		return nil
	}
}
