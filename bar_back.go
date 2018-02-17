package main

import (
   "fmt"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/postgres"
    "github.com/gorilla/mux"
    "net/http"
    "log"
    "encoding/json"
     "github.com/matryer/respond"
     "errors"
)


func main() {
  db, err := gorm.Open("postgres", "host=localhost user=bar_back dbname=bar_back sslmode=disable password=back_bar")
  defer db.Close()
  InitDBModel(db)
  InitRouter()
  if err != nil {
      panic(err)
   }
  log.Fatal(http.ListenAndServe(":8080", nil))
}


func InitRouter() {
  r := mux.NewRouter()
  r.HandleFunc("/spirits", SpiritsHandler)
  http.Handle("/", r)
}

type ok interface {
  OK() error
}

func (s *Spirit) OK() error {
  if len(s.Name) == 0 {
    return errors.New("Required:name")
  }
  if len(s.Family) == 0 {
    return errors.New("Required:family")
  }
  return nil
}

func decode(r *http.Request, v interface{}) error {
    if err := json.NewDecoder(r.Body).Decode(v); err != nil {
        return err
    }
    if validatable, ok := v.(ok); ok {
        return validatable.OK()
    }
    return nil
}

func SpiritsHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method == "GET" {
    fmt.Fprintf(w, "List of spirits")
  }else if r.Method == "POST" {
    var s Spirit
    if err := decode(r, &s); err != nil {
      respond.With(w, r, http.StatusBadRequest, err)
      return
    }
  respond.With(w, r, http.StatusOK, &s)
  }else {
    respond.With(w, r, http.StatusNotFound, "Not Found")
  }
}

type Spirit struct {
   Id     int64
   Name   string
   Family string
   Manufacturer string
}


func (s Spirit) String() string {
   return fmt.Sprintf("Spirit<%d %s %s %s>", s.Id, s.Name, s.Family, s.Manufacturer)
}

func InitDBModel(db *gorm.DB) {
   // Migrate the schema
  db.AutoMigrate(&Spirit{})

//   spirit1 := &Spirit{
//      Name:   "Campari",
//      Family: "aperitif",
//      Manufacturer: "Campari Group",
//   }
  
//  db.Create(spirit1)

//  var spirit Spirit
//  db.First(&spirit, 1)
//  fmt.Println(spirit)
    fmt.Println("Db migrated")
}

