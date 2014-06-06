package main

import "fmt"
import "strings"

func main() {
     var t1, t2 string
     t1 = "Brett Ratner"
     t2 = "FFC, BR, RH"
     fmt.Println(strings.Split(t1,","))
     fmt.Println(strings.Split(t2,","))
}