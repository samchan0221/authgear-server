// Copyright 2015-present Oursky Ltd.
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

package audit

import (
	"github.com/sirupsen/logrus"

	"github.com/skygeario/skygear-server/pkg/server/router"
)

// We do not acquire a logger from our logging package
// because the log level of this logger
// should not be affected the global log level
var trailLogger = logrus.New()

const (
	enabledLevel  = logrus.InfoLevel
	disabledLevel = logrus.PanicLevel
)

func init() {
	trailLogger.Formatter = &logrus.JSONFormatter{}
	trailLogger.Level = disabledLevel
}

type Event int

const (
	_ Event = iota

	// EventLoginSuccess represents Login Success
	EventLoginSuccess

	// EventLoginFailure represents Login Failure
	EventLoginFailure

	// EventLogout represents Logout
	EventLogout

	// EventSignup represents Signup
	EventSignup

	// EventChangePassword represents Change Password
	EventChangePassword

	// EventChangeRoles represents Change Roles
	EventChangeRoles
)

func (e Event) String() string {
	switch e {
	case EventLoginSuccess:
		return "login_success"
	case EventLoginFailure:
		return "login_failure"
	case EventLogout:
		return "logout"
	case EventSignup:
		return "signup"
	case EventChangePassword:
		return "change_password"
	case EventChangeRoles:
		return "change_roles"
	default:
		return ""
	}
}

type Entry struct {
	Event         Event
	Admin         bool
	AuthID        string
	Data          map[string]interface{}
	RemoteAddr    string
	XForwardedFor string
	XRealIP       string
	Forwarded     string
}

func (e Entry) WithRouterPayload(payload *router.Payload) Entry {
	// If we directly assign to e, we will have an ineffective
	// assignment lint error
	ee := e
	if payload != nil {
		if remoteAddr, ok := payload.Meta["remote_addr"].(string); ok {
			ee.RemoteAddr = remoteAddr
		}
		if xff, ok := payload.Meta["x_forwarded_for"].(string); ok {
			ee.XForwardedFor = xff
		}
		if xri, ok := payload.Meta["x_real_ip"].(string); ok {
			ee.XRealIP = xri
		}
		if forwarded, ok := payload.Meta["forwarded"].(string); ok {
			ee.Forwarded = forwarded
		}
	}
	return ee
}

func (e *Entry) toLogrusFields() logrus.Fields {
	return logrus.Fields{
		"event":                e.Event.String(),
		"auth_id":              e.AuthID,
		"data":                 e.Data,
		"tcp_remote_addr":      e.RemoteAddr,
		"http_x_forwarded_for": e.XForwardedFor,
		"http_x_real_ip":       e.XRealIP,
		"http_forwarded":       e.Forwarded,
	}
}

func Trail(entry Entry) {
	trailLogger.WithFields(entry.toLogrusFields()).Info("audit_trail")
}

func SetTrailEnabled(enabled bool) {
	if enabled {
		trailLogger.Level = enabledLevel
	} else {
		trailLogger.Level = disabledLevel
	}
}
