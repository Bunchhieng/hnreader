# hnreader

Stay up to date with the latest news in technology from your favourite programming sites, all aggregated into one feed.

hnreader (Hackernews Reader) allows you to open tech news feeds in your favorite browser from the command line.

[![Go Report Card](https://goreportcard.com/badge/github.com/Bunchhieng/hnreader)](https://goreportcard.com/report/github.com/Bunchhieng/hnreader) [![Build Status](https://travis-ci.org/Bunchhieng/hnreader.svg?branch=master)](https://travis-ci.org/Bunchhieng/hnreader)
[![Maintainability](https://api.codeclimate.com/v1/badges/ba5c7736f364c04b562c/maintainability)](https://codeclimate.com/github/Bunchhieng/hnreader/maintainability)
![GitHub](https://img.shields.io/github/license/mashape/apistatus.svg)

#### Installation

- Download the binary from the hndreader [GitHub Release Page](github.com/FrontSide/hnreader/releases)

  For example: To download and install version 1.1 for linux, you may run the following.

  ```
  $ sudo curl -L https://github.com/FrontSide/hnreader/releases/download/v1.1/hnreader-v1.1.linux.x86_64 > /usr/local/bin/hnreader && chmod +x /usr/local/bin/hnreader
  ```

- **or** install the Go package

  ```
  $ go get -u github.com/Bunchhieng/hnreader
  ```

  Note that **this option requires** you to have **golang** already
  installed. You can install go with your operation system's package manager or download it from [golang.org/dl/](https://golang.org/dl/).

  Don't forget to set your GOPATH and PATH environment variables:

  ```
  $ export GOPATH=/path/to/your/go/workspace
  $ export PATH=$GOPATH/bin:$PATH
  ```

#### Usage

To use hnreader with its default options (Opens 10 news sites with chrome), simply run:

```
$ hnreader r
```

To see all available flags for each command:

```
$ hnreader help r
```

OR

```
$ hnreader help rr
```

There are a number of customization options:

```
--tabs value, -t value Specify value of tabs (default: 10)
--browser value, -b value Specify browser
--source value, -s value Specify news source (one of "hn", "reddit", "lobsters") (default: "hn")
```

Examples with options:

```
$ hnreader r -t 31 -b "firefox"
$ hnreader r -b "brave" -s "reddit"
$ hnreader r -b "firefox" -s "reddit" -t 20
```

To use hnreader with a randomized source of news, run:

```
$ hnreader rr
```

The following options are available:

```
--tabs value, -t value Specify value of tabs (default: 10)
--browser value, -b value Specify browser
```

**Tip:** Create a bash alias (for linux and macOS), if you are going to run the same command every morning.
You can do so by adding the following line (with your preferred options) to the end of your `~/.bashrc` file:

```
alias hnr='hnreader r -b "firefox" -s "reddit" -t 30' >> ~/.bashrc
```

#### Contribution

Please see the [CONTRIBUTING.md](CONTRIBUTING.md)

#### Credits

- [urfave/cli](https://github.com/urfave/cli)
- [Fatih Arslan](https://github.com/fatih/color)
- [Martin Angers](https://github.com/PuerkitoBio/goquery)
- [skratchdot](https://github.com/skratchdot/open-golang)

#### License

The MIT License (MIT)

```

```
