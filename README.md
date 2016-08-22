Hacker News Reader
=========

Open multiple hacker news feed with your favorite browser using command line

#### Why?
Tired of clicking each link and for the sake of learning Golang. I read hacker news every morning, each every one of the link in the front page.

#### Installation and Usage
* Run this for get/install it:  
`go get -u github.com/bunchhieng/hnreader`
* From the root of a project/projects:
  * Run with default option to open 10 news with chrome:   
      `hnreader run`                
  * Run with option `-t = tabs` and `-b = browser`:     
      `hnreader run -t 7 -b "chrome"`

## Credits
* [urfave/cli](https://github.com/urfave/cli)
* [Fatih Arslan](https://github.com/fatih/color)
* [Martin Angers](https://github.com/PuerkitoBio/goquery)
* [skratchdot](https://github.com/skratchdot/open-golang)

#### License

The MIT License (MIT)
