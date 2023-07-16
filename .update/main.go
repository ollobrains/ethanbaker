package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Template for the README
var README_TEMPLATE = `Hello! I'm Ethan.
## About Me

I am a %v year old driven developer attending North Carolina State University who’s looking to master the many fields of computer science. My goal is to find innovative solutions to compelling problems across various disciplines. In addition to my knowledge of diverse coding modalities, I consistently work to improve my knowledge of math, science and other changing technologies.

## Projects

My favorite projects I have crafted:

%v

## Tech Stack

Skills I have learned and honed over the years:

* Frontend
  * HTML, CSS, JS
  * Typescript
  * Angular/AngularJS
* Backend
  * Golang
  * NodeJS
  * REST APIs
  * Web sockets
  * MySQL
* Terminal-Based Development
  * Bash
  * C, C++
* Data Science
  * Artificial Intelligence
  * Statistics
  * Python
* Development Processes
  * Git
  * Agile
  * Bash
* STEM
  * Mathematics
  * Physics
  * Biology
  * Chemistry

## Experience

Experiences and opportunities I've had throughout my journey:

* Software Development Intern for [Bio-Rad Laboratories](https://www.bio-rad.com/)
* Developer for [Institute for Transportation Research and Education](https://itre.ncsu.edu/)
* Middle School Coding Teacher for the [Community School of Davidson](https://www.csdspartans.org/)

## Contact

Want to get in touch? You can find me anywhere below:

* [contact@ethanbaker.dev](mailto:contact@ethanbaker.dev)
* [LinkedIn](https://www.linkedin.com/in/ethan-baker-802b2a183)
* [GitHub](https://github.com/ethanbaker)

<sub>Last updated on: %v</sub>
`

// Calculate my age so the repo stays updated
func calculateAge() string {
	// Get the dates to compare
	birthdate := time.Date(2003, time.February, 8, 0, 0, 0, 0, time.UTC)
	now := time.Now()

	// Find the age
	years := now.Year() - birthdate.Year()
	if now.YearDay() < birthdate.YearDay() {
		years--
	}

	return fmt.Sprint(years)
}

// Get a list of all my pinned repos
func findPinnedRepos() string {
	// Request the HTML page
	res, err := http.Get("https://github.com/ethanbaker")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("Status code error: %d %s\n", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalf("Error loading html document: %v\n", err)
	}

	output := ""

	// Find the pinned repos
	doc.Find(".pinned-item-list-item-content").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title, link, and description
		title := s.Find("a > span").Text()
		link := s.Find("a").AttrOr("href", "invalid")
		desc := s.Find(".pinned-item-desc").Text()

		output += fmt.Sprintf("* [%v](https://github.com%v): %v\n", title, link, strings.TrimSpace(desc))
	})

	return output
}

func main() {
	// Get the different components of the README and format it
	age := calculateAge()
	pinned := findPinnedRepos()
	last := time.Now().Format("Mon Jan 2 15:04 2006")

	readme := fmt.Sprintf(README_TEMPLATE, age, pinned, last)

	// Create a new file for the README
	err := os.WriteFile("../README.md", []byte(readme), 0644)
	if err != nil {
		log.Fatalf("Failed to create README file: %v\n", err)
	}
}
