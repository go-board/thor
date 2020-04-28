package web

import (
	"net/http"

	"github.com/go-board/x-go/xctx"
	"github.com/go-board/x-go/xnet/xhttp"
)

type AuthResult interface {
	AppId() string
	DeviceId() string
	UserId() *string
}

type Authenticator interface {
	Authenticate(r *http.Request) (authResult AuthResult, need bool, err error)
}

func AuthenticateMiddleware(a Authenticator) xhttp.Middleware {
	return xhttp.MiddlewareFn(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			authResult, need, err := a.Authenticate(request)
			if err != nil && need {
				writer.WriteHeader(http.StatusUnauthorized)
				return
			}
			h.ServeHTTP(writer, injectDataToRequest(request, authResult))
		})
	})
}

func injectDataToRequest(r *http.Request, data interface{}) *http.Request {
	ctx := xctx.NewTyped(r.Context())
	ctx.With(data)
	return r.WithContext(ctx)
}
