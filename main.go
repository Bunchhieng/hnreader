package main

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
	"github.com/jzelinskie/geddit"
	"github.com/skratchdot/open-golang/open"
	cli "gopkg.in/urfave/cli.v2"
)

// App information and constants
const (
	AppName        = "hnreader"
	AppVersion     = "v1.1"
	AppAuthor      = "Bunchhieng Soth"
	AppEmail       = "Bunchhieng@gmail.com"
	AppDescription = "Open multiple tech news feeds in your favorite browser through the command line."
	HackerNewsURL  = "https://news.ycombinator.com/news?p="
	LobstersURL    = "https://lobste.rs"
	DZoneURL       = "http://feeds.dzone.com/home"
	DevToURL       = "https://dev.to/feed"
	SteemItURL     = "https://api.steemjs.com/getState?path=/trending/technology&scope=content"
)

// Supported operating systems (GOOS)
const (
	OSDarwin  = "darwin"
	OSLinux   = "linux"
	OSWindows = "windows"
)

// Colors for console output
var blue = color.New(color.FgBlue, color.Bold).SprintFunc()
var yellow = color.New(color.FgYellow, color.Bold).SprintFunc()
var red = color.New(color.FgRed, color.Bold).SprintFunc()

// Rss decode RSS xml
type Rss struct {
	Item []RssItem `xml:"channel>item"`
}

// RssItem item with link to news
type RssItem struct {
	Link string `xml:"link"`
}

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

// Fetch gets news from the HackerNews
func (hn *HackerNewsSource) Fetch(count int) (map[int]string, error) {
	news := make(map[int]string)
	pages := count / 30
	
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	for i := 0; i <= pages; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, HackerNewsURL+strconv.Itoa(pages), nil)
		if err != nil {
			return nil, fmt.Errorf("creating request: %w", err)
		}
		
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("fetching page %d: %w", i, err)
		}
		defer resp.Body.Close()
		
		// Add response status code check
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unexpected status code %d from page %d", resp.StatusCode, i)
		}
		
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			handleError(err)
			continue
		}

		doc.Find("span.titleline > a").Each(func(i int, s *goquery.Selection) {
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

// Fetch gets news from the Reddit
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

	if err != nil {
		return news, err
	}

	for i, sub := range subs {
		news[i] = sub.URL
	}

	return news, nil
}

// LobstersSource fetches new stories from https://lobste.rs
type LobstersSource struct{}

// Fetch gets news from the Lobsters
func (l *LobstersSource) Fetch(count int) (map[int]string, error) {
	offset := float64(count) / float64(25)
	pages := int(math.Ceil(offset))
	news := make(map[int]string)
	newsIndex := 0

	for p := 1; p <= pages; p++ {
		url := fmt.Sprintf("%s/page/%d", LobstersURL, p)
		resp, err := http.Get(url)
		if err != nil {
			handleError(err)
			continue
		}

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			handleError(err)
			continue
		}

		doc.Find(".link a.u-url").Each(func(_ int, s *goquery.Selection) {
			href, exist := s.Attr("href")
			if !exist {
				fmt.Println(red("can't find any stories..."))
			}

			if newsIndex >= count {
				return
			}

			// if internal link
			if strings.HasPrefix(href, "/") {
				href = LobstersURL + href
			}

			news[newsIndex] = href
			newsIndex++
		})

		resp.Body.Close()
	}

	return news, nil
}

// DZoneSource fetches latest stories from http://feeds.dzone.com/home
type DZoneSource struct{}

// Fetch gets news from the DZone
func (l *DZoneSource) Fetch(count int) (map[int]string, error) {
	news := make(map[int]string)

	resp, err := http.Get(DZoneURL)
	if err != nil {
		return news, err
	}

	defer resp.Body.Close()

	doc := Rss{}
	d := xml.NewDecoder(resp.Body)

	if err := d.Decode(&doc); err != nil {
		return news, err
	}

	for i, item := range doc.Item {
		if i >= count {
			break
		}

		news[i] = item.Link
	}

	return news, nil
}

// DevToSource fetches latest stories from https://dev.to/
type DevToSource struct{}

