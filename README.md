# Hacker News Reader

[![Go Report Card](https://goreportcard.com/badge/github.com/Bunchhieng/hnreader)](https://goreportcard.com/report/github.com/Bunchhieng/hnreader) [![Build Status](https://travis-ci.org/Bunchhieng/hnreader.svg?branch=master)](https://travis-ci.org/Bunchhieng/hnreader)
[![Maintainability](https://api.codeclimate.com/v1/badges/ba5c7736f364c04b562c/maintainability)](https://codeclimate.com/github/Bunchhieng/hnreader/maintainability)
![GitHub](https://img.shields.io/github/license/mashape/apistatus.svg)

Open multiple news feed with your favorite browser using command line.

#### Why?

Stay up to update with new technology and waste all your morning reading everything on front page of popular programming sites.

#### Installation and Usage

- Run this for get/install it:  
  `go get -u github.com/Bunchhieng/hnreader`

- Install golang:
  - `brew install go` on Mac
  - Or download the binary from [here](https://golang.org/dl/)
- Make sure you have GOPATH set up properly:  
  `export GOPATH=/path/to/your/go/workspace`  
  `export PATH=$GOPATH/bin:$PATH`

- From the root of a project:

  - Run with default option to open 10 news with chrome:  
     `hnreader r`
  - Example on how to run with option `-t = tabs`, `-s = source` and `-b = browser`:
    - `hnreader r -t 31 -b "firefox"`
    - `hnreader r -b "brave" -s "reddit"`
    - `hnreader r -b "firefox" -s "reddit" -t 20`
  - bash alias is your friend if you're going to run the same command every morning
    - `alias hnr='hnreader r -b "firefox" -s "reddit" -t 30'`

#### Want to contribute?

Please see the [CONTRIBUTING.md](CONTRIBUTING.md)

#### Credits

- [urfave/cli](https://github.com/urfave/cli)
- [Fatih Arslan](https://github.com/fatih/color)
- [Martin Angers](https://github.com/PuerkitoBio/goquery)
- [skratchdot](https://github.com/skratchdot/open-golang)

#### License

The MIT License (MIT)
