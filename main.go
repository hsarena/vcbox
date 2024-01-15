/*
Copyright © 2023 NAME HERE hamid@hsarena.me
*/
package main

import (
	"log"
	"os"

	"github.com/hsarena/vcbox/cmd"
)

func main() {
	logFile, err := os.OpenFile("./vcbox.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)

	log.Println("VCBOX Started")
	cmd.Execute()
}
