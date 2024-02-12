package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load(".env")

	http.HandleFunc("/app/", handlerHttp)
	fs := http.FileServer(http.Dir("./static/"))

	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.ListenAndServe(":8080", nil)
}

func handlerHttp(w http.ResponseWriter, r *http.Request) {
	log.Println("In http handler")
	log.Println(r.URL)

	url := r.URL
	if url.Path == "/app/activity" {
		var a Activity
		err := json.NewDecoder(r.Body).Decode(&a)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ProcessActivity(a)
		return
	} else if url.Path == "/app/exchange_token" {
		log.Println("Exchanging token")
		user, err := HandleTokenExchange(url.Query().Get("code"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		token := GenerateUserToken(user.ID)
		cookie := http.Cookie{
			Name:  "user-jwt",
			Value: token,
			Path:  "/",
		}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/static/update-user.html", http.StatusSeeOther)
	} else {
		http.NotFound(w, r)
	}
}
