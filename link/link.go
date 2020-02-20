package link

import (
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strings"

	"golang.org/x/net/html"
)

// Link stores information about an HTML link.
type Link struct {
	Href    string
	Summary string
}

// NewLink constructs a Link from the href and content as chunks.
func NewLink(href string, chunks []string) Link {
	summary := strings.Join(chunks, " ")
	return Link{Href: href, Summary: summary}
}

// Links is a collection of links. 
// It is declared as a type to allow methods to be implemented.
type Links []Link

// AsDeclaration outputs a string suitable for inclusion as a declaration in source code
func (links *Links) AsDeclaration(stripPackageName bool) string {
	v := reflect.ValueOf(*links)
	collectionTypeName := fmt.Sprint(v.Type())
	packageName := strings.Split(collectionTypeName, ".")[0]
	prepare := func(s string) string {
		if !stripPackageName {
			return s
		}
		return strings.ReplaceAll(s, packageName+".", "")
	}
	// v := reflect.ValueOf(*links)
	// collectionTypeName := fmt.Sprint(v.Type())
	var sb strings.Builder
	sb.WriteString(prepare(fmt.Sprintf("\t%T{\n", *links)))
	for _, link := range *links {
		sb.WriteString(prepare(fmt.Sprintf("\t\t%#v,\n", link)))
	}
	sb.WriteString(fmt.Sprintf("\t}"))
	return sb.String()
}

// ScanLinksFromFile loads links from a file in HTML format
func ScanLinksFromFile(filepath string) Links {
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	links := ScanForLinks(f)
	return links
}

type linkState struct {
	depth  int
	href   string
	chunks []string
	links  Links
}

// ScanForLinks finds links in html implemented as anchor tags.
// It is implemented as a state machine working on the golang.org/x/net/html html.Tokenizer
func ScanForLinks(r io.Reader) Links {
	//links := Links(make([]Link, 0))
	state := linkState{
		chunks: make([]string, 0),
		links:  Links(make([]Link, 0)),
	}
	tokenizer := html.NewTokenizer(r)
	for {
		token := tokenizer.Next()
		done := handleToken(token, tokenizer, &state)
		if done {
			break
		}
	}
	return state.links
}

// handleToken processes a single token while updating the link state.  
// It returns whether the processing should be done.
func handleToken(token html.TokenType, tokenizer *html.Tokenizer, state *linkState) bool {

	switch token {
	case html.ErrorToken:
		msg := fmt.Sprintf("Error token encountered: %+v", token)
		log.Fatal(msg)
	case html.StartTagToken:
		tag, _ := tokenizer.TagName()
		tagname := string(tag)
		// fmt.Printf("Token: %v Tag name: %v Has attributes: %v \n", token, tagname, hasAttr)
		if tagname == "a" {
			state.depth++
			if state.depth == 1 {
				state.chunks = state.chunks[:0]
				var ok bool
				state.href, ok = GetAttrFromTokenizer(tokenizer, "href")
				if !ok {
					log.Print("No href found!")
				}
			}
		}
	case html.EndTagToken:
		tag, _ := tokenizer.TagName()
		tagname := string(tag)
		// fmt.Printf("Token: %v Tag name: %v Has attributes: %v \n", token, tagname, hasAttr)
		if tagname == "html" {
			return true
		}
		if tagname == "a" {
			state.depth--
			if state.depth == 0 {
				link := NewLink(state.href, state.chunks)
				state.links = append(state.links, link)
				state.chunks = state.chunks[:0]
			}
		}
	case html.TextToken:
		text := string(tokenizer.Text())
		// fmt.Printf("Token: %v Text: %v Length: %v \n", token, text, len(text))
		text = strings.TrimSpace(text)
		if text == "" {
			return false
		}
		if state.depth == 1 {
			//fmt.Printf("Adding text to content %v\n", text)
			state.chunks = append(state.chunks, text)
		}
	default:
		// fmt.Printf("Token: %v\n", token)
	}
	return false
}

// GetAttrFromTokenizer finds the value of a specified attribute.
// It should only be called when the tokenizer has received a start or end tag token.
func GetAttrFromTokenizer(tokenizer *html.Tokenizer, attrName string) (value string, ok bool) {
	for {
		key, val, more := tokenizer.TagAttr()
		attr := string(key)
		if attr == "href" {
			return string(val), true
		}
		if !more {
			// The desired attribute has not been specified for the tag.
			return "", false
		}
	}
}
