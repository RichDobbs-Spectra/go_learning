package main

import "testing"

func linksAreEqual(a, b []Link) bool {
	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestEx1(t *testing.T) {
	actual := scanLinksFromFile("ex1.html")
	expected := []Link{
		{
			Href:    "/other-page",
			Summary: "A link to another page",
		},
	}
	if !linksAreEqual(actual, expected) {
		// t.Logf("Expected %v, got %v", LinksToString(expected), LinksToString(actual))
		t.Logf("Expected %#v\n got %#v\n", expected, actual)
		t.Fail()
	}
}

func TestEx2(t *testing.T) {
	actual := scanLinksFromFile("ex2.html")
	expected := []Link{
		{
			Href:    "https://www.twitter.com/joncalhoun",
			Summary: "Check me out on twitter"},
		{
			Href:    "https://github.com/gophercises",
			Summary: "Gophercises is on Github !"},
	}
	if !linksAreEqual(actual, expected) {
		t.Logf("Expected: \n%#v\nGot: \n%#v\n", expected, actual)
		t.Fail()
	}
}
