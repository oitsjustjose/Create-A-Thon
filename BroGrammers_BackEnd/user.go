package main

import (
	"encoding/json"

	"io"

	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

type User struct {
	Email      string
	FirstName  string
	LastName   string
	UserName   string
	Password   string `json:"-"`
	Collection []Item
	IWantCol   []Item
}

func (user *User) key(ctx context.Context) *datastore.Key {
	return datastore.NewKey(ctx, "Users", user.Email, 0, nil)
}

func (user *User) save(ctx context.Context) error {
	key, err := datastore.Put(ctx, user.key(ctx), user)
	if err != nil {
		return err
	}
	log.Infof(ctx, "Saved user: j%v", key)
	//No error return nil
	return nil
}

func GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	user.Email = email

	key := user.key(ctx)
	err := datastore.Get(ctx, key, &user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func UpdateUser(ctx context.Context, email string, r io.ReadCloser) (*User, error) {
	user, err := GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	// this is a temporary instance that is built
	// from `r.Body`
	var uservar User
	err = json.NewDecoder(r).Decode(&uservar)
	if err != nil {
		return nil, err
	}

	// we only want specific fields to be updated
	user.Email = uservar.Email

	err = user.save(ctx)
	if err != nil {
		return nil, err
	}

	return user, nil

}
