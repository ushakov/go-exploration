package main

import "flag"
import "http"
import "fmt"

import "tiles"

var httpPort = flag.String("http", ":8081", "HTTP port to listen on")
var tilesPath = flag.String("tiles", "", "Path to tiles directory")

func test(c *http.Conn, req *http.Request) {
	fmt.Fprintf(c, "Hello, World at %s!\n", req.URL.Path);
}

func main() {
	flag.Parse();
	ts, err := tiles.Load(*tilesPath);
	if err != nil {
		panic(err.String());
	}

	http.Handle("/gettile", ts);
	http.Handle("/", http.FileServer("./static", ""));
	fmt.Println("Serving on ", *httpPort, "...");
	err = http.ListenAndServe(*httpPort, nil);
	if err != nil {
		panic("ListenAndServe: ", err.String());
	}
}
