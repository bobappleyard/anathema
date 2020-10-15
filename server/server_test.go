package server_test

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/bobappleyard/anathema/a"
	"github.com/bobappleyard/anathema/hterror"
	"github.com/bobappleyard/anathema/server"
)

var errNotFound = errors.New("not found")

// Begin by defining our model type. This is what the HTTP API is managing. In
// this example we have an extremely simple user entity.

type User struct {
	ID   int
	Name string
}

// Here we define a repository for storing users. You'll want to define all, or
// most, interactions in terms of these kinds of interfaces.

type UserRepository interface {
	GetUser(id int) (User, error)
	SetUser(id int, u User) error
	DelUser(id int) error
}

// Finally, we define a resource type. This is what really interacts with the
// framework. The important part here is that we have embedded server.Resource
// along with giving it a path field tag.
//
// The path tag is used by the framework to populate the fields of the resource
// struct during method handling. Here, the ID field has been mapped to a path
// segment.

type UserResource struct {
	a.Resource `path:"/users/{ID}"`
	ID         int
}

// The most straightforward way of setting up a resource is to define an Init
// method. This is called on line xxx when the resource type is introduced to
// the server, and maps HTTP verbs to operations on the resource type.

func (UserResource) Init(g server.Group) {
	g.GET(UserResource.GetUser)
	g.PUT(UserResource.PutUser)
	g.DELETE(UserResource.DelUser)
	g.Sub("photo").GET(UserResource.GetPhoto)
}

// Implement those HTTP verbs. Note that the arguments are furnished by the
// framework based on the declared type (typical DI scenario) and the return
// values of the methods are marshaled into appropriate responses (JSON by
// default).

func (r UserResource) GetUser(repo UserRepository) (User, error) {
	u, err := repo.GetUser(r.ID)
	if err == errNotFound {
		return User{}, hterror.WithStatusCode(http.StatusNotFound, err)
	}
	return u, err
}

func (r UserResource) GetPhoto() ([]byte, error) {
	return nil, nil
}

func (r UserResource) PutUser(user User, repo UserRepository) (User, error) {
	err := repo.SetUser(r.ID, user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (r UserResource) DelUser(repo UserRepository) error {
	err := repo.DelUser(r.ID)
	if err == errNotFound {
		return hterror.WithStatusCode(http.StatusNotFound, err)
	}
	return err
}

// This is a mock type for testing out our API.

type testRepository struct {
	users map[int]User
}

func (r testRepository) GetUser(id int) (User, error) {
	if u, ok := r.users[id]; ok {
		return u, nil
	}
	return User{}, errNotFound
}

func (r testRepository) SetUser(id int, user User) error {
	r.users[id] = user
	return nil
}

func (r testRepository) DelUser(id int) error {
	if _, ok := r.users[id]; ok {
		delete(r.users, id)
		return nil
	}
	return errNotFound
}

// This is a function for testing our API by firing a request at it and printing
// what it produces.

func runRequest(s *server.Server, method, path, body string) {
	var bodyReader io.Reader
	if body != "" {
		bodyReader = strings.NewReader(body)
	}
	r := httptest.NewRequest(method, path, bodyReader)
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.ServeHTTP(w, r)
	fmt.Print(w.Code)
	if w.Body.Len() != 0 {
		fmt.Print(" ", w.Body)
	}
	fmt.Println()
}

// Wiring everything together is simple.

func Example() {
	// repo := testRepository{map[int]User{}}

	s := server.New()
	// s.AddService(func() UserRepository { return repo })
	s.Resource(UserResource{})

	// Some test invocations

	runRequest(s, "PUT", "/users/1", `{"ID":1,"Name":"bob"}`)
	runRequest(s, "PUT", "/users/2", `{"ID":2,"Name":"jim"}`)
	runRequest(s, "GET", "/users/1", "")
	runRequest(s, "GET", "/users/3", "")
	runRequest(s, "DELETE", "/users/1", "")
	runRequest(s, "GET", "/users/1", "")
	runRequest(s, "GET", "/users/1/photo", "")

	// Output:
	// 200 {"ID":1,"Name":"bob"}
	// 200 {"ID":2,"Name":"jim"}
	// 200 {"ID":1,"Name":"bob"}
	// 404
	// 204
	// 404
	// 200 null
}
