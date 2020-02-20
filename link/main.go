package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	filepath := "ex1.html"
	filepath = "ex2.html"
	filepath = "ex3.html"
	filepath = "ex4.html"
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	links := scanForLinks(f)
	// fmt.Printf("%#v \n", links)
	fmt.Println(LinksToString(links))
}

type Link struct {
	Href    string
	Summary string
}


func LinksToString(links []Link) string {
	s, _ := json.MarshalIndent(links, "", "  ")
	return string(s)
}

func NewLink(href string, chunks []string) Link {
	summary := strings.Join(chunks, " ")
	return Link{Href: href, Summary: summary}
}

func GetAttrFromTokenizer(tokenizer *html.Tokenizer, attrName string) (value string, ok bool) {
	for {
		key, val, more := tokenizer.TagAttr()
		attr := string(key)
		if attr == "href" {
			return string(val), true
		}
		if !more {
			return "", false
		}
	}
}

func scanLinksFromFile(filepath string) []Link {
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	links := scanForLinks(f)
	return links	
}


// scanForLinks finds links in html implemented as anchor tags.
func scanForLinks(r io.Reader) []Link {
	links := make([]Link, 0)
	tokenizer := html.NewTokenizer(r)
	inLink := false
	href := ""
	chunks := make([]string, 0)
ParsingLoop:
	for {
		token := tokenizer.Next()
		switch token {
		case html.ErrorToken:
			msg := fmt.Sprintf("Error token encountered: %+v", token)
			log.Fatal(msg)
		case html.StartTagToken:
			tag, _ := tokenizer.TagName()
			tagname := string(tag)
			// fmt.Printf("Token: %v Tag name: %v Has attributes: %v \n", token, tagname, hasAttr)
			if tagname == "a" {
				inLink = true
				chunks = chunks[:0]
				var ok bool
				href, ok = GetAttrFromTokenizer(tokenizer, "href")
				if !ok {
					log.Print("No href found!")
				}
			}
		case html.EndTagToken:
			tag, _ := tokenizer.TagName()
			tagname := string(tag)
			// fmt.Printf("Token: %v Tag name: %v Has attributes: %v \n", token, tagname, hasAttr)
			if tagname == "html" {
				break ParsingLoop
			}
			if tagname == "a" {
				inLink = false
				links = append(links, NewLink(href, chunks))
				chunks = chunks[:0]
			}
		case html.TextToken:
			text := string(tokenizer.Text())
			// fmt.Printf("Token: %v Text: %v Length: %v \n", token, text, len(text))
			text = strings.TrimSpace(text)
			if text == "" {
				continue
			}
			if inLink {
				//fmt.Printf("Adding text to content %v\n", text)
				chunks = append(chunks, text)
			}
		default:
			// fmt.Printf("Token: %v\n", token)
		}
	}
	return links
}
