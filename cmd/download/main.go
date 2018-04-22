package main

import (
	"fmt"
	"os"

	"github.com/brandoneprice31/freedrive/setup"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println(`2 args are required:
 - path/to/freedrive/dir
 - path/to/download/folder
 		`)
		os.Exit(2)
	}

	fd := os.Args[1]
	m, err := setup.Manager(fd)
	if err != nil {
		panic(err)
	}

	dl := os.Args[2]
	m.Download(dl)
}
