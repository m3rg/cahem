package main

import (
	"fmt"
	"log"
)

var config Config

func main() {
	config = *readConfig()

	newScheduler(func() {
		log.Println("Data is being fetched.")
		err := Fetch()
		if err != nil {
			fmt.Println(err)
		}
	})
}
