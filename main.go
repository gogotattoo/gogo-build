package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/gogotattoo/common/models"
	"github.com/gogotattoo/common/util"
	qrcode "github.com/skip2/go-qrcode"

	gia "github.com/ipfs/go-ipfs-api"
	flag "github.com/jteeuwen/go-pkg-optarg"
)

const (
	defaultArtist string = "gogo"
)

var artist string
var workType string
var shortName string
var lang string

func init() {
	flag.Header("General")
	flag.Add("a", "artist", "the name of the artist", defaultArtist)
	flag.Add("t", "type", "tattoo/design/...etc", "tattoo")
	flag.Add("s", "shortName", "link name of the work", "")
	flag.Add("i", "input", "artist/type/name format, instead of -a -t -s", "")
	flag.Add("l", "language", "language, default is en, supports ru", "")

	for opt := range flag.Parse() {
		switch opt.Name {
		case "artist":
			artist = opt.String()
		case "type":
			workType = opt.String()
		case "shortName":
			shortName = opt.String()
		case "input":
			atn := opt.String()
			artist = strings.Split(atn, "/")[0]
			workType = strings.Split(atn, "/")[1]
			shortName = strings.Split(atn, "/")[2]
		case "language":
			lang = opt.String()
		}
	}
}

func main() {
	if len(os.Args) > 1 {
		renderIntoArticle()
		return
	}
	renderSite("gogo.tattoo", "http")
	renderSite("gogotattoo.github.io", "https")
}

func renderIntoArticle() {
	path := "../" + artist + "/content/" + workType + "/" + shortName + ".md"
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}
	//fmt.Print(string(dat))
	tomlStr, _ := util.ExtractTomlStr(bytes.NewReader(dat))
	var work models.Artwork
	toml.Unmarshal([]byte(tomlStr), &work)

	for i, tag := range work.Tags {
		work.Tags[i] = strings.Replace(tag, " ", "-", -1)
	}

	if len(lang) > 0 {
		renderArticleTemplate("templates/template.ru.html", work)
	} else {
		renderArticleTemplate("templates/template.html", work)
	}

}

func renderArticleTemplate(tmpl string, work models.Artwork) {
	funcMap := template.FuncMap{
		"title": strings.Title,
		"fdate": func(date string) string {
			t, e := time.Parse(time.RFC3339, date)
			if e != nil {
				panic(e)
			}
			return t.Format("2006/01/02")
		},
	}
	t, err := template.New(strings.Split(tmpl, "/")[1]).Funcs(funcMap).ParseFiles(tmpl)
	if nil != err {
		panic(err)
	}
	work.Link = artist
	if len(lang) > 0 {
		work.Link += "/" + lang
	}
	work.Link += "/" + workType + "/" + shortName

	qrLink := "http://gogo.tattoo/" + work.Link + "?utm_medium=qrcode&utm_source="
	source := "steemit"
	if lang == "ru" {
		source = "golos"
	}
	qrLink += source
	e := qrcode.WriteFile(qrLink, qrcode.Medium, 256, "qr_"+shortName+"_"+source+".png")
	var qrHash string
	if e == nil {
		sh := gia.NewShell("localhost:5001")
		qrHash, _ = sh.AddDir("qr_" + shortName + "_" + source + ".png")
	}
	if err := t.Execute(os.Stdout, struct {
		Work   models.Artwork
		Artist string
		QRhash string
	}{work, artist, qrHash}); err != nil {
		fmt.Println(err)
	}
}
func renderSite(out, protocol string) {
	os.Remove("gogo.tattoo")
	gopath := os.Getenv("GOPATH")
	gogotattooPrefix := "/src/github.com/gogotattoo/"
	artists := []string{"xizi", "gogo", "aid", "kate", "klimin", "jiaye"}
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
		addNewLocalizedFile(fName, ".zh-hant.md")
		addNewLocalizedFile(fName, ".ru.md")
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