// Fetch gets news from the Dev.To
func (l *DevToSource) Fetch(count int) (map[int]string, error) {
	news := make(map[int]string)

	resp, err := http.Get(DevToURL)
	if err != nil {
		return news, err
	}

	defer resp.Body.Close()

	doc := Rss{}
	d := xml.NewDecoder(resp.Body)

	if err := d.Decode(&doc); err != nil {
		return news, err
	}

	for i, item := range doc.Item {
		if i >= count {
			break
		}

		news[i] = item.Link
	}

	return news, nil
}

// SteemItSource fetches latest news from SteemIt
type SteemItSource struct{}

// Fetch gets data from SteemIt
func (l *SteemItSource) Fetch(count int) (map[int]string, error) {
	news := make(map[int]string)

	resp, err := http.Get(SteemItURL)
	if err != nil {
		return news, err
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return news, err
	}

	m := make(map[string]map[string]interface{})
	err = json.Unmarshal(data, &m)
	if err != nil {
		return news, err
	}

	idx := 0
	for _, value := range m {
		if idx >= count {
			break
		}

		str, ok := value["url"].(string)
		if ok {
			news[idx] = str
			idx++
		}
	}

	return news, nil
}

// Init initializes the app
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
	fmt.Println(blue(app.Description))
}

// Write logs to console
func (writer logWriter) Write(bytes []byte) (int, error) {
	return fmt.Print(yellow("[") + time.Now().UTC().Format("15:04:05") + yellow("]") + string(bytes))
}

// RunApp opens a browser with input tabs count
func RunApp(tabs int, browser string, src Fetcher) error {
	if tabs <= 0 {
		return fmt.Errorf("invalid number of tabs: %d", tabs)
	}
	
	news, err := src.Fetch(tabs)
	if err != nil {
		return fmt.Errorf("fetching news: %w", err)
	}
	
	if len(news) == 0 {
		return fmt.Errorf("no news items found")
	}
	
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
	
	normalizedTarget := strings.ToLower(strings.TrimSpace(target))
	browsers := map[string][]string{
		"chrome":  {"google", "chrome", "google-chrome"},
		"firefox": {"mozilla", "firefox"},
		"brave":   {"brave"},
		"safari":  {"safari"},
		"opera":   {"opera"},
	}
	
	// Direct match
	for browser, aliases := range browsers {
		for _, alias := range aliases {
			if normalizedTarget == alias {
				return getBrowserNameByOS(browser, runtime.GOOS)
			}
		}
	}
	
	// Partial match
	for browser, aliases := range browsers {
		for _, alias := range aliases {
			if strings.Contains(alias, normalizedTarget) || strings.Contains(normalizedTarget, alias) {
				return getBrowserNameByOS(browser, runtime.GOOS)
			}
		}
	}
	
	// If no match found, return empty string
	return ""
}

// getGoogleChromeNameForOS
func getGoogleChromeNameForOS(os string) string {
	switch os {
	case OSDarwin:
		return "Google Chrome"
	case OSLinux:
		return "google-chrome"
	case OSWindows:
		return "chrome"
	}
	return ""
}

// getFirefoxNameForOS
func getFirefoxNameForOS(os string) string {
	switch os {
	case OSDarwin:
		return "Firefox"
	case OSLinux:
		return "firefox"
	case OSWindows:
		return "firefox"
	}
	return ""
}

// getBraveNameForOS
func getBraveNameForOS(os string) string {
	switch os {
	case OSDarwin:
		return "Brave"
	case OSLinux:
		return "brave"
	case OSWindows:
		return "brave"
	}
	return ""
}

// getSafariNameForOS
func getSafariNameForOS(os string) string {
	switch os {
	case OSDarwin:
		return "Safari"
	case OSLinux:
		return "safari"
	case OSWindows:
		return "safari"
	}
	return ""
}

// getOperaNameForOS
func getOperaNameForOS(os string) string {
	switch os {
	case OSDarwin:
		return "Opera"
	case OSLinux:
		return "opera"
	case OSWindows:
		return "opera"
	}
	return ""
}

