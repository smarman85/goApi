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

type Person struct {
  Name string
  Age int
}

type Posts map[string]interface{}

type Post struct {
  Name string
  Desc string
}

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

func respondWithError(w http.ResponseWriter, code int, message string) {
    respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    response, _ := json.Marshal(payload)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(response)
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

func getPost(w http.ResponseWriter, r *http.Request) {

  vars := mux.Vars(r)
  id := vars["id"]
  f.Println(id)

  DB_CON := psqlCon()
  // create db pool
  db, err := sql.Open("postgres", DB_CON)
  if err != nil {
    log.Fatal("Failed to open DB connection: ", err)
  }
  
  row := db.QueryRow("SELECT * FROM posts WHERE id=$1", id)
  json.NewEncoder(w).Encode(row)
  f.Println(row)

}

func createPost(w http.ResponseWriter, r *http.Request) {
  //var p NewPost
  var p Posts

  err := json.NewDecoder(r.Body).Decode(&p)
  if err != nil {
      http.Error(w, err.Error(), http.StatusBadRequest)
      return
  }
  DB_CON := psqlCon()
  // create db pool
  db, err := sql.Open("postgres", DB_CON)
  if err != nil {
    log.Fatal("Failed to open DB connection: ", err)
  }

  input, _ := json.Marshal(p)
  // Do something with the Person struct...
  _, err = db.Exec("INSERT INTO posts (post_info) VALUES ($1) RETURNING ID;", input)
  if err != nil {
    log.Fatal(err)
  }
  if err != nil {
    log.Fatal("Failed to update db: ", err)
  }

  f.Fprintf(w, "%+v", string(input))

}

func deletePost(w http.ResponseWriter, r *http.Request) {

  vars := mux.Vars(r)
  id := vars["id"]
  DB_CON := psqlCon()
  // create db pool
  db, err := sql.Open("postgres", DB_CON)
  if err != nil {
    log.Fatal("Failed to open DB connection: ", err)
  }

  _, err = db.Exec("DELETE FROM posts WHERE id=$1", id)
  if err != nil {
    log.Fatal("Failed to remove from database: ", err)
  }

}

func main() {
  dbConn := psqlCon()
  f.Println(dbConn)
  router := mux.NewRouter()
  router.HandleFunc("/", home)
  router.HandleFunc("/posts", getPosts).Methods("GET")
  router.HandleFunc("/post/{id}", getPost).Methods("GET")
  router.HandleFunc("/new_post", createPost).Methods("POST")
  router.HandleFunc("/retract/{id}", deletePost).Methods("DELETE")
  log.Fatal(http.ListenAndServe(":8090", router))
}
