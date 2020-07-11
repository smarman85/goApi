package main

import (
  f "fmt"
  "log"
  "database/sql"
  _ "github.com/lib/pq"
  "github.com/gorilla/mux"
  "github.com/joho/godotenv"
  "encoding/json"
  "net/http"
  "os"
)

type postInfo struct {
  ID int
  INFO json.RawMessage
  //INFO []uint8
  //INFO map[string]string
}

type posts struct {
  Posts []postInfo
}

func psqlCon() string {
  err := godotenv.Load("psql.env")
  if err != nil {
    log.Fatal("Error Loading psql.env file")
  }
  user := os.Getenv("POSTGRES_USER")
  pass := os.Getenv("POSTGRES_PASSWORD")
  host := os.Getenv("POSTGRES_HOST")
  db   := os.Getenv("POSTGRES_DB")

  con := f.Sprintf(
    "postgres://%v:%v@%v:5432/%v?sslmode=disable",
    user,
    pass,
    host,
    db,
  )
  return con
}

func home(w http.ResponseWriter, r *http.Request) {
  msg := map[string]string {
    "message": "Welcome. Please have a look around",
  }
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(msg)
}

func queryPosts(postData *posts) error {
  DB_CON := psqlCon()
  // create db pool
  db, err := sql.Open("postgres", DB_CON)
  if err != nil {
    log.Fatal("Failed to open DB connection: ", err)
  }

  //rows, err := db.Query(`SELECT * FROM new_posts;`)
  rows, err := db.Query(`SELECT * FROM posts;`)
  //rows, err := db.Query(`SELECT id, post_info -> 'name' FROM posts;`)
  if err != nil {
    return err
  }
  defer db.Close()

  for rows.Next() {
    post := postInfo{}
    err := rows.Scan(
      &post.ID,
      &post.INFO,
    )
    if err != nil {
      return err
    }
    postData.Posts = append(postData.Posts, post)
  }
  err = rows.Err()
  if err != nil {
    return err
  }
  return nil

}

func getPosts(w http.ResponseWriter, r *http.Request) {
  postData := posts{}
  err := queryPosts(&postData)
  if err != nil {
    http.Error(w, err.Error(), 500)
    return
  }

  out, err := json.Marshal(postData)
  if err != nil {
    http.Error(w, err.Error(), 500)
    return
  }

  f.Fprint(w, string(out))
}

func main() {
  dbConn := psqlCon()
  f.Println(dbConn)
  router := mux.NewRouter()
  router.HandleFunc("/", home)
  router.HandleFunc("/posts", getPosts).Methods("GET")
  log.Fatal(http.ListenAndServe(":8090", router))
}
