package main

import(
	"fmt"
	"github.com/orchestrate-io/gorc"
)

func main() {
     apiKey := "9edc13f3-b67a-4d1d-bf96-a8159822d44f"
     c := gorc.NewClient(apiKey)                     // why no error return val?
     if err := c.Ping(); err != nil {
       fmt.Println("Something's wrong cap'n!")
     } else {
       fmt.Println("A-OK!")
     }
}