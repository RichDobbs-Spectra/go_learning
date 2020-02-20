package main

import (
	"fmt"
	"go_learning/link"
	"log"
	"os"
)

func main() {
	filepath := "ex1.html"
	filepath = "ex2.html"
	filepath = "../ex3.html"
	// filepath = "ex4.html"
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	links := link.ScanForLinks(f)
	// fmt.Printf("%#v \n", links)
	// fmt.Println(link.LinksToString(links))
	fmt.Println(links.AsDeclaration(true))
	expected := link.Links{
		link.Link{Href:"#", Summary:"Login"},
		link.Link{Href:"/lost", Summary:"Lost? Need help?"},
		link.Link{Href:"https://twitter.com/marcusolsson", Summary:"@marcusolsson"},
	}
	fmt.Println(expected.AsDeclaration(true))
	fmt.Println(expected.AsDeclaration(false))
}
