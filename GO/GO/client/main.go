package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
)

type Item struct {
	Title string
	Body  string
}

type Credentials struct {
	Username string
	Password string
}

func main() {
	var reply Item
	var db []Item

	client, err := rpc.DialHTTP("tcp", ":4040")
	if err != nil {
		log.Fatal("Connection error: ", err)
	}
	defer client.Close()

	// Login Menu
	fmt.Println("Login:")
	fmt.Print("Enter Username: ")
	var username string
	fmt.Scanln(&username)

	fmt.Print("Enter Passkey: ")
	var password string
	fmt.Scanln(&password)

	// Validate credentials
	credentials := Credentials{username, password}
	var auth bool
	err = client.Call("API.Authenticate", credentials, &auth)
	if err != nil {
		log.Fatal("Authentication error: ", err)
	}

	if !auth {
		fmt.Println("Invalid credentials. Exiting...")
		return
	}

	// Main Menu
	for {
		var choice int
		fmt.Println("Choose an option:")
		fmt.Println("1. Add Item")
		fmt.Println("2. Delete Item")
		fmt.Println("3. Edit Item")
		fmt.Println("4. Display Database")
		fmt.Println("5. Exit")
		fmt.Print("Enter your choice: ")
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			var title, body string
			fmt.Print("Enter Title: ")
			fmt.Scanln(&title)
			fmt.Print("Enter Body: ")
			fmt.Scanln(&body)
			item := Item{title, body}
			client.Call("API.AddItem", item, &reply)
		case 2:
			var title, body string
			fmt.Print("Enter Title: ")
			fmt.Scanln(&title)
			fmt.Print("Enter Body: ")
			fmt.Scanln(&body)
			item := Item{title, body}
			client.Call("API.DeleteItem", item, &reply)
		case 3:
			var title, body string
			fmt.Print("Enter Title: ")
			fmt.Scanln(&title)
			fmt.Print("Enter New Body: ")
			fmt.Scanln(&body)
			item := Item{title, body}
			client.Call("API.EditItem", item, &reply)
		case 4:
			client.Call("API.GetDB", "", &db)
			fmt.Println("Database: ", db)
		case 5:
			fmt.Println("Exiting...")
			os.Exit(0)
		default:
			fmt.Println("Invalid choice")
		}

		fmt.Println("Operation completed successfully")
	}
}
