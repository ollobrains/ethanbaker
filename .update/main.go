package main

import (
    "errors"
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

I am a %v year old driven developer attending North Carolina State University
who's looking to master the many fields of computer science. My goal is to find
innovative solutions to compelling problems across various disciplines. In addition
to my knowledge of diverse coding modalities, I consistently work to improve my
knowledge of math, science and other changing technologies.

## Projects

My favorite projects I have crafted:

%v

## Tech Stack

Skills I have learned and honed over the years:

* Frontend
  * HTML, CSS, JS
  * TypeScript
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

// repoCountLimit is an optional limit for how many pinned repos to display.
const repoCountLimit = 10 // set to 0 for no limit

func main() {
    // 1. Determine GitHub username
    githubUser := os.Getenv(GITHUB_USER_ENV_KEY)
    if githubUser == "" {
        githubUser = "ethanbaker" // default user
    }

    // 2. Calculate user’s age
    ageString := calculateAge()

    // 3. Asynchronously fetch pinned repos
    var pinnedRepos string
    var pinnedErr error

    wg := sync.WaitGroup{}
    wg.Add(1)
    go func() {
        defer wg.Done()
        pinnedRepos, pinnedErr = findPinnedRepos(githubUser)
        if pinnedErr != nil {
            log.Printf("[WARN] Could not fetch pinned repos: %v\n", pinnedErr)
            pinnedRepos = fallbackPinnedMsg
        }
    }()
    wg.Wait()

    // 4. Generate final README content
    lastUpdated := time.Now().Format("Mon Jan 2 15:04 2006")
    readmeContent := fmt.Sprintf(README_TEMPLATE, ageString, pinnedRepos, lastUpdated)

    // 5. Write content to README.md
    if err := os.WriteFile("../README.md", []byte(readmeContent), 0644); err != nil {
        log.Fatalf("[FATAL] Failed to create README file: %v\n", err)
    }

    log.Printf("[INFO] README updated successfully for user: %s\n", githubUser)
}

// calculateAge determines your age based on a fixed birthdate.
func calculateAge() string {
    birthdate := time.Date(2003, time.February, 8, 0, 0, 0, 0, time.UTC)
    now := time.Now()

    years := now.Year() - birthdate.Year()
    // If today's day of year is before the birthdate day of year, subtract one year.
    if now.YearDay() < birthdate.YearDay() {
        years--
    }
    return fmt.Sprint(years)
}

// findPinnedRepos scrapes pinned repositories from a user's GitHub profile
// and returns formatted bullet items, or an error if unsuccessful.
func findPinnedRepos(githubUser string) (string, error) {
    profileURL := fmt.Sprintf("https://github.com/%s", githubUser)
    res, err := http.Get(profileURL)
    if err != nil {
        return "", fmt.Errorf("could not fetch: %s, error: %w", profileURL, err)
    }
    defer res.Body.Close()

    if res.StatusCode != http.StatusOK {
        return "", fmt.Errorf("non-200 status code: %d %s", res.StatusCode, res.Status)
    }

    doc, err := goquery.NewDocumentFromReader(res.Body)
    if err != nil {
        return "", fmt.Errorf("error loading HTML doc: %w", err)
    }

    var output strings.Builder
    counter := 0

    doc.Find(".pinned-item-list-item-content").Each(func(i int, s *goquery.Selection) {
        if repoCountLimit > 0 && counter >= repoCountLimit {
            return
        }
        title := s.Find("a > span").Text()
        link := s.Find("a").AttrOr("href", "")
        desc := s.Find(".pinned-item-desc").Text()

        // Trim leading/trailing whitespace from description
        desc = strings.TrimSpace(desc)

        if link != "" && title != "" {
            // Format a bullet item
            output.WriteString(fmt.Sprintf("* [%s](https://github.com%s): %s\n", title, link, desc))
            counter++
        }
    })

    // If no pinned repos found, return fallback
    if output.Len() == 0 {
        return "", errors.New("no pinned repositories found")
    }
    return output.String(), nil
}
