package main

import(
	"fmt"
	"github.com/orchestrate-io/gorc"
	"net/url"
//	"encoding/json"
	"os"
)

func main() {
  var err error

  apiKey := os.Getenv("GORCKEY")
  collection := "moobies"

  k := "Face/Off"   // The offending key
  var v TestData
  v.Year = 1997
  v.Rating = "R"
  
  c := gorc.NewClient(apiKey)
  if err = c.Ping(); err != nil {
     panic(err)
  }
  fmt.Println("Ping Success! We're in!")
  
  fmt.Println("PUTting unescaped 'Face/Off'...")
  _, err = c.Put(collection, k, v)
  if err != nil {
    fmt.Println(k,"failed with an error:", err)
  }

  fmt.Println("PUTting escaped 'Face/Off'...")
  k = url.QueryEscape(k)
  _, err = c.Put(collection, k, v)
  if err != nil {
     fmt.Println(k, "failed with an error:", err)
  }
}

type TestData struct {
     Year int
     Rating string
}

