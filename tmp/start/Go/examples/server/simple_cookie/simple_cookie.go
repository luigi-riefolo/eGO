package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
)

// User ...
type User struct {
	ID   int
	Name string
}

// UserSession ...
var UserSession Session

var sessionStore = sessions.NewCookieStore([]byte("our-social-network-application"))

// LoginPOST ...
func LoginPOST(w http.ResponseWriter, r *http.Request) {
	validateSession(w, r)

	u := User{}
	name := r.FormValue("user_name")
	pass := r.FormValue("user_password")
	password := weakPasswordHash(pass)
	err := database.QueryRow("SELECT user_id, user_name FROM users WHERE user_name=? and user_password=?", name,
		password).Scan(&u.ID, &u.Name)
	if err != nil {
		fmt.Fprintln(w, err)
		u.ID = 0
		u.Name = ""
	} else {
		updateSession(UserSession.ID, u.ID)
		fmt.Fprintln(w, u.Name)
	}
}

func getSessionUID(sid string) int {
	user := User{}
	err := database.QueryRow("SELECT user_id FROM sessions WHERE session_id=?", sid).Scan(user.ID)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return user.ID
}

func updateSession(sid string, uid int) {
	const timeFmt = "2006-01-02T15:04:05.999999999"
	tstamp := time.Now().Format(timeFmt)
	_, err := database.Exec("INSERT INTO sessions SET session_id=?,user_id=?, session_update=? ON DUPLICATE KEY UPDATE user_id=?,session_update=?", sid, uid, tstamp, uid, tstamp)
	if err != nil {
		fmt.Println(err)
	}
}

func generateSessionID() string {
	sid := make([]byte, 24)
	_, err := io.ReadFull(rand.Reader, sid)
	if err != nil {
		log.Fatal("Could not generate session id")
	}
	return base64.URLEncoding.EncodeToString(sid)
}

func cookieHandler(w http.ResponseWriter, r *http.Request) {
	var cookieStore = sessions.NewCookieStore([]byte("ideally, some random piece of entropy"))
	session, err := cookieStore.Get(r, "mystore")
	if err != nil {
		log.Fatalln(err)
	}
	if value, exists := session.Values["hello"]; exists {
		fmt.Fprintln(w, value)
	} else {
		session.Values["hello"] = "(world)"
		session.Save(r, w)
		fmt.Fprintln(w, "We just set the value!")
	}
}

func validateSession(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, "app-session")
	if err != nil {
		log.Fatalln(err)
	}
	if sid, valid := session.Values["sid"]; valid {
		currentUID := getSessionUID(sid.(string))
		updateSession(sid.(string), currentUID)
		UserSession.ID = string(currentUID)
	} else {
		newSID := generateSessionID()
		session.Values["sid"] = newSID
		session.Save(r, w)
		UserSession.ID = newSID
		updateSession(newSID, 0)
	}
	fmt.Println(session.ID)
}

func main() {
	http.HandleFunc("/test", cookieHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
