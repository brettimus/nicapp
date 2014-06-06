/* Tests GET and write of Face-Off */
package main

import(
	"fmt"
	"os"
	"io"
//	"log"
	"bufio"
	"encoding/json"
	"github.com/orchestrate-io/gorc"
	"strings"
)

func main() {
     film := "Face-Off"
     apiKey := os.Getenv("GORCKEY")
     c := gorc.NewClient(apiKey)
     if err := c.Ping(); err != nil {
     	fmt.Println("Ping success!\n")
     }
     
     filmData, err := get_film_data(c, film)
     if err == nil {
     	fmt.Println("> GET of", film, " successful.")
     }
     err = write_film_data(film, filmData)
     if err == nil {
     	fmt.Println("> Write of", film, "data successful.")
     } else {
       fmt.Println("> WRITE FAILED! error:", err)
     }
}

func write_film_data(title string, data map[string]interface{}) (err error) {
     oFilename := title + ".json"
     oFile := os.Stdout
     oFile, err = os.Create(oFilename)
     defer oFile.Close()
     if err != nil {
       return err
     }
     writer := bufio.NewWriter(oFile)
     defer writer.Flush()

     var toWrite []byte
     toWrite, err = json.Marshal(data)
     if err != nil {
       return err
     }
     _, err = writer.Write(toWrite)
     if err != nil {
     	fmt.Println("WARNING! Writing data of ", oFilename, "failed.")
     }
     return err
}

/* This is wasteful of API calls (one GET per film) */
/* Consider running a search instead */

func get_film_data(c *gorc.Client, film string) (data map[string]interface{}, err error) {

     var result *gorc.KVResult
     result, err = c.Get("moobies",film)
     if err != nil {
     	fmt.Println("Getting", film, "caused an error!")
	return nil, err
     }
     data = make(map[string]interface{})
     err = result.Value(&data)
     if err != nil {
       fmt.Println("Marshalling data for", film, "failed!")
       return nil, err
     }
     return data, err
}

func read_films(iFilename string) (films []string, err error) {
     iFile := os.Stdin
     iFile, err = os.Open(iFilename)
     if err != nil {
     	return nil, err
     }
     reader := bufio.NewReader(iFile)
     var line string
     for {
          line, err = reader.ReadString('\n')
	  if err == io.EOF {
	     err = nil
	     return films, err
	  } else if err != nil {
	    return nil, err
	  }
     	  title := parse_title(line)
	  films = append(films, title)
     }
}

func parse_title(line string) (title string) {
     title = strings.Split(line,",")[0]
     return title
}
