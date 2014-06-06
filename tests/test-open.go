/* FAIL! 
   Check out:
   https://groups.google.com/forum/#!topic/golang-nuts/UZ8W5IH95wg
 */
package main

import(
	"fmt"
	"os"
//	"io"
	"bufio"
)

func main() {
     fmt.Println("hi")
     file := os.Stdout
     var err error
     file, err = os.Create("/Users/bbeutell/go-prax/nichApp/test-file.txt")
     if err != nil {
       fmt.Println("Heyyyyy bad stuff")
     }
     writer := bufio.NewWriter(file)
     _, err = writer.WriteString("dooooood")
     writer.Flush()
     file.Close()
     fmt.Println("Final error:",err)
}