package main

import (
	"encoding/json"
	"html/template"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

	"github.com/julienschmidt/httprouter"
)

var t *template.Template

// func main being used by appengine

func init() {
	r := httprouter.New()
	http.Handle("/", r)
	r.GET("/", indexHandler)
	r.GET("/form/login", loginHandler)
	r.GET("/form/signup", signupHandler)
	r.GET("/userstatus", userStatusHandler)
	r.GET("/items/additem", addItemHandler)
	r.GET("/items/view/:ID", detailHandler)
	r.GET("/items", viewAllHandler)
	r.GET("/user/wishlist", wishListHandler)
	r.GET("/matchkins", gameHandler)
	r.GET("/user/friends", friendsHandler)

	r.POST("/api/userAddItem", userAddItem)
	r.POST("/api/userAddItemIwant", userAddWishListItem)
	r.POST("/api/editItem", createItem)
	r.POST("/api/checkusername", checkUserName)
	r.POST("/api/createuser", createUser)
	r.POST("/api/login", loginProcess)
	r.GET("/api/logout", logout)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("public/"))))

	t = template.Must(template.ParseGlob("templates/html/*.html"))
}

func indexHandler(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(req)

	//Get User Infomation
	memItem, err := getSession(req)
	if err != nil {
		log.Infof(ctx, "User request status from not authenticated source.")
		http.Redirect(res, req, "/form/login", 303)
		return
	}

	var sd SessionData
	if err == nil {
		json.Unmarshal(memItem.Value, &sd)
		sd.LoggedIn = true

		log.Infof(ctx, "Session Data: wheres my", sd.Collection)
	}
	user, err := GetUserByEmail(ctx, sd.Email)
	log.Infof(ctx, "Test: ", user.Collection)
	if err != nil {
		return
	}
	sd.Items = user.Collection
	t.ExecuteTemplate(res, "index.html", &sd)
}

func wishListHandler(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(req)

	//Get User Infomation
	memItem, err := getSession(req)
	if err != nil {
		log.Infof(ctx, "User request status from not authenticated source.")
		http.Redirect(res, req, "/form/login", 303)
		return
	}

	var sd SessionData
	if err == nil {
		json.Unmarshal(memItem.Value, &sd)
		sd.LoggedIn = true

		log.Infof(ctx, "Session Data: wheres my", sd.IWantCol)
	}
	user, err := GetUserByEmail(ctx, sd.Email)
	log.Infof(ctx, "Test: ", user.IWantCol)
	if err != nil {
		return
	}
	sd.Items = user.IWantCol
	t.ExecuteTemplate(res, "my_wishlist.html", &sd)
}

func loginHandler(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	serveTemplate(res, req, "login_screen.html")
}

func friendsHandler(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	serveTemplate(res, req, "not_implemented.html")
}

func signupHandler(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	serveTemplate(res, req, "signup.html")
}

func gameHandler(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	serveTemplate(res, req, "matchkins_game.html")
}

func userStatusHandler(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(req)

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
		log.Infof(ctx, "Session Data:", sd)
	}

	t.ExecuteTemplate(res, "userStatus.html", &sd)
}

//Item Handlers
func addItemHandler(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	serveTemplate(res, req, "itemEdit.html")
}

func detailHandler(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	ctx := appengine.NewContext(req)
	item, err := GetItemFromRequest(req, ps)

	if err != nil {
		log.Errorf(ctx, "DetialHandler error: %v", err)
	}

	t.ExecuteTemplate(res, "item.html", &item)
}

func viewAllHandler(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	ctx := appengine.NewContext(req)
	items, err := GetItems(ctx)
	if err != nil {
		return
	}
	//Get User Infomation
	memItem, err := getSession(req)
	if err != nil {
		log.Infof(ctx, "User not logged in viewing all items.")

	}

	var sd SessionData
	if err == nil {
		json.Unmarshal(memItem.Value, &sd)
		sd.LoggedIn = true
	}

	sd.Items = items
	log.Infof(ctx, "Test: %v", items)
	t.ExecuteTemplate(res, "items.html", &sd)
}
