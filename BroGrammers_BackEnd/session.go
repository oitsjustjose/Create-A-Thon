package main

import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/memcache"
)

type SessionData struct {
	User
	LoggedIn  bool
	LoginFail bool
	Items     []Item
}

func getSession(req *http.Request) (*memcache.Item, error) {
	ctx := appengine.NewContext(req)

	cookie, err := req.Cookie("session")
	if err != nil {
		return &memcache.Item{}, err
	}

	item, err := memcache.Get(ctx, cookie.Value)
	if err != nil {
		return &memcache.Item{}, err
	}
	return item, nil
}
