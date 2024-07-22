package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"ascii-art-web-stylize/utils"
)

func main() {
	if len(os.Args) != 1 {
		log.Println("Usage: <go run .> <go run main.go>")
		return
	}
	http.HandleFunc("/", utils.ServeIndex)
	http.HandleFunc("/ascii-art", utils.GenerateASCIIArt)

	port := ":8000"
	portNum, err := strconv.Atoi(port[1:])
	if err != nil {
		fmt.Printf("Error: Unable to covert %v to integer\n", port[1:])
		return
	}
	if portNum < 1024 || portNum > 65535 {
		fmt.Println("Error: The port you are using is either reserved or doesn't exist")
		return
	}
	finalPort := ":" + strconv.Itoa(portNum)
	log.Printf("Server running at http://localhost%v\n", finalPort)
	http.ListenAndServe(finalPort, nil)
}
