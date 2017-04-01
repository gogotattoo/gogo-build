package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	renderSite("gogo.tattoo", "http")
	renderSite("gogotattoo.github.io", "https")
}

func renderSite(out, protocol string) {
	os.Remove("gogo.tattoo")
	gopath := os.Getenv("GOPATH")
	gogotattooPrefix := "/src/github.com/gogotattoo/"
	artists := []string{"xizi", "gogo", "aid"}
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	renderMainSite(gopath+gogotattooPrefix+"/gogo.tattoo", dir+"/"+out, protocol+"://"+out+"/")
	for _, artist := range artists {
		renderArtistSite(gopath+gogotattooPrefix+artist, dir+"/"+out+"/"+artist, protocol+"://"+out+"/"+artist)
	}
}
func renderMainSite(source, destination, baseURL string) {
	cmd := exec.Command("hugo",
		"--destination", destination,
		"--source", source,
		"--baseUrl", baseURL,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

var filesToDelete []string
var filesToCopy []string

func collectLocalizedFiles(path string, info os.FileInfo, err error) error {
	if strings.HasSuffix(path, ".md") {
		filesToCopy = append(filesToCopy, path)
	}
	return nil
}

func addNewLocalizedFile(path, res string) {
	if len(path) == 0 {
		return
	}
	newPath := path[:len(path)-3] + res
	er := copyFile(newPath, path)
	if er == nil {
		filesToDelete = append(filesToDelete, newPath)
	}
}
func renderArtistSite(source, destination, baseURL string) {
	filesToDelete = make([]string, 100)
	filesToCopy = make([]string, 100)
	filepath.Walk(source+"/content/", collectLocalizedFiles)
	fmt.Println(filesToCopy)
	for _, fName := range filesToCopy {
		//fmt.Println(fName)
		addNewLocalizedFile(fName, ".zh.md")
		//addNewLocalizedFile(fName, ".zh-hant.md")
		//addNewLocalizedFile(fName, ".ru.md")
	}
	renderMainSite(source, destination, baseURL)
	for _, fName := range filesToDelete {
		os.Remove(fName)
	}
}
func copyFile(dst, src string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	cerr := out.Close()
	if err != nil {
		return err
	}
	return cerr
}
