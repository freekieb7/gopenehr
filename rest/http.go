package rest

import (
	"log"
	"net/http"
)

type HandleFunc func(w http.ResponseWriter, r *http.Request)
type GroupFunc func(group *Router)
type MiddlewareFunc func(next http.Handler) HandleFunc

func RecoverMiddleware() MiddlewareFunc {
	return func(next http.Handler) HandleFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if recover := recover(); recover != nil {
					log.Println(recover)
					return
				}
			}()

			next.ServeHTTP(w, r)
		}
	}
}

type Route struct {
	Path        string
	Method      string
	Handler     HandleFunc
	Middlewares []MiddlewareFunc
}

type Router struct {
	Path        string
	Routes      []Route
	Middlewares []MiddlewareFunc
	Routers     []Router
}

func NewRouter() Router {
	return Router{
		Path:    "",
		Routes:  make([]Route, 0),
		Routers: make([]Router, 0),
	}
}

func (r *Router) Any(path, method string, handler HandleFunc, middlewares ...MiddlewareFunc) {
	r.Routes = append(r.Routes, Route{
		Path:        path,
		Method:      method,
		Handler:     handler,
		Middlewares: middlewares,
	})
}

func (r *Router) OPTIONS(path string, handler HandleFunc, middlewares ...MiddlewareFunc) {
	r.Any(path, http.MethodOptions, handler)
}

func (r *Router) GET(path string, handler HandleFunc, middlewares ...MiddlewareFunc) {
	r.Any(path, http.MethodGet, handler)
}

func (r *Router) POST(path string, handler HandleFunc, middlewares ...MiddlewareFunc) {
	r.Any(path, http.MethodPost, handler)
}

func (r *Router) PUT(path string, handler HandleFunc, middlewares ...MiddlewareFunc) {
	r.Any(path, http.MethodPut, handler)
}

func (r *Router) DELETE(path string, handler HandleFunc, middlewares ...MiddlewareFunc) {
	r.Any(path, http.MethodDelete, handler)
}

func (r *Router) Group(path string, grouper GroupFunc, middlewares ...MiddlewareFunc) {
	subRouter := NewRouter()
	subRouter.Path = path

	// Adds routes to group
	grouper(&subRouter)

	// Add routes to parent router
	for _, route := range subRouter.Routes {
		subRouter.Any(route.Path, route.Method, route.Handler, middlewares...)
	}

	r.Routers = append(r.Routers, subRouter)
}

type Server struct {
	Name   string
	Router Router
}

func NewServer(name string) Server {
	return Server{
		Name:   name,
		Router: NewRouter(),
	}
}

func (s *Server) ListenAndServe(addr string) error {
	routeTable := mergeRoutes(s.Router)

	for path, routes := range routeTable {
		// Method aware request handler
		http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			for _, route := range routes {
				if r.Method != route.Method {
					continue
				}

				route.Handler(w, r)
				return
			}

			w.WriteHeader(http.StatusMethodNotAllowed)
		})
	}

	return http.ListenAndServe(addr, nil)
}

func mergeRoutes(router Router) map[string][]Route {
	routeTable := make(map[string][]Route, 0)

	// Add direct routes to table
	for _, route := range router.Routes {
		// Create new route with combined router and route configs
		path := router.Path + route.Path
		middlewares := append(router.Middlewares, route.Middlewares...)

		routeTable[path] = append(routeTable[path], Route{
			Path:        path,
			Method:      route.Method,
			Handler:     route.Handler,
			Middlewares: middlewares,
		})
	}

	// Add indirect (sub router) routes to table
	for _, subRouter := range router.Routers {
		subRouterRouteTable := mergeRoutes(subRouter)

		for _, routes := range subRouterRouteTable {
			for _, route := range routes {
				// Create new route with combined router and route configs
				path := router.Path + route.Path
				middlewares := append(router.Middlewares, route.Middlewares...)

				routeTable[path] = append(routeTable[path], Route{
					Path:        path,
					Method:      route.Method,
					Handler:     route.Handler,
					Middlewares: middlewares,
				})
			}
		}
	}

	return routeTable
}
