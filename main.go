package main

//TODO: Major flaw in the way i was doing things
//Come back when i can get market price for multiple items at once.
//Mainly steam trading cards.
import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/solovev/steam_go"
)

type User struct {
	Username string
	Password string
	Steamid  string
	Friends  map[string]uint64
}

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
	http.HandleFunc("/price", myPrice)
	http.Handle("/favicon.io", http.NotFoundHandler())
	http.ListenAndServe(":9089", nil)
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
func myPrice(w http.ResponseWriter, r *http.Request) {
	newPrice, err := getPrice(730, "AK-47 | Aquamarine Revenge (Minimal Wear)")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("New", newPrice.LowestPrice)
	tpl.ExecuteTemplate(w, "price.gohtml", newPrice)
}
func getInventory(w http.ResponseWriter, r *http.Request) {
	if steamID == "" {
		w.Write([]byte("You have to sign in first"))
	} else {
		//TODO: move to another folder
		resp, err := http.Get("http://steamcommunity.com/inventory/" + steamID + "/730/2?l=english&count=5000")
		if err != nil {
			fmt.Println(err)
		}
		body, err := ioutil.ReadAll(resp.Body)

		var inventory Inventory
		json.NewDecoder(strings.NewReader(string(body))).Decode(&inventory)
		for _, itemName := range inventory.Descriptions {
			fmt.Println(itemName.MarketHashName)
		}
		myPrice, err := getPrice(730, "R8 Revolver | Bone Mask (Field-Tested)")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(myPrice.LowestPrice)
		fmt.Println()

		tpl.ExecuteTemplate(w, "inventory.gohtml", inventory.Descriptions)
	}
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
	user := User{
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
