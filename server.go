package dorsey

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

// Params is used to keep track of parameters that may be passed in via
// named portions of the URL path.
type Params map[string]string

// A wrapper for the http.Request type with some additional information.
type Request struct {
	*http.Request
	pathParts []string
	// Params passed in via named parts of the URL
	UrlParams     Params
	matchedPrefix string
}

// Get a named parameter from the URL.  For instance, if the handler is added at "/a/:blah/portion",
// then getting a requests url parameter named "blah" from /a/something/portion would return
// "something"
func (r *Request) GetURLParam(name string) string {
	if value, found := r.UrlParams[name]; found {
		return value
	}
	return ""
}

// Get a PUT / POST / GET param.  Pass in a default which will be returned if the param
// is not set.
func (r *Request) GetParam(name, def string) string {
	if r := r.FormValue(name); r != "" {
		return r
	}
	return def
}

// Get a HEADER value passed in.
func (r *Request) GetHeader(name string) string {
	return r.Request.Header.Get(name)
}

// A wrapper for http.ResponseWriter with some additiona information.
type ResponseWriter struct {
	http.ResponseWriter
	*http.Request
	rendered bool
	// Pass context variables on to future handlers
	Context map[string]interface{}
}

// Render a value.  This should be called once and only once.  The value can be either
// a string or a byte array.  The content type should be set automagically.
func (w *ResponseWriter) Render(value interface{}) {
	if w.rendered {
		log.Panicln("Render / Redirect functions can only be called at most once.")
	}

	switch value.(type) {
	case string:
		fmt.Fprint(w.ResponseWriter, value)
	case []byte:
		fmt.Fprint(w.ResponseWriter, value)
	default:
		log.Panicln("No idea how to render passsed value:", value)
	}
	w.rendered = true
}

// Set the given header.
func (w *ResponseWriter) SetHeader(name, value string) {
	w.ResponseWriter.Header().Set(name, value)
}

// Set the response code (200 by default).
func (w *ResponseWriter) SetResponseCode(code int) {
	w.ResponseWriter.WriteHeader(code)
}

// Temporary redirect
func (w *ResponseWriter) Redirect(path string) {
	if w.rendered {
		log.Panicln("Render / Redirect functions can only be called at most once.")
	}
	http.Redirect(w.ResponseWriter, w.Request, path, http.StatusTemporaryRedirect)
	w.rendered = true
}

// Return a 401 and request a username / password
func (w *ResponseWriter) Unauthorized() {
	w.SetHeader("WWW-Authenticate", "Basic realm=\""+w.Request.Header.Get("Host")+"\"")
	w.Error("Not Authorized", 401)
}

// Render an error.
func (w *ResponseWriter) Error(error string, code int) {
	w.SetHeader("Content-Type", "text/html; charset=utf-8")
	w.SetResponseCode(code)
	w.Render(fmt.Sprintf("<html><body><h1>HTTP %d Error</h1></body></html>", code))
	log.Println("ERROR: " + error)
}

// Render a 500 internal server error.
func (w *ResponseWriter) InternalError(error string) {
	w.Error(error, http.StatusInternalServerError)
}

// Render a file.
func (w *ResponseWriter) RenderFile(value string) {
	contents, error := ioutil.ReadFile(value)
	if error != nil {
		w.InternalError("Cannot read \"" + value + "\": " + error.Error())
	} else {
		w.Render(string(contents))
	}
}

// Acts as a http.Handler for use by the default http.Server type
type Server struct {
	routes routeTable
}

// Create a http.Server and set it's handler with the given Server.
func (s *Server) Run(hostport string) error {
	hs := &http.Server{
		Addr:           hostport,
		Handler:        s,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Println("Starting server on", hostport)
	return hs.ListenAndServe()
}

// Add the given HandlerFunc at the set path to handle GET requests
func (s *Server) Get(path interface{}, hs ...HandlerFunc) {
	s.AddRoute("GET", path, hs)
}

// Add the given HandlerFunc at the set path to handle POST requests
func (s *Server) Post(path interface{}, hs ...HandlerFunc) {
	s.AddRoute("POST", path, hs)
}

// Add the given HandlerFunc at the set path to handle PUT requests
func (s *Server) Put(path interface{}, hs ...HandlerFunc) {
	s.AddRoute("PUT", path, hs)
}

// Add the given HandlerFunc at the set path to handle DELETE requests
func (s *Server) Delete(path interface{}, hs ...HandlerFunc) {
	s.AddRoute("DELETE", path, hs)
}

// If you have a custom route you'd like to add (for instance, with a custom http method),
// you can use:
//
//        s.AddRoute("XGET", "/some/path", aHandlerFunc)
//
func (s *Server) AddRoute(method string, path interface{}, hs []HandlerFunc) {
	switch path.(type) {
	case string:
		s.routes.pathMatchers = append(s.routes.pathMatchers, makeStringMatcher(path.(string)))
		s.routes.handlerFuncs = append(s.routes.handlerFuncs, hs)
		s.routes.methods = append(s.routes.methods, method)
	default:
		log.Fatalf("Unknown path type %v", path)
	}
}

// Primary handler function.  All routing is done here.
func (h *Server) ServeHTTP(hw http.ResponseWriter, hr *http.Request) {
	log.Println("Handling request for", hr.URL.Path)

	r := &Request{
		Request:   hr,
		pathParts: strings.Split(hr.URL.Path, "/"),
		UrlParams: make(Params),
	}

	rw := &ResponseWriter{
		ResponseWriter: hw,
		Request:        hr,
		rendered:       false,
		Context:        make(map[string]interface{}),
	}

	for index, handlerfuncs := range h.routes.handlerFuncs {
		if h.routes.methods[index] == hr.Method && h.routes.pathMatchers[index].match(r) {
			h.routes.pathMatchers[index].extractParams(r)
			for _, handlerfunc := range handlerfuncs {
				handlerfunc(rw, r)
				if rw.rendered {
					return
				}
			}
			if !rw.rendered {
				rw.InternalError("Render never called")
			}
		}
	}
	rw.Error("File not found: "+r.URL.String(), http.StatusNotFound)
}

// Create a new Server
func New() *Server {
	return &Server{}
}
