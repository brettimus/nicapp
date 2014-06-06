package main

import(
	"os"
	"fmt"
	"log"
	"strings"
	"net/http"
	"html/template"
	"github.com/orchestrate-io/gorc"

)

var apiKey = "9edc13f3-b67a-4d1d-bf96-a8159822d44f"

var templates = template.Must(template.ParseFiles("./templates/results.html"))

func main() {
     http.HandleFunc("/", home)
     http.HandleFunc("/results", results)
     http.HandleFunc("/related", related)
     if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
     	log.Fatal("Server did not start. Error:", err)
     }
}

func home(w http.ResponseWriter, r *http.Request) {
     err := r.ParseForm() //
     if err != nil {
       log.Fatal("Oops. Parsing form failed.", err)
     }
     http.ServeFile(w, r, r.URL.Path[1:])
}

func related(w http.ResponseWriter, r *http.Request) {
     // TODO
}

func results(w http.ResponseWriter, r *http.Request) {
     err := r.ParseForm()
     if err != nil {
     	log.Fatal("Oops. Parsing form failed.", err)
     }
     value := r.URL.Query() // type: map[string][]string

     qry := construct_query(value)
     fmt.Println(qry)

     c := gorc.NewClient(apiKey)
     if err = c.Ping(); err != nil {
     	// Some form of error page?
	// for now just print line
	fmt.Println("Ping failed! Error:", err)
     }

     res, err := c.Search("moobies", qry, 100, 0)
     if err != nil {
       fmt.Println("> Search threw an error! Here:", err)
       return
     } 
     fmt.Println(res)

     qValues := make([]map[string]interface{}, len(res.Results))
     for i, result := range res.Results {
     	 result.Value(&qValues[i])
	 qValues[i]["Title"] = result.Path.Key
	 fmt.Println(qValues[i])
     }
     // Don't have to deal with "HasNext" because we have fewer than 100 movies in the database

     err = templates.ExecuteTemplate(w, "results.html", qValues)
     if err != nil {
     	http.Error(w, err.Error(), http.StatusInternalServerError)
	return
     }
}

/* I am embarrassed to have written this function */
/* There must be a better way */
func construct_query(value map[string][]string) string {
     oQ := ""
     oQ += "("

     gRes, gEmp := q_helper("Genre:", value["gen"])
     rRes, rEmp := mpaa_helper("Rated:", value["rating"])
     iRes, iEmp := imdb_helper(value["imdb"])
     fmt.Println(iRes,iEmp)

     if q := value["q"]; len(q) != 0 && q[0] != "" {
     	oQ += "\"" + q[0] + "\"~4"
	if !gEmp || !rEmp || !iEmp {
	   oQ += " AND "
	}
     } 

     if gEmp && rEmp && iEmp {
     	if q:= value["q"]; len(q) != 0 && q[0] != "" {
	   return oQ + ")"
	} else {
	   // Return all movies
	   return "*"
	}
     } 
     oQ  += "("

     switch {
     	    case !gEmp && !rEmp && !iEmp : oQ += gRes + " AND " + rRes + " AND " + iRes
     	    case !gEmp && !rEmp          : oQ += gRes + " AND " + rRes
	    case !gEmp && !iEmp          : oQ += gRes + " AND " + iRes
	    case !rEmp && !iEmp          : oQ += rRes + " AND " + iRes
	    case !gEmp                   : oQ += gRes
	    case !rEmp                   : oQ += rRes
	    case !iEmp                   : oQ += iRes
     }
     return oQ + "))"
}

/*** CREATE MAKE_HELPER CLOSURE ***/

// Returns the string to concatenate to the query, and a boolean about whether or not it's empty
func q_helper(pfix string, input []string) (string, bool) {
     var res string;
     if len(input) == 0 || input[0] == "" {
     	return res, true
     }

     const qry string = "value."
     res += qry + pfix 
     for i, v := range input {
     	 if i == 0 {
	    res += "("
	 }
	 res += "\"" + v + "\""

	 if i != len(input) -1 {
	    res += " OR "
	 } else {
	    res += ")"
	 }	 
     }
     return res, false
}

func mpaa_helper(pfix string, input []string) (string, bool) {
     var res string;
     if len(input) == 0 || input[0] == "" {
     	return res, true
     }

     only_pg := is_in("PG", input) && !is_in("PG-13",input)
     // THEN WHAT?!

     const qry string = "value."
     res += qry + pfix 
     for i, v := range input {
     	 if i == 0 {
	    res += "("
	 }

	 res += "\"" + v + "\""

	 if i != len(input) -1 {
	    res += " OR "
	 } else {
	    res += ")"
	 }	 
     }

     if only_pg {
       res += " AND -value.Rated:\"PG-13\""
     }
     return res, false
}

func is_in(s string, a []string) bool {
     for _, thing := range a {
         if s == thing {
            return true
         }
     }
     return false
}

/* queries a range */
func imdb_helper(input []string) (string, bool) {
     var res string;
     if len(input) == 0 || input[0] == "" {
     	return res, true
     }
     res += "value.ImdbRating:"
     for i, v := range input {
     	 if i == 0 {
	    res += "("
	 }

	 if !is_valid_imdb(v) {
	    res += ""
	 } else {
 	    res += "[" + strings.Replace(v,"-"," TO ",1) + "]"
	 }

	 if i != len(input) -1 {
	    res += " OR "
	 } else {
	    res += ")"
	 }
     }
     return res, false
}

func is_valid_imdb(v string) (res bool) {
     switch {
     	    case v == "0-5"  : res = true
	    case v == "5-6"  : res = true
	    case v == "6-7"  : res = true
	    case v == "7-8"  : res = true
	    case v == "8-10" : res = true
	    default          : return res       // res is initialized as false
     }
     return res
}

type FilmData struct {
     Year        int
     Director    []string
     Genre       []string
     Poster      string
     Rated       string
     ImdbRating  float64
     ImdbID      string
}

/* example from 10 things you didn't know about Go */
/* http://talks.golang.org/2012/10things.slide#4 */
/*
type Item struct {
  Title string
  URL   string
}
type Response struct {
  Data struct {
    Children []struct {
      Data Item
    }
  }
}
*/

