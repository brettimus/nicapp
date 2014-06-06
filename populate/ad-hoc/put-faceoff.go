package main

import(
	"fmt"
	"os"
	"io/ioutil"
	"net/http"
	"github.com/orchestrate-io/gorc"
	"encoding/json"
	"log"
	"strings"
	"strconv"
)

func main() {
     var err error
     fo := "Face-Off"
     apiKey := os.Getenv("GORCKEY")
     c := gorc.NewClient(apiKey)
     if err = c.Ping(); err != nil {
       fmt.Println("Ping failed!!! NOOOOOO")
     }

     var data []byte
     data, err = get_data(get_URL())
     if err != nil {
       fmt.Println("Something went wrong getting the data:", err)
       log.Fatal(err)
     }
     var fd FilmData
     err = json.Unmarshal(data, &fd)
     if err != nil {
       fmt.Println("Error unmarshalling")
       log.Fatal(err)
     }
     var fdf FilmDataForc
     fdf, err = FD_to_FDF(fd)
     if err != nil {
       fmt.Println("Conversion to FilmDataForc failed with error:", err)
       log.Fatal(err)
     }
     _, err = c.PutIfAbsent("moobies", fo, fdf)
     if err != nil {
        fmt.Println("PutIfAbsent failed with error:", err)
	log.Fatal(err)
     }
}

func get_data(URL string) ([]byte, error) {
     req, err := http.NewRequest("GET",URL,nil)
     if err != nil {
       fmt.Println("Making request failed!")
       return nil, err
     }
     client := &http.Client{}
     resp, err := client.Do(req)
     if err != nil {
       fmt.Println("'Do'ing request failed!")
       return nil, err
     }
     defer resp.Body.Close()

     body, err := ioutil.ReadAll(resp.Body) // Reads into []byte
     if err != nil {
       return nil, err
     }
     return body, nil
}

func get_URL() string {
     return "http://omdbapi.com/?t=Face%2FOff"
}


/* struct into which we unmarshal JSON */
type FilmData struct {
     Year string
     Director string
     Genre string
     Poster string
     Rated string
     ImdbRating string
     ImdbID string
}

type FilmDataForc struct {
     Year int
     Director []string
     Genre []string
     Poster string
     Rated string
     ImdbRating float64
     ImdbID string
}



func FD_to_FDF(fd FilmData) (fdf FilmDataForc, err error) {
     var int_year int64
     int_year, err = strconv.ParseInt(fd.Year, 10, 64)
     if err != nil {
       fmt.Println("Error parsing int for fd.Year")
     }
     fdf.Year = int(int_year)
     fdf.Director = strings.Split(fd.Director,",")
     fdf.Genre = strings.Split(fd.Genre,",")
     fdf.Poster = fd.Poster
     fdf.Rated = fd.Rated
     fdf.ImdbRating, err = strconv.ParseFloat(fd.ImdbRating,64)
     if err != nil {
        fmt.Println("Error parsing float64 for ImdbRating")
     }
     fdf.ImdbID = fd.ImdbID

     return fdf, err
}
