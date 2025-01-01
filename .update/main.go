package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// README_TEMPLATE is the template for the generated README.md content.
var README_TEMPLATE = `Hello! I'm Ethan.

## About Me

I am a %v year old driven developer attending North Carolina State University who's looking to master the many fields of computer science. My goal is to find innovative solutions to compelling problems across various disciplines. In addition to my knowledge of diverse coding modalities, I consistently work to improve my knowledge of math, science and other changing technologies.

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

// fallbackPinnedMsg is used if no pinned repos are found or if an error occurs.
const fallbackPinnedMsg = "* No pinned repositories found or there was an issue loading them.\n"

// GITHUB_USER_ENV_KEY is the environment variable key for overriding the GitHub username.
const GITHUB_USER_ENV_KEY = "GITHUB_USER"

// repoCountLimit is an optional limit to how many pinned repos we display (if desired).
const repoCountLimit = 10 // or set to 0 if unlimited

// calculateAge determines the user's age based on a fixed birth date.
func calculateAge() string {
	birthdate := time.Date(2003, time.February, 8, 0, 0, 0, 0, time.UTC)
	now := time.Now()

	years := now.Year() - birthdate.Year()
	if now.YearDay() < birthdate.YearDay() {
		years--
	}
	return fmt.Sprint(years)
}

// findPinnedRepos scrapes the GitHub user's pinned repositories from the GitHub profile.
// Returns a formatted string of repo info or a fallback message if none found.
func findPinnedRepos(githubUser string) string {
	profileURL := fmt.Sprintf("https://github.com/%s", githubUser)

	res, err := http.Get(profileURL)
	if err != nil {
		log.Printf("[WARN] Could not fetch pinned repos from: %s, err: %v\n", profileURL, err)
		return fallbackPinnedMsg
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Printf("[WARN] Non-200 status code: %d %s\n", res.StatusCode, res.Status)
		return fallbackPinnedMsg
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Printf("[WARN] Error loading HTML doc: %v\n", err)
		return fallbackPinnedMsg
	}

	output := ""
	counter := 0

	doc.Find(".pinned-item-list-item-content").Each(func(i int, s *goquery.Selection) {
		if repoCountLimit > 0 && counter >= repoCountLimit {
			return
		}
		title := s.Find("a > span").Text()
		link := s.Find("a").AttrOr("href", "invalid")
		desc := s.Find(".pinned-item-desc").Text()

		output += fmt.Sprintf("* [%v](https://github.com%v): %v\n", title, link, strings.TrimSpace(desc))
		counter++
	})

	if output == "" {
		return fallbackPinnedMsg
	}
	return output
}

func main() {
	// Overriding GitHub username from environment variable, if present
	githubUser := os.Getenv(GITHUB_USER_ENV_KEY)
	if githubUser == "" {
		githubUser = "ethanbaker" // default user
	}

	var pinnedRepos string
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		pinnedRepos = findPinnedRepos(githubUser)
	}()

	age := calculateAge()
	last := time.Now().Format("Mon Jan 2 15:04 2006")

	wg.Wait()

	readmeContent := fmt.Sprintf(README_TEMPLATE, age, pinnedRepos, last)

	err := os.WriteFile("../README.md", []byte(readmeContent), 0644)
	if err != nil {
		log.Fatalf("[FATAL] Failed to create README file: %v\n", err)
	}

	log.Printf("[INFO] README updated successfully for user: %s\n", githubUser)
}
