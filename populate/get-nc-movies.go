/* Creates file 'nfs.txt' containing newline-separated list of Nic Cage films */
package main

import(
	"fmt"
	"github.com/hailiang/html-query"
 	. "github.com/hailiang/html-query/expr"
 	"io"
 	"bufio"
 	"net/http"
 	"os"
)

func main() {
  writeFilms(getFilms(), "nfs.txt")
}

func getFilms() (ncFilms []string) {
  r := get("http://www.imdb.com/name/nm0000115/")
  defer r.Close()
  root, err := query.Parse(r)
  checkError(err)
  // Janky Hack 1 - Find() is a BFS, so the second call grabs the div we need 
  root.Find(Div, Id("filmography")).Find(Div, Class("filmo-category-section")).Descendants(Class("filmo-row")).For(func(item *query.Node) {
    if ncFilm := item.B().Ahref().Text(); ncFilm != nil {
      if ncFY := item.Span(Class("year_column")).Text(); ncFY != nil {
      	 // Janky Hack 2 there are issues with having a slash in the title. We rid ourselves of it here!
      	 if *ncFilm == "Face/Off" {
	    *ncFilm = "Face-Off"
	 }
      	 ncFilms = append(ncFilms, *ncFilm + "," + (*ncFY)[3:7])      // split is to unescape year txt
      }
    }
  })
  return ncFilms
}

func writeFilms(films []string, oFilename string) {
  oFile := os.Stdout
  var err error
  oFile, err = os.Create(oFilename)
  defer oFile.Close()                // Will close with panic if checkError fails
  checkError(err)     
  writer := bufio.NewWriter(oFile)
  defer writer.Flush()              // This could fail -- consider more defensive appraoch (p. 36)
  /* iterate through films, print  */
  for _, film := range films {
    _, err := writer.WriteString(film)
    checkError(err)
    fmt.Fprint(writer,"\n")            // Is there a way to write newLines without this?
    fmt.Println(film, "was written.")
  }
}

func get(url string) io.ReadCloser {
  resp, err := http.Get(url)
  checkError(err)
  return resp.Body
}

func checkError(err error) {
     if err != nil {
       panic(err)
     }
}