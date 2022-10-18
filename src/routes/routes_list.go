package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/forgeronvirtuel/cache-server/src/cache"
	"io"
	"log"
	"net/http"
	"strings"
)

type HttpError struct {
	Status  int               `json:"status"`
	Message string            `json:"message,omitempty"`
	Details map[string]string `json:"details,omitempty"`
}

type CacheServerHttpHandler func(w http.ResponseWriter, r *http.Request, route *Route) (status int, body []byte, err *HttpError)

type Route struct {
	Path string
	// TODO : better implementation bc a map here is slow
	MethodsHandler map[string]CacheServerHttpHandler
	Cache          *cache.LockedCache
	Logger         *log.Logger
	KeyHandler     *cache.Root
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
	if _, err := w.Write(body); err != nil {
		panic(err)
	}
}

func write404NotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}

func CreateRouteList(logger *log.Logger) []*Route {
	cche := cache.NewLockedCache(nil)
	keyhandler := cache.NewRoot()
	return []*Route{
		{
			Path: "/add",
			MethodsHandler: map[string]CacheServerHttpHandler{
				http.MethodPost: PostAddValueHandler,
			},
			Logger:     logger,
			Cache:      cche,
			KeyHandler: keyhandler,
		},
		{
			Path: "/get",
			MethodsHandler: map[string]CacheServerHttpHandler{
				http.MethodGet: GetValueHandler,
			},
			Logger:     logger,
			Cache:      cche,
			KeyHandler: keyhandler,
		},
		{
			Path: "/search",
			MethodsHandler: map[string]CacheServerHttpHandler{
				http.MethodGet: GetSearchHandler,
			},
			Logger:     logger,
			Cache:      cche,
			KeyHandler: keyhandler,
		},
	}
}

// GetValueHandler retrieves a value from the cache and send it.
func GetValueHandler(_ http.ResponseWriter, r *http.Request, route *Route) (status int, body []byte, err *HttpError) {
	key := strings.Replace(r.URL.String(), route.Path, "", 1)
	value, ok := route.Cache.GetWithStatus(key)
	if !ok {
		return http.StatusNotFound, nil, nil
	}
	route.Logger.Printf("Key = `%s`, value = `%s`", key, string(value))
	return http.StatusOK, value, nil
}

// PostAddValueHandler add a value to the cache.
func PostAddValueHandler(_ http.ResponseWriter, r *http.Request, route *Route) (status int, body []byte, err *HttpError) {
	key := strings.Replace(r.URL.String(), route.Path, "", 1)
	buf := bytes.NewBuffer(nil)
	_, ioerr := io.Copy(buf, r.Body)
	if ioerr != nil {
		return http.StatusBadRequest, nil, &HttpError{
			Status:  http.StatusBadRequest,
			Message: "Cannot read request's body",
			Details: map[string]string{
				"ioerr": ioerr.Error(),
			},
		}
	}
	route.Logger.Printf("Key = `%s`, value = `%s`", key, string(buf.Bytes()))
	route.Cache.Add(key, buf.Bytes())
	route.KeyHandler.Root(key)
	return http.StatusOK, nil, nil
}

// GetSearchHandler returns the list of keys currently registered
func GetSearchHandler(_ http.ResponseWriter, r *http.Request, route *Route) (int, []byte, *HttpError) {
	paths := route.KeyHandler.GetAllPaths("")
	return http.StatusOK, []byte(strings.Join(paths, "\n")), nil
}
