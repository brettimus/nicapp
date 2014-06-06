package main

import "fmt"
//import "strings"
//import "os"
//import "io"
import "io/ioutil"
//import "bufio"
import "net/url"
import "net/http"
import "encoding/json"
import "log"

func main() {
     //URL := "http://omdbapi.com/?t=The+Family+Man&y=2000"
     URL := make_url("The Family Man","2000")
     fmt.Println(URL)
     data, err := get_data(URL)
     if err != nil {
       fmt.Println("error in first")
       log.Fatal(err)     	
     }
     var fd FilmData
     err = json.Unmarshal(data, &fd)
     if err != nil {
       fmt.Println("error in second")
       log.Fatal(err)
     }
     fmt.Println("Hi")
     fmt.Println(fd)
}

/* Calls the omdbapi */
func get_data(URL string) ([]byte, error) {
  req, err := http.NewRequest("GET", URL, nil)
  if err != nil {
    return nil, err
  }

  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
    return nil, err
  }
  defer resp.Body.Close()

  body, err := ioutil.ReadAll(resp.Body)   // Reads into a byte array
  if err != nil {
    return nil, err
  }
  return body, nil
}

/* Constructs URL for omdbapi query */
func make_url(title, year string) (URL string) {
  base_url := "http://omdbapi.com/"
  URL = base_url + "?t=" + url.QueryEscape(title) + "&y=" + year

  return URL
}

/* struct in which we store unmarshalled JSON */
type FilmData struct {
     Title string 
     Year string
     Director string
     Genre string
     Poster string
     Rated string
     ImdbRating string
     ImdbID string
}