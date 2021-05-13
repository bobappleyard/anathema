package router

import (
	"errors"
	"net/http"
)

var ErrInvalidSegment = errors.New("invalid path segment")

type Route struct {
	names    []string
	segments []segment
	handler  http.Handler
}

type Match struct {
	Request *http.Request
	Route   *Route
	Values  []string
}

type segment interface {
	match(string) bool
	buildNames([]string) []string
	buildValues([]string, string) []string
}

func ParseRoute(path string) (*Route, error) {
	segments := splitPath(path)
	res := &Route{
		segments: make([]segment, len(segments)),
	}
	for i, seg := range segments {
		s, err := parseSegment(seg)
		if err != nil {
			return nil, err
		}
		res.segments[i] = s
		res.names = s.buildNames(res.names)
	}
	return res, nil
}

func (r *Route) Handler() http.Handler {
	return r.handler
}

func (r *Route) WithHandler(h http.Handler) *Route {
	return &Route{r.names, r.segments, h}
}

func (r *Route) SubRoute(name string) *Route {
	segments := make([]segment, len(r.segments)+1)
	copy(segments, r.segments)
	segments[len(r.segments)] = &fixedSegment{name}
	return &Route{r.names, segments, nil}
}

func (r *Route) Names() []string {
	return r.names
}

func (m Match) GetValue(name string) string {
	for i, n := range m.Route.names {
		if name == n {
			return m.Values[i]
		}
	}
	return ""
}

func (r *Route) match(path []string) (Match, bool) {
	// First try to reject the route testing the constant parts
	for i, p := range path {
		segment := r.segments[i]
		if !segment.match(p) {
			return Match{}, false
		}
	}
	// We have a match, populate a match structure with the remaining parts of
	// the path.
	var vals []string
	for i, p := range path {
		segment := r.segments[i]
		vals = segment.buildValues(vals, p)
	}
	return Match{Route: r, Values: vals}, true
}

func parseSegment(seg string) (segment, error) {
	if seg[0] == '{' {
		if seg[len(seg)-1] != '}' {
			return nil, ErrInvalidSegment
		}
		return &anythingSegment{seg[1 : len(seg)-1]}, nil
	}
	return &fixedSegment{seg}, nil
}

type anythingSegment struct {
	name string
}

func (s *anythingSegment) match(p string) bool {
	return true
}

func (s *anythingSegment) buildNames(names []string) []string {
	return append(names, s.name)
}

func (s *anythingSegment) buildValues(values []string, seg string) []string {
	return append(values, seg)
}

type fixedSegment struct {
	value string
}

func (s *fixedSegment) match(p string) bool {
	return p == s.value
}

func (s *fixedSegment) buildNames(names []string) []string {
	return names
}

func (s *fixedSegment) buildValues(values []string, seg string) []string {
	return values
}