// getBrowserNameByOS normilizes browser name
func getBrowserNameByOS(browserFromCLI, os string) string {
	switch browserFromCLI {
	case "google", "chrome":
		return getGoogleChromeNameForOS(os)
	case "mozilla", "firefox":
		return getFirefoxNameForOS(os)
	case "brave":
		return getBraveNameForOS(os)
	case "safari":
		return getSafariNameForOS(os)
	case "opera":
		return getOperaNameForOS(os)
	}
	return ""
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

// removeIndex removes specific index from the slice
func removeIndex(slice []cli.Flag, s int) []cli.Flag {
	return append(slice[:s], slice[s+1:]...)
}

// getAllFlags return all flags for the command line
func getAllFlags(includeSource bool) []cli.Flag {
	flags := []cli.Flag{
		&cli.UintFlag{
			Name:    "tabs",
			Value:   10,
			Aliases: []string{"t"},
			Usage:   "Specify number of tabs\t",
		},
		&cli.StringFlag{
			Name:    "browser",
			Value:   "",
			Aliases: []string{"b"},
			Usage:   "Specify browser (one of \"chrome\", \"brave\", \"safari\", \"firefox\", \"opera\")\t",
		},
		&cli.StringFlag{
			Name:    "source",
			Value:   "hn",
			Aliases: []string{"s"},
			Usage:   "Specify news source (one of \"hn\", \"reddit\", \"lobsters\", \"dzone\", \"devto\", \"steemit\")\t",
		},
	}

	if !includeSource {
		flags = removeIndex(flags, 2)
	}

	return flags
}

// getAllActions return all action for the command line
func getAllActions(c *cli.Context) error {
	cfg := NewConfig()
	
	// Validate tabs
	tabs := c.Int("tabs")
	if tabs <= 0 || tabs > cfg.MaxTabs {
		return fmt.Errorf("invalid number of tabs (must be between 1 and %d)", cfg.MaxTabs)
	}
	
	// Get and validate source
	var srcName string
	if c.Command.Name == "random" {
		rand.Seed(time.Now().UnixNano())
		srcName = cfg.Sources[rand.Intn(len(cfg.Sources))]
	} else {
		srcName = c.String("source")
	}
	
	src, err := validateSource(srcName, cfg)
	if err != nil {
		return err
	}
	
	// Run the app
	return RunApp(tabs, c.String("browser"), src)
}

// Add configuration struct
type Config struct {
	MaxTabs     int
	DefaultTabs int
	Sources     []string
	Browsers    []string
	HTTPTimeout time.Duration
}

func NewConfig() *Config {
	return &Config{
		MaxTabs:     50,
		DefaultTabs: 10,
		Sources:     []string{"hn", "reddit", "lobsters", "dzone", "devto", "steemit"},
		Browsers:    []string{"chrome", "firefox", "brave", "safari", "opera"},
		HTTPTimeout: 10 * time.Second,
	}
}

// Add source validation
func validateSource(srcName string, cfg *Config) (Fetcher, error) {
	srcName = strings.ToLower(strings.TrimSpace(srcName))
	
	for _, validSource := range cfg.Sources {
		if srcName == validSource {
			switch srcName {
			case "hn":
				return new(HackerNewsSource), nil
			case "reddit":
				return new(RedditSource), nil
			case "lobsters":
				return new(LobstersSource), nil
			case "dzone":
				return new(DZoneSource), nil
			case "devto":
				return new(DevToSource), nil
			case "steemit":
				return new(SteemItSource), nil
			}
		}
	}
	
	return nil, fmt.Errorf("invalid source: %s", srcName)
}

// Add structured logging
type Logger struct {
	*log.Logger
}

func NewLogger() *Logger {
	return &Logger{
		Logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.Printf("[INFO] "+format, v...)
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.Printf("[ERROR] "+format, v...)
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
				Flags:   getAllFlags(true),
				Action:  getAllActions,
				Before: func(c *cli.Context) error {
					app.Information()
					checkGoPath()
					return nil
				},
			},
			{
				Name:    "random",
				Aliases: []string{"rr"},
				Usage:   "Start hnreader with a randomized source of news",
				Flags:   getAllFlags(false),
				Action:  getAllActions,
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
