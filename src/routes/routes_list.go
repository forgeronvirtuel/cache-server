package routes

import (
	"encoding/json"
	"fmt"
	"github.com/forgeronvirtuel/cache-server/src/cache"
	"log"
	"net/http"
)

type HttpError struct {
	Status  int               `json:"status"`
	Message string            `json:"message,omitempty"`
	Details map[string]string `json:"details,omitempty"`
}

type CacheServerHttpHandler func(w http.ResponseWriter, r *http.Request, route *Route) (status int, body interface{}, err *HttpError)

type Route struct {
	Path           string
	MethodsHandler map[string]CacheServerHttpHandler
	Cache          *cache.LockedCache
	Logger         *log.Logger
}

func (route *Route) HandleHttp(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if p := recover(); p != nil {
			data, err := json.Marshal(HttpError{
				Status:  500,
				Message: "panic detected",
				Details: map[string]string{
					"err": fmt.Sprintf("%s", p),
				},
			})
			if err != nil {
				panic(err)
			}
			w.WriteHeader(http.StatusInternalServerError)
			if _, err := w.Write(data); err != nil {
				panic(err)
			}
			route.Logger.Println(p)
		}
	}()

	route.Logger.Printf("Route called: %s. [%s] URL: %s", route.Path, r.Method, r.URL)

	handler, ok := route.MethodsHandler[r.Method]
	if !ok {
		route.Logger.Println("No handler found")
		write404NotFound(w)
		return
	}

	st, body, httpErr := handler(w, r, route)
	w.WriteHeader(st)
	if httpErr != nil {

		// Encode & write error comming from the route
		data, err := json.Marshal(httpErr)

		// If an error happens while we write the error => panic
		if err != nil {
			panic(err)
		}
		if _, err := w.Write(data); err != nil {
			panic(err)
		}
		return
	}

	// Encode and write body generated from the route
	data, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	if _, err := w.Write(data); err != nil {
		panic(err)
	}
}

func write404NotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}

func CreateRouteList(logger *log.Logger) []*Route {
	cache := cache.NewLockedCache(nil)
	return []*Route{
		{
			Path: "/add/",
			MethodsHandler: map[string]CacheServerHttpHandler{
				http.MethodPost: PostAddValueHandler,
			},
			Logger: logger,
			Cache:  cache,
		},
		{
			Path: "/get/",
			MethodsHandler: map[string]CacheServerHttpHandler{
				http.MethodGet: GetValueHandler,
			},
			Logger: logger,
			Cache:  cache,
		},
	}
}

// GetValueHandler retrieves a value from the cache and send it.
func GetValueHandler(w http.ResponseWriter, r *http.Request, route *Route) (status int, body interface{}, err *HttpError) {
	//return http.StatusOK, "my_value", nil

	return http.StatusBadRequest, "my_value", &HttpError{
		Status:  http.StatusBadRequest,
		Message: "Internal server error.",
	}
}

// GetValueHandler retrieves a value from the cache and send it.
func PostAddValueHandler(w http.ResponseWriter, r *http.Request, route *Route) (status int, body interface{}, err *HttpError) {
	//return http.StatusOK, "my_value", nil

	return http.StatusOK, "my_value", &HttpError{
		Status:  500,
		Message: "Internal server error.",
	}
}
