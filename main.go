package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"sort"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
	"github.com/jzelinskie/geddit"
	"github.com/skratchdot/open-golang/open"
	"github.com/texttheater/golang-levenshtein/levenshtein"
	"gopkg.in/urfave/cli.v2"
)

// App information and constants
const (
	AppName        = "hnreader"
	AppVersion     = "v1.1"
	AppAuthor      = "Bunchhieng Soth"
	AppEmail       = "Bunchhieng@gmail.com"
	AppDescription = "Open multiple hacker news in your favorite browser with command line."
	HackerNewsURL  = "https://news.ycombinator.com/news?p="
)

// Colors for console output
var blue = color.New(color.FgBlue, color.Bold).SprintFunc()
var yellow = color.New(color.FgYellow, color.Bold).SprintFunc()
var red = color.New(color.FgRed, color.Bold).SprintFunc()

type logWriter struct{}

// App contains author information
type App struct {
	Name, Version, Email, Description, Author string
}

// Fetcher retrieves stories from a source.
type Fetcher interface {
	Fetch(count int) (map[int]string, error)
}

// HackerNewsSource fetches new stories from news.ycombinator.com.
type HackerNewsSource struct{}

func (hn *HackerNewsSource) Fetch(count int) (map[int]string, error) {
	news := make(map[int]string)
	// 30 news per page
	pages := count / 30
	for i := 0; i <= pages; i++ {
		resp, err := http.Get(HackerNewsURL + strconv.Itoa(pages))
		handleError(err)
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		handleError(err)
		doc.Find("a.storylink").Each(func(i int, s *goquery.Selection) {
			href, exist := s.Attr("href")
			if !exist {
				fmt.Println(red("can't find any stories..."))
			}
			news[i] = href
		})
	}

	return news, nil
}

// RedditSource fetches new stories from reddit.com/r/programming.
type RedditSource struct{}

func (rs *RedditSource) Fetch(count int) (map[int]string, error) {
	news := make(map[int]string)

	s := geddit.NewSession(fmt.Sprintf("desktop:com.github.Bunchhieng.%s:%s", AppName, AppVersion))
	subs, err := s.SubredditSubmissions(
		"programming",
		geddit.HotSubmissions,
		geddit.ListingOptions{
			Count: count,
			Limit: count,
		},
	)
	handleError(err)

	for i, sub := range subs {
		news[i] = sub.URL
	}
	return news, nil
}

// Init initalizes the app
func Init() *App {
	return &App{
		Name:        AppName,
		Version:     AppVersion,
		Description: AppDescription,
		Author:      AppAuthor,
		Email:       AppEmail,
	}
}

// Information prints out app information
func (app *App) Information() {
	fmt.Println(blue(app.Name) + " - " + blue(app.Version))
	fmt.Println(blue(app.Description) + "\n")
}

func (writer logWriter) Write(bytes []byte) (int, error) {
	return fmt.Print(yellow("[") + time.Now().UTC().Format("15:04:05") + yellow("]") + string(bytes))
}

//RunApp opens a browser with input tabs count
func RunApp(tabs int, browser string, src Fetcher) error {
	news, err := src.Fetch(tabs)
	handleError(err)

	browser = findBrowser(browser)

	// To store the keys in slice in sorted order
	var keys []int
	for k := range news {
		keys = append(keys, k)
	}
	// Sort map keys
	sort.Ints(keys)

	for _, k := range keys {
		if k == tabs {
			break
		}

		var err error
		if browser == "" {
			fmt.Println(red("Trying default browser..."))
			err = open.Run(news[k])
		} else {
			err = open.RunWith(news[k], browser)
			if err != nil {
				fmt.Printf(red("%s is not found on this computer, trying default browser...\n"), browser)
				err = open.Run(news[k])
			}
		}

		if err != nil {
			os.Exit(1)
		}
	}
	return nil
}

func findBrowser(target string) string {
	if target == "" {
		return ""
	}
	browsers := []string{"google", "chrome", "mozilla", "firefox", "brave"}
	shortest := -1
	word := ""
	for _, browser := range browsers {
		distance := levenshtein.DistanceForStrings([]rune(browser), []rune(target), levenshtein.DefaultOptions)
		if distance == 0 {
			word = browser
			break
		}
		if distance <= shortest || shortest < 0 {
			shortest = distance
			word = browser
		}
	}

	return getBrowserNameByOS(word)
}

func getBrowserNameByOS(k string) string {
	browser := ""
	switch k {
	case "google", "chrome":
		switch runtime.GOOS {
		case "darwin":
			browser = "Google Chrome"
		}
	case "mozilla", "firefox":
		switch runtime.GOOS {
		case "darwin":
			browser = "Firefox"
		}
	case "brave":
		switch runtime.GOOS {
		case "darwin":
			browser = "Brave"
		}
	}

	return browser
}

// checkGoPath checks for GOPATH
func checkGoPath() error {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		log.Fatal(red("$GOPATH isn't set up properly..."))
	}
	return nil
}

// handleError go convention
func handleError(err error) error {
	if err != nil {
		fmt.Println(red(err.Error()))
	}
	return nil
}

func init() {
	log.SetFlags(0)
	log.SetOutput(new(logWriter))
}

func main() {
	app := Init()

	cli := &cli.App{
		Name:    app.Name,
		Version: app.Version,
		Authors: []*cli.Author{
			{
				Name:  app.Author,
				Email: app.Email,
			},
		},
		Usage: app.Description,
		Commands: []*cli.Command{
			{
				Name:    "run",
				Aliases: []string{"r"},
				Usage:   "Start hnreader with default option (10 news and chrome browser)",
				Flags: []cli.Flag{
					&cli.UintFlag{
						Name:    "tabs",
						Value:   10,
						Aliases: []string{"t"},
						Usage:   "Specify value of tabs\t",
					},
					&cli.StringFlag{
						Name:    "browser",
						Value:   "",
						Aliases: []string{"b"},
						Usage:   "Specify broswer\t",
					},
					&cli.StringFlag{
						Name:    "source",
						Value:   "hn",
						Aliases: []string{"s"},
						Usage:   "Specify news source (one of \"hn\", \"reddit\")\t",
					},
				},
				Action: func(c *cli.Context) error {
					var src Fetcher

					switch c.String("source") {
					case "hn":
						src = new(HackerNewsSource)
					case "reddit":
						src = new(RedditSource)
					default:
						return handleError(fmt.Errorf("invalid source: %s", c.String("source")))
					}

					return handleError(RunApp(c.Int("tabs"), c.String("browser"), src))
				},
				Before: func(c *cli.Context) error {
					app.Information()
					checkGoPath()
					return nil
				},
			},
		},
	}

	cli.Run(os.Args)
}
