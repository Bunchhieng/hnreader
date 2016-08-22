package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
	"github.com/skratchdot/open-golang/open"
	"gopkg.in/urfave/cli.v2"
)

var Blue = color.New(color.FgBlue, color.Bold).SprintFunc()
var Yellow = color.New(color.FgYellow, color.Bold).SprintFunc()
var Red = color.New(color.FgRed, color.Bold).SprintFunc()

type logWriter struct{}

const (
	AppName        = "hnreader"
	AppVersion     = "v1.0"
	AppEmail       = "Bunchhieng@gmail.com"
	AppDescription = "Open multiple hacker news in your favorite browser with command line."
	AppAuthor      = "Bunchhieng Soth"
)

type App struct {
	Name, Version, Email, Description, Author string
}

func init() {
	log.SetFlags(0)
	log.SetOutput(new(logWriter))
}

func Init() *App {
	return &App{
		Name:        AppName,
		Version:     AppVersion,
		Description: AppDescription,
		Author:      AppAuthor,
		Email:       AppEmail,
	}
}

func (app *App) Information() {
	fmt.Println(Blue(app.Name) + " - " + Blue(app.Version))
	fmt.Println(Blue(app.Description) + "\n")
}

func (writer logWriter) Write(bytes []byte) (int, error) {
	return fmt.Print(Yellow("[") + time.Now().UTC().Format("15:04:05") + Yellow("]") + string(bytes))
}

func main() {
	app := Init()

	header := func() error {
		app.Information()
		gopath := os.Getenv("GOPATH")
		if gopath == "" {
			log.Fatal(Red("$GOPATH isn't set up properly."))
		}
		return nil
	}

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
				Usage:   "start hnreader with default option (10 news and chrome browser)",
				Flags: []cli.Flag{
					&cli.UintFlag{Name: "tabs", Value: 10, Aliases: []string{"t"}, Usage: "Specify value of tabs \t"},
					&cli.StringFlag{Name: "browser", Value: "chrome", Aliases: []string{"b"}, Usage: "Specify broswer \t"},
				},
				Action: func(c *cli.Context) error {
					return handle(runApp(int(c.Int("tabs")), string(c.String("browser"))))
				},
				Before: func(c *cli.Context) error {
					header()
					return nil
				},
			},
		},
	}

	cli.Run(os.Args)
}

func runApp(tabs int, browser string) error {
	news, err := getStories(tabs)
	// To store the keys in slice in sorted order
	var keys []int
	for k := range news {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	handle(err)

	for _, k := range keys {
		if k == tabs {
			break
		}
		err := open.RunWith(news[k], browser)
		if err != nil {
			fmt.Printf("%s is not found on this computer...", browser)
			os.Exit(1)
		}
	}
	return nil
}

func getStories(count int) (map[int]string, error) {
	news := make(map[int]string)
	// 30 news per page
	pages := count / 30
	for i := 0; i <= pages; i++ {
		doc, err := goquery.NewDocument("https://news.ycombinator.com/news?p=" + strconv.Itoa(pages))
		handle(err)
		doc.Find("a.storylink").Each(func(i int, s *goquery.Selection) {
			href, exist := s.Attr("href")
			if !exist {
				fmt.Println(Red("can't find any stories..."))
			}
			news[i] = href
		})
	}

	return news, nil
}

func handle(err error) error {
	if err != nil {
		fmt.Println(Red(err.Error()))
		return nil
	}
	return nil
}
