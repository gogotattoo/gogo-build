package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	renderSite("gogo.tattoo", "http")
	//renderSite("gogotattoo.github.io", "https")

}

func renderSite(out, protocol string) {
	gopath := os.Getenv("GOPATH")
	gogotattooPrefix := "/src/github.com/gogotattoo/"
	artists := []string{"xizi"}
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

func addLocalizedFiles(path string, info os.FileInfo, err error) error {
	newPath := path[:len(path)-3] + ".zh.md"
	er := copyFile(path, newPath)
	if er == nil {
		filesToDelete = append(filesToDelete, newPath)
	}
	return nil
}

func renderArtistSite(source, destination, baseURL string) {
	filepath.Walk(source+"/content/design", addLocalizedFiles)
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
