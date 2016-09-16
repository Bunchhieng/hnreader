Hacker News Reader
=========
[![Go Report Card](https://goreportcard.com/badge/github.com/Bunchhieng/hnreader)](https://goreportcard.com/report/github.com/Bunchhieng/hnreader)  [![Build Status](https://travis-ci.org/Bunchhieng/hnreader.svg?branch=master)](https://travis-ci.org/Bunchhieng/hnreader)

Open multiple hacker news feed with your favorite browser using command line.

#### Why?
I tend to read all news on the front page of Hacker News. This allow me to run one command and make life easy.

#### Installation and Usage
* Run this for get/install it:  
`go get -u github.com/Bunchhieng/hnreader`

* Make sure you have GOPATH set up properly:    
`export GOPATH=/path/to/your/go/workspace`                    
`export PATH=$GOPATH/bin:$PATH`

* From the root of a project:
  * Run with default option to open 10 news with chrome:   
      `hnreader run`                
  * Run with option `-t = tabs` and `-b = browser`:     
      `hnreader run -t 31 -b "firefox"`

#### Credits
* [urfave/cli](https://github.com/urfave/cli)
* [Fatih Arslan](https://github.com/fatih/color)
* [Martin Angers](https://github.com/PuerkitoBio/goquery)
* [skratchdot](https://github.com/skratchdot/open-golang)

#### License

The MIT License (MIT)
