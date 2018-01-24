package main

import (
	"crypto/sha1"
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

const (
	// DBHost ...
	DBHost = "127.0.0.1"
	// DBPort ...
	DBPort = ":3306"
	// DBUser ...
	DBUser = "root"
	// DBPass ...
	DBPass = "password!"
	// DBDbase ...
	DBDbase = "cms"
	port    = ":10443"
	maxLen  = 150
)

var database *sql.DB

// Comment ...
type Comment struct {
	ID          int
	Name        string
	Email       string
	CommentText string
}

// Page ...
type Page struct {
	ID         int
	Title      string
	RawContent string
	Content    template.HTML
	Date       string
	Comments   []Comment
	Session    Session
}

// User ...
type User struct {
	ID   int
	Name string
}

// Session ...
type Session struct {
	ID              string
	Authenticated   bool
	Unauthenticated bool
	User            User
}

// JSONResponse ...
type JSONResponse struct {
	Fields map[string]string
}

// Cookie ...
type Cookie struct {
	Name       string
	Value      string
	Path       string
	Domain     string
	Expires    time.Time
	RawExpires string
	MaxAge     int
	Secure     bool
	HTTPOnly   bool
	Raw        string
	Unparsed   []string
}

// TruncatedText ...
func (p Page) TruncatedText() template.HTML {
	if len(p.RawContent) > maxLen {
		return template.HTML(p.RawContent[:maxLen] + " ...")
	}
	return template.HTML(p.RawContent)
}

// ServeIndex ...
func ServeIndex(w http.ResponseWriter, r *http.Request) {
	//w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")

	log.Print("PROTO!!!!!!!!!!--------------")
	var Pages = []Page{}
	pages, err := database.Query("SELECT page_title,page_content,page_date,page_guid FROM pages ORDER BY ? DESC", "page_date")
	if err != nil {
		fmt.Fprintln(w, err)
	}
	defer pages.Close()
	for pages.Next() {
		thisPage := Page{}
		pages.Scan(&thisPage.Title, &thisPage.RawContent, &thisPage.Date,
			&thisPage.ID)
		fmt.Println("PAGE ID: ", thisPage.ID)
		thisPage.Content = template.HTML(thisPage.RawContent)
		Pages = append(Pages, thisPage)
	}

	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, Pages)
}

// RedirectToHTTPS ...
/*func RedirectToHTTPS(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}
*/

// RedirIndex redirects requests from "/" to "/home".
func RedirIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}

// ServePage ...
func ServePage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println("VARS: ", vars)
	pageGUID := vars["guid"]
	thisPage := Page{}
	fmt.Println(pageGUID)
	err := database.QueryRow("SELECT id,page_title,page_content,page_date FROM pages WHERE page_guid=?", pageGUID).Scan(&thisPage.ID, &thisPage.Title,
		&thisPage.RawContent, &thisPage.Date)
	if err != nil {
		http.Error(w, http.StatusText(404), http.StatusNotFound)
		log.Println(err)
		return
	}
	thisPage.Content = template.HTML(thisPage.RawContent)
	//thisPage.GUID = pageGUID
	/*comments, err := database.Query("SELECT id, comment_name as Name, comment_email, comment_text FROM comments WHERE page_id=?", thisPage.ID)
	if err != nil {
		log.Println(err)
	}
	for comments.Next() {
		var comment Comment
		comments.Scan(&comment.ID, &comment.Name, &comment.Email,
			&comment.CommentText)
		thisPage.Comments = append(thisPage.Comments, comment)
	}*/
	t, err := template.ParseFiles("templates/blog.html")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("PAGE: %#v\n", thisPage)
	t.Execute(w, thisPage)
}

// LoginPOST ...
func LoginPOST(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/login.html")
	if err != nil {
		log.Fatalln(err)
	}
	t.Execute(w, nil)
}

// APIPage ...
func APIPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pageGUID := vars["guid"]
	thisPage := Page{}
	fmt.Println(pageGUID)
	err := database.QueryRow("SELECT page_title,page_content,page_date FROM pages WHERE page_guid=?", pageGUID).Scan(&thisPage.Title,
		&thisPage.RawContent, &thisPage.Date)
	thisPage.Content = template.HTML(thisPage.RawContent)
	if err != nil {
		http.Error(w, http.StatusText(404), http.StatusNotFound)
		log.Println(err)
		return
	}
	APIOutput, err := json.Marshal(thisPage)
	fmt.Println(APIOutput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, thisPage)
}

