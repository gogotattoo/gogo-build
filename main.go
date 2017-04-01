package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	gopath := os.Getenv("GOPATH")
	gogotattooPrefix := "/src/github.com/gogotattoo/"
	artists := []string{"gogo", "aid", "xizi"}
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)
	runHugo(gopath+gogotattooPrefix+"/gogo.tattoo", dir+"/gogo.tattoo/", "http://gogo.tattoo/")
	for _, artist := range artists {
		runHugo(gopath+gogotattooPrefix+artist, dir+"/gogo.tattoo/"+artist, "http://gogo.tattoo/"+artist)
	}
	runHugo(gopath+gogotattooPrefix+"/gogotattoo.github.io", dir+"/gogotattoo.github.io/", "https://gogotattoo.github.io/")
	for _, artist := range artists {
		runHugo(gopath+gogotattooPrefix+artist, dir+"/gogotattoo.github.io/"+artist, "https://gogotattoo.github.io/"+artist)
	}
}

func runHugo(source, destination, baseURL string) {
	cmd := exec.Command("hugo",
		"--destination", destination,
		"--source", source,
		"--baseUrl", baseURL,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
