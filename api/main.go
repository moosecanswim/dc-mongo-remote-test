package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
)

type Post struct {
	Text      string    `json:"text" bson:"text"`
	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
}

var posts *mgo.Collection

func main() {
	// Connect to mongo
	viper.SetConfigName("settings")
	viper.AddConfigPath("$GOPATH/src/app/")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/go/src/app")

	viper.WatchConfig()
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln(fmt.Errorf("fatal error config file %s", err))
	}

	username := viper.GetString("mongo_user")
	password := viper.GetString("mongo_password")
	database := viper.GetString("mongo_db")
	address := viper.GetString("mongo_host")

	log.Println("my name is")
	log.Println(viper.GetString("name"))

	log.Println("Viper config address")
	log.Println(address)

	log.Println(fmt.Sprintf("mongo_user: %s mongo_password: %s, mongo_db: %s, mongo_host: %s", username, password, database, address))

	var connURL string
	if username == "" || password == "" {
		log.Println("No username or password")
		connURL = fmt.Sprintf("mongodb://%s/%s", address, database)
	} else {
		log.Println("Authenticating with username and password")
		connURL = fmt.Sprintf("mongodb://%s:%s@%s/%s", username, password, address, database)
	}
	log.Println(connURL)
	session, err := mgo.Dial(connURL)
	if err != nil {
		log.Fatalln(err)
		log.Fatalln("mongo err")
		os.Exit(1)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	// Get posts collection
	posts = session.DB("app").C("posts")

	// Set up routes
	r := mux.NewRouter()
	r.HandleFunc("/posts", createPost).
		Methods("POST")
	r.HandleFunc("/posts", readPosts).
		Methods("GET")

	http.ListenAndServe(":8080", cors.AllowAll().Handler(r))
	log.Println("Listening on port 8080...")
}

func createPost(w http.ResponseWriter, r *http.Request) {
	// Read body
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Read post
	post := &Post{}
	err = json.Unmarshal(data, post)
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}
	post.CreatedAt = time.Now().UTC()

	// Insert new post
	if err := posts.Insert(post); err != nil {
		responseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseJSON(w, post)
}

func readPosts(w http.ResponseWriter, r *http.Request) {
	result := []Post{}
	if err := posts.Find(nil).Sort("-created_at").All(&result); err != nil {
		responseError(w, err.Error(), http.StatusInternalServerError)
	} else {
		responseJSON(w, result)
	}
}

func responseError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func responseJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
