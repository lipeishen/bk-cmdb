/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package logics

import (
	"gopkg.in/redis.v5"
	"net/http"

	"configcenter/src/auth/extensions"
	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/errors"
	"configcenter/src/common/language"
	"configcenter/src/common/util"
	"configcenter/src/thirdpartyclient/esbserver"
)

type Logics struct {
	*backbone.Engine
	esbServ     esbserver.EsbClientInterface
	ErrHandle   errors.DefaultCCErrorIf
	cache       *redis.Client
	header      http.Header
	rid         string
	ownerID     string
	user        string
	ccErr       errors.DefaultCCErrorIf
	ccLang      language.DefaultCCLanguageIf
	AuthManager *extensions.AuthManager
}

// NewLogics get logics handle
func NewLogics(b *backbone.Engine, header http.Header, authManager *extensions.AuthManager) *Logics {
	lang := util.GetLanguage(header)
	return &Logics{
		Engine:      b,
		header:      header,
		rid:         util.GetHTTPCCRequestID(header),
		ccErr:       b.CCErr.CreateDefaultCCErrorIf(lang),
		ccLang:      b.Language.CreateDefaultCCLanguageIf(lang),
		user:        util.GetUser(header),
		ownerID:     util.GetOwnerID(header),
		AuthManager: authManager,
	}
}

// NewFromHeader new Logic from header
func (lgc *Logics) NewFromHeader(header http.Header) *Logics {
	lang := util.GetLanguage(header)
	rid := util.GetHTTPCCRequestID(header)
	if rid == "" {
		if lgc.rid == "" {
			rid = util.GenerateRID()
		} else {
			rid = lgc.rid
		}
		header.Set(common.BKHTTPCCRequestID, rid)
	}
	newLgc := &Logics{
		header:  header,
		Engine:  lgc.Engine,
		rid:     rid,
		cache:   lgc.cache,
		esbServ: lgc.esbServ,
		user:    util.GetUser(header),
		ownerID: util.GetOwnerID(header),
	}
	// if language not exist, use old language
	if lang == "" {
		newLgc.ccErr = lgc.ccErr
		newLgc.ccLang = lgc.ccLang
	} else {
		newLgc.ccErr = lgc.CCErr.CreateDefaultCCErrorIf(lang)
		newLgc.ccLang = lgc.Language.CreateDefaultCCLanguageIf(lang)
	}
	return newLgc
}

type Logic struct {
	*backbone.Engine
}
