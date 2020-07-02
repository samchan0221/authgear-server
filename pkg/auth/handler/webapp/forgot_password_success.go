package webapp

import (
	"net/http"

	"github.com/authgear/authgear-server/pkg/db"
	"github.com/authgear/authgear-server/pkg/httproute"
)

func ConfigureForgotPasswordSuccessRoute(route httproute.Route) httproute.Route {
	return route.
		WithMethods("OPTIONS", "GET").
		WithPathPattern("/forgot_password/success")
}

type ForgotPasswordSuccessProvider interface {
	GetForgotPasswordSuccess(w http.ResponseWriter, r *http.Request) (func(error), error)
}

type ForgotPasswordSuccessHandler struct {
	Provider ForgotPasswordSuccessProvider
	Database *db.Handle
}

func (h *ForgotPasswordSuccessHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.Database.WithTx(func() error {
		if r.Method == "GET" {
			writeResponse, err := h.Provider.GetForgotPasswordSuccess(w, r)
			writeResponse(err)
			return err
		}
		return nil
	})
}