// APICommentPost ...
func APICommentPost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("COMMENT POST!!!")
	fmt.Printf("REQ: %#v\n", r)
	var commentAdded bool
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	name := r.FormValue("name")
	email := r.FormValue("email")
	comments := r.FormValue("comments")
	formID := r.FormValue("id")
	fmt.Printf("INSERT INTO comments SET page_id ='%s', comment_name='%s', comment_email='%s', comment_text='%s'\n", formID, name, email, comments)
	res, err := database.Exec("INSERT INTO comments SET page_id=?, comment_name=?, comment_email=?, comment_text=?", formID, name, email, comments)
	if err != nil {
		log.Println(err)
	}
	ID, err := res.LastInsertId()
	if err != nil {
		log.Fatalln("--------------", err)
		commentAdded = false
	} else {
		commentAdded = true
	}
	fmt.Println("ID:::::::::::::::::::::::::: ", string(ID))
	commentAddedBool := strconv.FormatBool(commentAdded)
	var resp JSONResponse
	resp.Fields = make(map[string]string)
	resp.Fields["id"] = fmt.Sprintf("%d", ID)
	resp.Fields["added"] = commentAddedBool

	//jsonResp, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	//fmt.Fprintln(w, jsonResp)
	fmt.Fprintln(w, resp.Fields)
}

// APICommentPut ...
func APICommentPut(w http.ResponseWriter, r *http.Request) {
	fmt.Println("SET")
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Println(id)
	name := r.FormValue("name")
	email := r.FormValue("email")
	comments := r.FormValue("comments")
	res, err := database.Exec("UPDATE comments SET comment_name=?, comment_email=?, comment_text=? WHERE comment_id=?", name, email,
		comments, id)
	fmt.Println(res)
	if err != nil {
		log.Println(err)
	}
	var resp JSONResponse
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, jsonResp)
}

func weakPasswordHash(password string) []byte {
	hash := sha1.New()
	io.WriteString(hash, password)
	return hash.Sum(nil)
}

// RegisterPOST ...
func RegisterPOST(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	name := r.FormValue("user_name")
	email := r.FormValue("user_email")
	pass := r.FormValue("user_password")
	pageGUID := r.FormValue("referrer")
	// pass2 := r.FormValue("user_password2")
	gure := regexp.MustCompile("[^A-Za-z0-9]+")
	guid := gure.ReplaceAllString(name, "")
	password := weakPasswordHash(pass)
	res, err := database.Exec("INSERT INTO users SET user_name=?,user_guid=?, user_email=?, user_password=?", name, guid, email,
		password)
	fmt.Println(res)
	if err != nil {
		fmt.Fprintln(w, err)
	} else {
		http.Redirect(w, r, "/page/"+pageGUID, 301)
	}
}

func main() {
	// Database
	dbConn := fmt.Sprintf("%s:%s@/%s", DBUser, DBPass, DBDbase)
	fmt.Println(dbConn)
	db, err := sql.Open("mysql", dbConn)
	if err != nil {
		log.Println("Couldn't connect to" + DBDbase)
		log.Println(err)
	}
	database = db
	// Routes
	routes := mux.NewRouter()
	routes.HandleFunc("/api/pages", APIPage).
		Methods("GET").
		Schemes("https")
	routes.HandleFunc("/api/page/{guid:[0-9a-zA\\-]+}", APIPage).
		Methods("GET").
		Schemes("https")
	routes.HandleFunc("/api/comments", APICommentPost).
		Methods("POST")
	routes.HandleFunc("/api/comments/{id:[\\w\\d\\-]+}", APICommentPut).
		Methods("PUT")

	routes.HandleFunc("/register", RegisterPOST).
		Methods("POST").
		Schemes("https")
	routes.HandleFunc("/login", LoginPOST).
		Methods("POST").
		Schemes("https")

	routes.HandleFunc("/page/{guid:[0-9a-zA\\-]+}", ServePage)
	routes.HandleFunc("/home", ServeIndex)
	routes.HandleFunc("/", RedirIndex)

	http.Handle("/", routes)

	// Certs
	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}
	_ = cfg
	srv := &http.Server{
		Addr:    port,
		Handler: routes,
		//TLSConfig: cfg,
		//TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}

	log.Fatal(srv.ListenAndServe())
	//log.Fatal(srv.ListenAndServeTLS("certs/server.pem", "certs/server.key"))
}
