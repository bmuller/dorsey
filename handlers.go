package dorsey

import (
	"encoding/base64"
	"net/http"
	"strings"
)

func wrapHttpHandler(h http.Handler, stripPrefix bool) HandlerFunc {
	return func(rw *ResponseWriter, r *Request) {
		if stripPrefix {
			http.StripPrefix(r.matchedPrefix, h).ServeHTTP(rw.ResponseWriter, r.Request)
		} else {
			h.ServeHTTP(rw.ResponseWriter, r.Request)
		}
		rw.rendered = true
	}
}

// Serve the directory at the given path.  Note that the URL path this is attached to should end in a '/'.
// For instance:
//        server.Get("/a/path/", dorsey.DirectoryHandler("/tmp/static/files"))
func DirectoryHandler(path string) HandlerFunc {
	return wrapHttpHandler(http.FileServer(http.Dir(path)), true)
}

// Serve the given file.  Obviously, the file must exist.
func FileHandler(path string) HandlerFunc {
	return func(rw *ResponseWriter, r *Request) {
		rw.RenderFile(path)
	}
}

// Permanently redirect to the given url.
func PermanentRedirectHandler(url string) HandlerFunc {
	return wrapHttpHandler(http.RedirectHandler(url, http.StatusMovedPermanently), false)
}

// Temporarily redirect to the given url.
func RedirectHandler(url string) HandlerFunc {
	return wrapHttpHandler(http.RedirectHandler(url, http.StatusTemporaryRedirect), false)
}

// A function that is given a username and password and returns true if the user should
// be authenticated.
type AuthFunction func(user, password string) bool

// Wrap an AuthFunction as a HTTP Basic Auth handler.  If the user isn't authenticated in the given
// AuthFunction, then no further handlers will be called.
func BasicAuthHandler(authFunc AuthFunction) HandlerFunc {
	return func(rw *ResponseWriter, r *Request) {
		as := strings.SplitN(r.GetHeader("Authorization"), " ", 2)
		if len(as) != 2 || as[0] != "Basic" {
			rw.Unauthorized()
		} else if b, err := base64.StdEncoding.DecodeString(as[1]); err != nil {
			rw.Unauthorized()
		} else if pair := strings.SplitN(string(b), ":", 2); len(pair) != 2 || !authFunc(pair[0], pair[1]) {
			rw.Unauthorized()
		}
	}
}
