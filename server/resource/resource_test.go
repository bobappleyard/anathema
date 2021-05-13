package resource

import (
	"context"
	"github.com/bobappleyard/anathema/assert"
	"github.com/bobappleyard/anathema/component/di"
	"github.com/bobappleyard/anathema/server/a"
	"github.com/bobappleyard/anathema/server/router"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestHappyPath(t *testing.T) {
	s := di.GetScope(di.EnterScope(context.Background(), ""))
	s.AddRule(&tagRule{
		encodings: []FieldEncoding{
			new(intEncoding),
			new(stringEncoding),
			new(methodEncoding),
		},
		sources: []TagSource{
			new(path),
			new(get),
			new(head),
		},
	})
	s.AddRule(di.Instance(struct{ a.Resource }{}))

	r, err := router.ParseRoute("/users/{client}/{name}")
	assert.NoError(t, err)

	rq := new(http.Request)
	rq.URL, _ = url.Parse("/?since=2020-02-01T00:00:00Z&count=10")
	rq.Header = map[string][]string{
		"Authorization": {"Bearer abc"},
	}

	m := router.Match{
		Request: rq,
		Route:   r,
		Values:  []string{"1234", "Bob"},
	}
	s.AddRule(di.Instance(m))

	var resource struct {
		a.Resource `path:"/users/{client}/{name}"`

		ClientID int       `path:"client"`
		Name     string    `path:"name"`
		Since    time.Time `get:"since"`
		Start    int       `get:"start"`
		Count    int       `get:"count"`
		Auth     string    `head:"Authorization"`
	}
	err = s.Furnish(&resource)
	assert.NoError(t, err)

	tt, _ := time.Parse(time.RFC3339, "2020-02-01T00:00:00Z")

	assert.Equal(t, resource.ClientID, 1234)
	assert.Equal(t, resource.Name, "Bob")
	assert.Equal(t, resource.Since, tt)
	assert.Equal(t, resource.Start, 0)
	assert.Equal(t, resource.Count, 10)
	assert.Equal(t, resource.Auth, "Bearer abc")
}
