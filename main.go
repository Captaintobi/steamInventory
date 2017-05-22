package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/solovev/steam_go"
)

type user struct {
	Username string
	Password string
	Steamid  string
	Friends  map[string]uint64
}

const API_KEY = "CD3D14A8D01B5E68C7384C946B3A6631"

var steamID string
var isLoggedIn bool
var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func main() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/signup", signUpPage)
	http.HandleFunc("/login", loginPage)
	http.HandleFunc("/invi", getInventory)
	http.Handle("/favicon.io", http.NotFoundHandler())

	http.ListenAndServe(":8080", nil)
}
func loginPage(w http.ResponseWriter, r *http.Request) {
	opID := steam_go.NewOpenId(r)
	var err error
	switch opID.Mode() {
	case "":
		http.Redirect(w, r, opID.AuthUrl(), 301)
	case "cancel":
		w.Write([]byte("Authorization cancelled"))
	default:
		steamID, err = opID.ValidateAndGetId()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		// Do whatever you want with steam id
	}
	http.Redirect(w, r, "/signup", http.StatusMovedPermanently)
}
func homePage(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "index.gohtml", "HElla")

}
func getInventory(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("http://steamcommunity.com/inventory/76561198035405427/730/2?l=english&count=5000")
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	var inventory Inventory
	json.NewDecoder(resp.Body).Decode(&inventory)
	for _, itemName := range inventory.Descriptions {
		fmt.Println(itemName.MarketHashName)
	}
	tpl.ExecuteTemplate(w, "inventory.gohtml", inventory.Descriptions)
}

func signUpPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.ServeFile(w, r, "templates/signUp.gohtml")
		return
	}
	getCookie(w, r)
	//The formvalue parameter is the name of the input in the
	//html page
	//TODO: Move to another file
	username := r.FormValue("username")
	password := r.FormValue("password")
	user := user{
		username,
		password,
		steamID,
		nil,
	}
	/*
		steamid, _ := strconv.ParseInt(user.Steamid, 10, 64)
		friends, err := friendToMap(uint64(steamid), "friend", API_KEY)
		if err != nil {
			fmt.Println(err)
		}
		user.Friends = friends
	*/

	tpl.ExecuteTemplate(w, "index.gohtml", user)
}
