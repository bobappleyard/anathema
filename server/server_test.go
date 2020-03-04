package server_test

import (
	"errors"
	"fmt"
	"github.com/bobappleyard/anathema/hterror"
	"github.com/bobappleyard/anathema/server"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
)

type User struct {
	ID   int
	Name string
}

type UserRepository interface {
	GetUser(id int) (User, error)
	SetUser(id int, u User) error
}

type UserResource struct {
	server.Resource `path:"/users/{ID}"`
	ID              int
}

func (UserResource) Init(g server.Group) {
	g.GET(UserResource.GetUser)
	g.PUT(UserResource.PutUser)
}

func (r UserResource) GetUser(repo UserRepository) (User, error) {
	return repo.GetUser(r.ID)
}

func (r UserResource) PutUser(user User, repo UserRepository) (User, error) {
	err := repo.SetUser(r.ID, user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

type testRepository struct {
	users map[int]User
}

func (r testRepository) GetUser(id int) (User, error) {
	if u, ok := r.users[id]; ok {
		return u, nil
	}
	return User{}, hterror.WithStatusCode(http.StatusNotFound, errors.New("not found"))
}

func (r testRepository) SetUser(id int, user User) error {
	r.users[id] = user
	return nil
}

func runRequest(s *server.Server, method, path, body string) {
	var bodyReader io.Reader
	if body != "" {
		bodyReader = strings.NewReader(body)
	}
	r := httptest.NewRequest(method, path, bodyReader)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, r)
	fmt.Println(w.Code, w.Body)
}

func Example() {
	repo := testRepository{map[int]User{}}

	s := server.New()
	s.AddService(func() UserRepository { return repo })
	s.Resource(UserResource{})

	runRequest(s, "PUT", "/users/1", `{"ID":1,"Name":"bob"}`)
	runRequest(s, "PUT", "/users/2", `{"ID":2,"Name":"jim"}`)
	runRequest(s, "GET", "/users/1", "")
	runRequest(s, "GET", "/users/3", "")

	// Output:
	// 200 {"ID":1,"Name":"bob"}
	// 200 {"ID":2,"Name":"jim"}
	// 200 {"ID":1,"Name":"bob"}
	// 404
}
