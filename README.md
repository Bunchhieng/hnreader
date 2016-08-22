Hacker News Reader
=========

Open multiple hacker news feed with your favorite browser using command line.

#### Why?
Tired of clicking each link and for the sake of learning Golang. I read hacker news every morning, each every one of the link in the front page.

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
      `hnreader run -t 7 -b "firefox"`

#### Credits
* [urfave/cli](https://github.com/urfave/cli)
* [Fatih Arslan](https://github.com/fatih/color)
* [Martin Angers](https://github.com/PuerkitoBio/goquery)
* [skratchdot](https://github.com/skratchdot/open-golang)

#### License

The MIT License (MIT)
