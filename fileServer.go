package main

import (
	"fmt"
	//"fmt"
	"log"
	"net/http"
	"flag"
	"os"
)

func main()  {
	var root = flag.String("r", "./", "The files root directory.")
	var port = flag.String("p", ":8080", "The port file server listened.")
	flag.Parse()

	if (*port)[0:1] != ":" {
		*port = ":" + *port
	}

	rInfo, err := os.Stat(*root)
	if err == nil && !rInfo.IsDir() {
		log.Fatalf("%s is not directory.", *root)
	}
	if err != nil && os.IsNotExist(err) {
		log.Fatalf("%s is not exist.", *root)
	}

	fmt.Printf("Root:%s, Port:%s", *root, *port)
	http.Handle("/", http.FileServer(http.Dir(*root)))
	log.Fatal(http.ListenAndServe(*port, nil))
}
