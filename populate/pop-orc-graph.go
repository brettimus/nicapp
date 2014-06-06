package main

import(
	"fmt"
	"os"
	"io"
	"log"
	"bufio"
	"encoding/json"
	"github.com/orchestrate-io/gorc"
	"strings"
	"strconv"
	"math"
)

func main() {
     films, err := read_films("nfs.txt")
     if err != nil {
     	log.Fatal(err)
     }
     fmt.Println(">>> Read Success! \n\n", films, "\n")
     
     var filmData []FilmGraphData
     filmData, err = get_film_data(films)
     if err != nil {
       fmt.Println(">>> get_film_data failed! Error:", err)
       log.Fatal(err)
     }
     fmt.Println(">>> Here's what the data look like: \n", filmData, "\n\n")
     
     apiKey := os.Getenv("GORCKEY")
     c := gorc.NewClient(apiKey)
     if err = c.Ping(); err != nil {
     	fmt.Println("Ping failed! Error: ", err)
	log.Fatal(err)
     } 
     fmt.Println("Ping success! \n\n")

     gRels, _ := find_edges(filmData)
     /* fmt.Println(gRels) */
     for _, gRel := range gRels {
     	 err = put_edge(c,gRel.From, gRel.To, gRel.Type)
	 if err != nil {
	   fmt.Println(">>> Error putting relation", gRel.Type, "from", gRel.From, "to", gRel.To, ".")
	 } else {
	   fmt.Println(">>> Success!")
	 }
     	 err = put_edge(c,gRel.To, gRel.From, gRel.Type)
	 if err != nil {
	   fmt.Println(">>> Error putting relation", gRel.Type, "from", gRel.To, "to", gRel.From, ".")
	 } else {
	   fmt.Println(">>> Success!")
	 }
     }
}

func put_edge(c *gorc.Client, k1, k2, kind string) (err error) {
     err = c.PutRelation("moobies", k1, kind, "moobies", k2)
     return err
}

func find_edges(filmData []FilmGraphData) (conns []GraphCon, err error) {
     var match bool
     var con GraphCon
     for i := 0; i < len(filmData) - 1; i++ {
       for j := i + 1; j < len(filmData); j++ {
       	   f1 := filmData[i]
	   f2 := filmData[j]

	   /* TODO
	     
	     [] Error handling
	     [X] Create a connection type
	     [X] Resolve shared genre connections
	     [] Pass pointer to res to each function
	     ([] Fix logic by adding "is comedy" to JSON?)
	   
	   */

	   match, con = find_s_gen(f1,f2)
	   if match {
	     conns = append(conns, con)
	   }
	   match, con = is_s_imdb(f1,f2)
	   if match {
	      conns = append(conns,con)
	   }
       }
     }
     return conns, err
}


// func write_parsed_data() {}

func get_film_data(films []string) (res []FilmGraphData, err error) {
     for _, f := range films {
     	 fmt.Println(f)
	 var filmData FilmDataForc
	 filmData, err = fetch_local_json(f)
	 if err != nil {
	   fmt.Println("Call to fetch_local_json(",f,") failed.")
	   return nil, err
	 }

	 var temp FilmGraphData
	 temp.title = f
	 temp.director = filmData.Director
	 temp.genre = filmData.Genre
	 temp.mpaa = filmData.Rated
	 temp.imdb = filmData.ImdbRating

	 res = append(res, temp)
     }
     return res, err
}

func fetch_local_json(film string) (data FilmDataForc, err error) {
     filename := film + ".json"
     file := os.Stdin
     file, err = os.Open(filename)
     defer file.Close()
     if err != nil {
        fmt.Println("Opening the file failed!")
     	return data, err
     }
     reader := bufio.NewReader(file)
     
     var da_json []byte
     da_json, _, err = reader.ReadLine()
     err = json.Unmarshal(da_json, &data)
     if err != nil {
       fmt.Println("Unmarshalling", filename, "failed:", err)
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

func find_s_dir(f1, f2 FilmGraphData) (match bool, res GraphCon) {
     for _, dir := range f1.director {
     	 if is_in(dir, f2.director) {
	    match = true
	    res.From = f1.title
	    res.To = f2.title
	    res.Type = "s_Director"
	 }
     }
     return match, res
}

func find_s_gen(f1, f2 FilmGraphData) (match bool, res GraphCon) {
     var count int64
     count = 0
     for _, gen := range f1.genre {
     	 if is_in(gen, f2.genre) {
	    match = true
	    count += 1
	 }
     }
     res.From = f1.title
     res.To = f2.title
     fmt.Println(count)
     tipo := "s_Genre_" + strconv.FormatInt(count,10)
     res.Type = tipo
     return match, res
}
func is_s_mpaa(f1, f2 FilmGraphData) (match bool, res GraphCon) {
     if f1.mpaa == f2.mpaa {
       match = true
       res.From = f1.title
       res.To = f2.title
       res.Type = "s_Rated"
     } 
     return match, res
}
func is_s_imdb(f1, f2 FilmGraphData) (match bool, res GraphCon) {
     if f1.imdb == 0 || f2.imdb == 0 {
     	return match, res
     }

     if rDiff := math.Abs(f1.imdb - f2.imdb); rDiff <= .3 {
     	match = true
	res.From = f1.title	
	res.To = f2.title
	switch match {
	  case rDiff <= .1: res.Type = "s_Imdb_strong"
	  case rDiff <= .2: res.Type = "s_Imdb_medium"
	  default: res.Type = "s_Imdb_weak"
	}
     }
     return match, res
}

func is_in(s string, a []string) bool {
     for _, thing := range a {
     	 if s == thing {
	    return true
	 }
     }
     return false
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

type FilmGraphData struct {
     title string
     director []string
     genre []string
     mpaa string
     imdb float64
}

type GraphCon struct {
     From string
     To string
     Type string
}
