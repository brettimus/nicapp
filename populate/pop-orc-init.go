/* Gets the data, puts it in orchestrate */
/* Failed for "Face/Off", probably due to filepath */
package main

import( 
	"fmt"
	"strings"
	"strconv"
	"os"
	"io"
 	"io/ioutil"
	"bufio"
	"net/url"
	"net/http"
	"encoding/json"
	"github.com/orchestrate-io/gorc"
	"log"
)

func main() {
     fmt.Println("Begin...")
     films, err := read_films("nfs.txt")     
     if err != nil {
     	panic(err)
     }  
/* TEST    
     test, err := FD_to_FDF(films["The Family Man"])
     fmt.Println(test)
*/

     filmz := make(map[string]FilmDataForc)
     for k, v := range films {
     	 filmz[k], _ = FD_to_FDF(v)
	 fmt.Println(k)
	 fmt.Println(filmz[k])
     }
     apiKey := "9edc13f3-b67a-4d1d-bf96-a8159822d44f"
     err = pop_orc(apiKey, filmz)
     if err != nil {
     	log.Fatal(err)
     }
}

func pop_orc(apiKey string, films map[string]FilmDataForc) (err error) {
     c := gorc.NewClient(apiKey)
     if err = c.Ping(); err != nil {
       log.Fatal(err)
     } else {
       fmt.Println("Ping Success!")
     }

     for k, v := range films {
       path, err := c.PutIfAbsent("moobies",k,v)
       if err != nil {
         fmt.Println(err)
       } else {
       	 fmt.Println(path)
       }
       
     }
/* TEST
     path, err := c.PutIfAbsent("moobies","The Family Man",films["The Family Man"])
     fmt.Println(path)
*/
     return err
}

func read_films(iFilename string) (fds map[string]FilmData, err error) {
     /* Prep file read */
     iFile := os.Stdin
     if iFile, err = os.Open(iFilename); err != nil {
     	log.Fatal(err)       
     }
     defer iFile.Close()
     reader := bufio.NewReader(iFile)
     /* Gotta make that map */
     fds = make(map[string]FilmData)

     /*
     	(1) Loop through file, 
        (2) grab URL [make_url()], 
	(3) call api to get data [get_data()], 
	(4) unmarshal data
      */

     // (1) //
     for {
     	 var line, title, year string
	 line, err = reader.ReadString('\n')
	 if err == io.EOF {
	    err = nil        // End of file isn't really an error for us
	    return fds, err
	 } else if err != nil {
	   return nil, err
	 }

     // (2) //
	 title = strings.Split(line, ",")[0]
	 year = strings.Split(line, ",")[1][:4]
         URL := make_url(title, year)
	 fmt.Println("Following the URL:", URL)
     // (3) //
         var data []byte
         data, err = get_data(URL)
	 if err != nil {
	   return nil, err
	 }
     // (4) //
     	var fd FilmData
	err = json.Unmarshal(data, &fd)
	if err != nil {
	  fmt.Println("error Unmarshalling")
	  log.Fatal(err)
	}
	fds[title] = fd
     }
}

/* Calls the omdbapi */
func get_data(URL string) ([]byte, error) {
  req, err := http.NewRequest("GET", URL, nil)
  if err != nil {
    fmt.Println("Error creating request.")
    return nil, err
  }

  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
    fmt.Println("Error 'Do'ing request")
    fmt.Println(err)
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