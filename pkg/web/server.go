package web

import (
	"net/http"

	"github.com/go-board/x-go/xctx"
	"github.com/go-board/x-go/xnet/xhttp"
	"github.com/julienschmidt/httprouter"
)

type Server struct {
	router      *httprouter.Router
	middlewares []xhttp.Middleware
}

func New(middlewares ...xhttp.Middleware) *Server {
	return &Server{
		router:      httprouter.New(),
		middlewares: middlewares,
	}
}

func (s *Server) Run(addr string) error {
	return http.ListenAndServe(addr, s.router)
}

func (s *Server) Group(path string, middlewares ...xhttp.Middleware) *Route {
	newMiddlewares := make([]xhttp.Middleware, len(s.middlewares)+len(middlewares))
	copy(newMiddlewares, s.middlewares)
	copy(newMiddlewares[len(s.middlewares):], middlewares)
	return &Route{
		s:           s,
		path:        path,
		middlewares: newMiddlewares,
	}
}

func (s *Server) Handle(method string, path string, h http.Handler, middlewares ...xhttp.Middleware) {
	s.router.Handle(method, path, func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		ctx := xctx.NewTyped(request.Context())
		ctx.With(params)
		request = request.WithContext(ctx)
		h = xhttp.ComposeMiddleware(h, append(s.middlewares, middlewares...)...)
		h.ServeHTTP(writer, request)
	})
}

func (s *Server) Get(path string, h http.Handler, middlewares ...xhttp.Middleware) {
	s.Handle(http.MethodGet, path, h, middlewares...)
}

func (s *Server) Post(path string, h http.Handler, middlewares ...xhttp.Middleware) {
	s.Handle(http.MethodPost, path, h, middlewares...)
}
func (s *Server) Put(path string, h http.Handler, middlewares ...xhttp.Middleware) {
	s.Handle(http.MethodPut, path, h, middlewares...)
}
func (s *Server) Delete(path string, h http.Handler, middlewares ...xhttp.Middleware) {
	s.Handle(http.MethodDelete, path, h, middlewares...)
}
func (s *Server) Patch(path string, h http.Handler, middlewares ...xhttp.Middleware) {
	s.Handle(http.MethodPatch, path, h, middlewares...)
}
func (s *Server) Head(path string, h http.Handler, middlewares ...xhttp.Middleware) {
	s.Handle(http.MethodHead, path, h, middlewares...)
}
func (s *Server) Options(path string, h http.Handler, middlewares ...xhttp.Middleware) {
	s.Handle(http.MethodOptions, path, h, middlewares...)
}
func (s *Server) Trace(path string, h http.Handler, middlewares ...xhttp.Middleware) {
	s.Handle(http.MethodTrace, path, h, middlewares...)
}
func (s *Server) Connect(path string, h http.Handler, middlewares ...xhttp.Middleware) {
	s.Handle(http.MethodConnect, path, h, middlewares...)
}

type Route struct {
	s           *Server
	path        string
	middlewares []xhttp.Middleware
}

func (r *Route) Group(path string, middlewares ...xhttp.Middleware) *Route {
	newMiddlewares := make([]xhttp.Middleware, len(r.middlewares)+len(middlewares))
	copy(newMiddlewares, r.middlewares)
	copy(newMiddlewares[len(r.middlewares):], middlewares)
	return &Route{
		s:           r.s,
		path:        r.path + path,
		middlewares: newMiddlewares,
	}
}

func (r *Route) Handle(method string, path string, h http.Handler, middlewares ...xhttp.Middleware) {
	r.s.Handle(method, path, h, middlewares...)
}

func (r *Route) Get(path string, h http.Handler, middlewares ...xhttp.Middleware) {
	r.s.Handle(http.MethodGet, path, h, middlewares...)
}

func (r *Route) Post(path string, h http.Handler, middlewares ...xhttp.Middleware) {
	r.s.Handle(http.MethodPost, path, h, middlewares...)
}

func (r *Route) Put(path string, h http.Handler, middlewares ...xhttp.Middleware) {
	r.s.Handle(http.MethodPut, path, h, middlewares...)
}

func (r *Route) Delete(path string, h http.Handler, middlewares ...xhttp.Middleware) {
	r.s.Handle(http.MethodDelete, path, h, middlewares...)
}

func (r *Route) Patch(path string, h http.Handler, middlewares ...xhttp.Middleware) {
	r.s.Handle(http.MethodPatch, path, h, middlewares...)
}

func (r *Route) Head(path string, h http.Handler, middlewares ...xhttp.Middleware) {
	r.s.Handle(http.MethodHead, path, h, middlewares...)
}

func (r *Route) Options(path string, h http.Handler, middlewares ...xhttp.Middleware) {
	r.s.Handle(http.MethodOptions, path, h, middlewares...)
}

func (r *Route) Trace(path string, h http.Handler, middlewares ...xhttp.Middleware) {
	r.s.Handle(http.MethodTrace, path, h, middlewares...)
}

func (r *Route) Connect(path string, h http.Handler, middlewares ...xhttp.Middleware) {
	r.s.Handle(http.MethodConnect, path, h, middlewares...)
}
