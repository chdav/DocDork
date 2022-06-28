package main

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var (
	domain       = flag.String("d", "", "Specify the domain to search.")
	filetype     = flag.String("e", "", "Specify the extension type. No argument will query every type. Options: docx, pptx, xlsx, pdf")
	metadataBool = flag.Bool("m", false, "\nUsed to output metadata. (Only works with open XML extensions)")
	filetypes    = []string{"docx", "pptx", "xlsx", "pdf"}
	fileCount    = 0
)

var versions = map[string]string{
	"16": "2016",
	"15": "2013",
	"14": "2010",
	"12": "2007",
	"11": "2003",
}

type properties struct {
	Creator        string `xml:"creator"`
	LastModifiedBy string `xml:"lastModifiedBy"`
	Application    string `xml:"Application"`
	Company        string `xml:"Company"`
	Version        string `xml:"AppVersion"`
}

func handler(i int, search *goquery.Selection) {
	url, exists := search.Find("a").Attr("href")
	if !exists {
		return
	}

	fmt.Println("[Doc Found] " + url)
	res, err := http.Get(url)
	if err != nil {
		return
	}
	fileCount++

	if *metadataBool == true {
		buf, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return
		}

		reader, err := zip.NewReader(bytes.NewReader(buf), int64(len(buf)))
		if err != nil {
			return
		}

		property, err := newProperties(reader)
		if err != nil {
			return
		}

		fmt.Println("[Metadata]\n  Creator\t\t: " + property.Creator + "\n  Last Modified By\t: " + property.LastModifiedBy + "\n  Company\t\t: " + property.Company + "\n  Application\t\t: " + property.Application + "\n  Version\t\t: " + property.getVersion())
	}

	defer res.Body.Close()
}

func dorkRequest(filetype string) {
	query := fmt.Sprintf(
		"site:%s && filetype:%s",
		*domain,
		filetype)

	search := fmt.Sprintf("https://www.bing.com/search?q=%s", url.QueryEscape(query))

	client := &http.Client{}

	req, err := http.NewRequest("GET", search, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("User-Agent", "Golang_Spider_Bot/3.0")

	res, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()

	if err != nil {
		return
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Panicln(err)
	}

	element := "html body.b_respl div#b_content main ol#b_results li.b_algo div.b_title h2"

	doc.Find(element).Each(handler)
}

func newProperties(reader *zip.Reader) (*properties, error) {
	var props properties
	for _, file := range reader.File {
		if file.Name == "docProps/core.xml" {
			err := decodeFile(file, &props)
			if err != nil {
				return nil, err
			}
		} else if file.Name == "docProps/app.xml" {
			err := decodeFile(file, &props)
			if err != nil {
				return nil, err
			}
		}
	}
	return &props, nil
}

func decodeFile(file *zip.File, properties interface{}) error {
	closer, err := file.Open()
	if err != nil {
		return err
	}
	defer closer.Close()
	err = xml.NewDecoder(closer).Decode(&properties)
	if err != nil {
		return err
	}
	return nil
}

func (props *properties) getVersion() string {
	splitProps := strings.Split(props.Version, ".")

	if len(splitProps) < 2 {
		return "Unk"
	}
	version, ok := versions[splitProps[0]]
	if !ok {
		return "Unk"
	}
	return version
}

func banner() {
	banner := `
	
		__ \                __ \                |    
		|   |   _ \    __|  |   |   _ \    __|  |  / 
		|   |  (   |  (     |   |  (   |  |       <     
		___/  \___/  \___| ____/  \___/  _|    _|\_\

			DocDork (v0.1)
				
		`

	fmt.Println(banner)
}

func main() {

	flag.Parse()
	if *domain == "" {
		fmt.Printf("Missing arguments.\n\n")
		flag.Usage()
		os.Exit(1)
	} else if *filetype == "pdf" && *metadataBool == true {
		fmt.Printf("Metadata not currently supported for PDFs.\n\n")
		flag.Usage()
		os.Exit(1)
	}

	banner()

	fmt.Println("\nDomain: *." + *domain)
	if *filetype == "" {
		for i := 0; i < len(filetypes); i++ {
			fmt.Println("[Info] Seaching for " + filetypes[i] + " extensions.")
			dorkRequest(filetypes[i])
		}
	} else {
		fmt.Println("[Info] Seaching for " + *filetype + " extensions.")
		dorkRequest(*filetype)
	}

	fmt.Println("[Results] DocDork complete at " + time.Now().UTC().Format("15:04:05 UCT, 02 Jan 2006"))
	fmt.Println("[Results] " + fmt.Sprint(fileCount) + " files found for *." + *domain + "!")
}
