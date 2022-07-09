package swr

import (
	"flag"
	"log"
)

var editMode = flag.Bool("editmode", false, "Used to run the server in editor mode for offline world building.")

func Init() {
	log.Printf("Init\n")
}

func Main() {

	flag.Parse()

	log.Printf("Starting Server version %s\n", version)

	DB().Load()
	if *editMode {
		Editor()
	} else {
		ServerStart(Config().Addr)
	}

	DB().Save()
}

func Editor() {

}
