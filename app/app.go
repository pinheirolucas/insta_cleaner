package app

import (
	"net/http"
)

func init() {
	http.HandleFunc("/tasks/unfollow", unfollow)
	http.Handle("/auth/login", post(http.HandlerFunc(login)))
	http.Handle("/users/following", get(auth(http.HandlerFunc(following))))
}
