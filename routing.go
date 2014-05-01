package dorsey

import (
	"strings"
)

// HandlerFunc is a function that can render a response to a request.
type HandlerFunc func(*ResponseWriter, *Request)

type routeTable struct {
	pathMatchers []pathMatcher
	handlerFuncs [][]HandlerFunc
	methods      []string
}

type pathMatcher interface {
	match(r *Request) bool
	extractParams(r *Request)
}

type stringMatcher struct {
	value string
	parts []string
}

func (s stringMatcher) match(r *Request) bool {
	if len(r.pathParts) < len(s.parts) {
		return false
	}

	for index, part := range s.parts {
		if !strings.HasPrefix(part, ":") && r.pathParts[index] != part && part != "" {
			return false
		}
	}

	return true
}

func (s stringMatcher) extractParams(r *Request) {
	r.matchedPrefix = strings.Join(r.pathParts[:len(s.parts)], "/")
	for index, part := range s.parts {
		if strings.HasPrefix(part, ":") && len(r.pathParts[index]) > 1 {
			r.URLParams[part[1:]] = r.pathParts[index]
		}
	}
}

func makeStringMatcher(path string) stringMatcher {
	return stringMatcher{value: path, parts: strings.Split(path, "/")}
}
