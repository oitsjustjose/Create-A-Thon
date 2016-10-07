package main

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

//http://stevenlu.com/posts/2015/03/23/google-datastore-with-golang/ datastore information
//Item sturct used for items models
type Item struct {
	ID          int64 `datastore:"-"`
	Title       string
	Description string
	ImageURL    string
}

func (item *Item) key(ctx context.Context) *datastore.Key {

	//Datastore to determine the ID
	if item.ID == 0 {
		return datastore.NewIncompleteKey(ctx, "Items", nil)
	}

	return datastore.NewKey(ctx, "Items", "", item.ID, nil)
}

func (item *Item) save(ctx context.Context) error {

	key, err := datastore.Put(ctx, item.key(ctx), item)
	if err != nil {
		return err
	}

	//Datastore does not add the ID to struct
	item.ID = key.IntID()
	return nil

}

func GetItem(ctx context.Context, id int64) (*Item, error) {
	var item Item
	item.ID = id

	key := item.key(ctx)
	err := datastore.Get(ctx, key, &item)

	if err != nil {
		return nil, err
	}

	item.ID = key.IntID()

	return &item, nil
}

func GetItemFromRequest(req *http.Request, ps httprouter.Params) (*Item, error) {
	ctx := appengine.NewContext(req)
	id, err := strconv.ParseInt(ps.ByName("ID"), 10, 64)
	if err != nil {
		return nil, err
	}

	item, err := GetItem(ctx, id)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func GetItems(ctx context.Context) ([]Item, error) {
	q := datastore.NewQuery("Items").Order("Title")

	var items []Item
	keys, err := q.GetAll(ctx, &items)
	if err != nil {
		return nil, err
	}

	// you'll see this a lot because instances
	// do not have this by default
	for i := 0; i < len(items); i++ {
		items[i].ID = keys[i].IntID()
	}

	return items, nil
}
