package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"sync"
)

type Item struct {
	Title string
	Body  string
}

type Credentials struct {
	Password string
}

type API struct {
	mu          sync.Mutex
	Database    []Item                 // Database is now part of the API struct
	initialized bool                   // flag to check if the database is initialized
	itemLocks   map[string]*sync.Mutex // Map to track item locks
}

var initializedOnce sync.Once // Ensures initialization only happens once
var validPassword = "golang"  // Hardcoded valid password

func (a *API) initDatabase() {
	a.Database = make([]Item, 0)               // Initialize the database slice
	a.initialized = true                       // Set the flag to true after initialization
	a.itemLocks = make(map[string]*sync.Mutex) // Initialize the item locks map
}

func (a *API) GetDB(empty string, reply *[]Item) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	*reply = a.Database
	return nil
}

func (a *API) GetByName(title string, reply *Item) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	var getItem Item

	for _, val := range a.Database {
		if val.Title == title {
			getItem = val
			break
		}
	}

	*reply = getItem

	return nil
}

func (a *API) AddItem(item Item, reply *Item) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.Database = append(a.Database, item)
	*reply = item
	return nil
}

func (a *API) EditItem(item Item, reply *Item) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Check if the item lock exists, create if not
	lock, ok := a.itemLocks[item.Title]
	if !ok {
		lock = &sync.Mutex{}
		a.itemLocks[item.Title] = lock
	}

	lock.Lock()
	defer lock.Unlock()

	var found bool

	for idx, val := range a.Database {
		if val.Title == item.Title {
			a.Database[idx] = Item{item.Title, item.Body}
			*reply = a.Database[idx]
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("item not found: %s", item.Title)
	}

	return nil
}

func (a *API) DeleteItem(item Item, reply *Item) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	var found bool

	for idx, val := range a.Database {
		if val.Title == item.Title && val.Body == item.Body {
			found = true
			*reply = a.Database[idx]
			a.Database = append(a.Database[:idx], a.Database[idx+1:]...)
			break
		}
	}

	if !found {
		return fmt.Errorf("item not found: %s", item.Title)
	}

	return nil
}

func (a *API) Authenticate(creds Credentials, reply *bool) error {
	if creds.Password == validPassword {
		*reply = true
	} else {
		*reply = false
	}

	return nil
}

func main() {
	api := new(API)

	initializedOnce.Do(api.initDatabase) // Initialize the database only once

	err := rpc.Register(api)
	if err != nil {
		log.Fatal("error registering API", err)
	}

	rpc.HandleHTTP()

	listener, err := net.Listen("tcp", ":4040")

	if err != nil {
		log.Fatal("Listener error", err)
	}
	log.Printf("serving rpc on port %d", 4040)
	http.Serve(listener, nil)

	if err != nil {
		log.Fatal("error serving: ", err)
	}
}
