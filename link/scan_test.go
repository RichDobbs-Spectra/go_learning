package link

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


func handleFileTest(t *testing.T, filePath string, expected Links) {
	actual := ScanLinksFromFile(filePath)
	if !linksAreEqual(actual, expected) {
		t.Logf("Expected: \n%v\nGot: \n%v\n", expected.AsDeclaration(true), actual.AsDeclaration(true))
		t.Fail()
	}
}



func TestEx1(t *testing.T) {
	expected := Links{
		{
			Href:    "/other-page",
			Summary: "A link to another page",
		},
	}
	handleFileTest(t, "ex1.html", expected)
}

func TestEx2(t *testing.T) {
	expected := Links{
		{
			Href:    "https://www.twitter.com/joncalhoun",
			Summary: "Check me out on twitter"},
		{
			Href:    "https://github.com/gophercises",
			Summary: "Gophercises is on Github !"},
	}
	handleFileTest(t, "ex2.html", expected)
}

func TestEx3(t *testing.T) {
	expected := Links{
		{Href: "#", Summary: "Login"},
		{Href: "/lost", Summary: "Lost? Need help?"},
		{Href: "https://twitter.com/marcusolsson", Summary: "@marcusolsson"},
	}
	handleFileTest(t, "ex3.html", expected)
}


func TestEx4(t *testing.T) {
	expected := Links{
		Link{Href: "/dog-cat", Summary: "dog cat"},
	}
	handleFileTest(t, "ex4.html", expected)
}


func TestNestLink(t *testing.T) {
	expected := Links{
		Link{Href: "#", Summary: "Something here and here"},
	}
	handleFileTest(t, "nestLink.html", expected)
}

