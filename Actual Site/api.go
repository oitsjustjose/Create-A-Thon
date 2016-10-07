package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/nu7hatch/gouuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
)

func checkUserName(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(req)
	bs, err := ioutil.ReadAll(req.Body)
	sbs := string(bs)
	log.Infof(ctx, "REQUEST BODY: %v", sbs)
	var user User
	key := datastore.NewKey(ctx, "Users", sbs, 0, nil)
	err = datastore.Get(ctx, key, &user)
	// if there is an err, there is NO user
	log.Infof(ctx, "ERR: %v", err)
	if err != nil {
		// there is an err, there is a NO user
		fmt.Fprint(res, "false")
		return
	} else {
		fmt.Fprint(res, "true")
	}
}

func createUser(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(req)
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.FormValue("password")), bcrypt.DefaultCost)
	if err != nil {
		log.Errorf(ctx, "error creating password: %v", err)
		http.Error(res, err.Error(), 500)
		return
	}
	user := User{
		Email:     req.FormValue("email"),
		UserName:  req.FormValue("username"),
		FirstName: req.FormValue("firstname"),
		LastName:  req.FormValue("lastname"),
		Password:  string(hashedPass),
	}
	err = user.save(ctx)
	if err != nil {
		log.Errorf(ctx, "Error adding user.")
		log.Errorf(ctx, "error adding todo: %v", err)
		http.Error(res, err.Error(), 500)
		return
	}

	createSession(res, req, user)
	// redirect
	http.Redirect(res, req, "/", 302)
}

func loginProcess(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(req)

	user, err := GetUserByEmail(ctx, req.FormValue("email"))
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.FormValue("password"))) != nil {
		// failure logging in
		var sd SessionData
		sd.LoginFail = true
		t.ExecuteTemplate(res, "login_screen.html", sd)
		return
	} else {
		//

		// success logging in
		createSession(res, req, *user)
		// redirect
		http.Redirect(res, req, "/", 302)
	}
}

func createSession(res http.ResponseWriter, req *http.Request, user User) {
	ctx := appengine.NewContext(req)
	// SET COOKIE
	id, _ := uuid.NewV4()
	cookie := &http.Cookie{
		Name:  "session",
		Value: id.String(),
		Path:  "/",
		//		UNCOMMENT WHEN DEPLOYED:
		//		Secure: true,
		//		HttpOnly: true,
	}
	http.SetCookie(res, cookie)

	// SET MEMCACHE session data (sd)
	json, err := json.Marshal(user)
	if err != nil {
		log.Errorf(ctx, "error marshalling during user creation: %v", err)
		http.Error(res, err.Error(), 500)
		return
	}
	sd := memcache.Item{
		Key:   id.String(),
		Value: json,
		//		Expiration: time.Duration(20*time.Minute),
		Expiration: time.Duration(20 * time.Minute),
	}
	memcache.Set(ctx, &sd)
}

func logout(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(req)

	cookie, err := req.Cookie("session")
	// cookie is not set
	if err != nil {
		http.Redirect(res, req, "/", 302)
		return
	}

	// clear memcache
	sd := memcache.Item{
		Key:        cookie.Value,
		Value:      []byte(""),
		Expiration: time.Duration(1 * time.Microsecond),
	}
	memcache.Set(ctx, &sd)

	// clear the cookie
	cookie.MaxAge = -1
	http.SetCookie(res, cookie)

	// redirect
	http.Redirect(res, req, "/", 302)
}

// Item Handling
func createItem(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	//Create Context
	ctx := appengine.NewContext(req)

	//Create Item from Form Values
	item := Item{
		Title:       req.FormValue("title"),
		Description: req.FormValue("description"),
		ImageURL:    req.FormValue("imageURL"),
	}

	//Save item
	err := item.save(ctx)
	if err != nil {
		log.Errorf(ctx, "Could not save Item: %v", err)
	}

	//Redirect to added Item
	http.Redirect(res, req, fmt.Sprintf("/items/view/%d", item.ID), http.StatusFound)

}

func userAddItem(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(req)

	//Find Item
	bs, err := ioutil.ReadAll(req.Body)
	sbs := string(bs)
	log.Infof(ctx, "REQUEST BODY: %v", sbs)
	id, err := strconv.ParseInt(sbs, 10, 64)
	item, err := GetItem(ctx, id)
	if err != nil {
		log.Errorf(ctx, "Could Retrieve Item with ID Call: %v", err)
	}

	//Get User Infomation
	memItem, err := getSession(req)
	if err != nil {
		log.Infof(ctx, "User request status from not authenticated source.")
		http.Error(res, "Not logged in", http.StatusForbidden)
	}

	var sd SessionData
	if err == nil {
		json.Unmarshal(memItem.Value, &sd)
		sd.LoggedIn = true
	}

	//Added item to slice.
	user, err := GetUserByEmail(ctx, sd.Email)
	for _, colItm := range user.Collection {
		if colItm.ImageURL == item.ImageURL { //Not the best comparison....
			log.Errorf(ctx, "User already has item. %v", item.ID)
			return
		}
	}
	user.Collection = append(user.Collection, *item)
	if err != nil {
		log.Errorf(ctx, "No User by that email: %v", err)
		return
	}
	user.save(ctx)

	if err != nil {
		log.Errorf(ctx, "Error updating user: %v", err)
		return
	}
	fmt.Fprint(res, "true")
}

func userAddWishListItem(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(req)

	//Find Item
	bs, err := ioutil.ReadAll(req.Body)
	sbs := string(bs)
	log.Infof(ctx, "REQUEST BODY: %v", sbs)
	id, err := strconv.ParseInt(sbs, 10, 64)
	item, err := GetItem(ctx, id)
	if err != nil {
		log.Errorf(ctx, "Could Retrieve Item with ID Call: %v", err)
	}

	//Get User Infomation
	memItem, err := getSession(req)
	if err != nil {
		log.Infof(ctx, "User request status from not authenticated source.")
		http.Error(res, "Not logged in", http.StatusForbidden)
	}

	var sd SessionData
	if err == nil {
		json.Unmarshal(memItem.Value, &sd)
		sd.LoggedIn = true
	}

	//Added item to slice.
	user, err := GetUserByEmail(ctx, sd.Email)
	for _, colItm := range user.IWantCol {
		if colItm.ImageURL == item.ImageURL { //Need to find a better comparison, Maybe actual save IDs to the data store.
			log.Errorf(ctx, "User already has item. %v", item.ID)
			return
		}
	}
	user.IWantCol = append(user.IWantCol, *item)
	if err != nil {
		log.Errorf(ctx, "No User by that email: %v", err)
		return
	}
	user.save(ctx)

	if err != nil {
		log.Errorf(ctx, "Error updating user: %v", err)
		return
	}
	fmt.Fprint(res, "true")
}
