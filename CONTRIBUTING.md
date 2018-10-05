# Contributing to hnreader

Contributing to hnreader isn't limited to just filing bugs, users are more than welcomed to make suggestions, report any issue they may find, and make pull requests to help make hnreader better.

## Working on hnreader

### Prerequisites

- [Git](https://git-scm.com/)
- [Go](https://golang.org/dl/)
- [Godep](https://github.com/golang/dep)
- Your favorite text editor or IDE. (i.e: Atoms, Visual Studio Code or Intellij)

### Getting hnreader

1. Fork the repo
2. Clone the repo with the follow command

```
$ git clone https://github.com/YOUR-GITHUB-PROFILE-NAME/hnreader.git
$ cd hnreader
```

3. Install dependencies if needed: `dep ensure`

### Please pay attention to

1. Open an issue describing the feature/bug you wish to contribute first to start a discussion, explain why, what and how
2. Follow the syntax rules
3. Write tests covering 100% of the library code you produce
4. One PR per feature/fix unless you follow [standard-version](https://github.com/conventional-changelog/standard-version) commit guidelines

### Using branches

When working on any issue on Github, it's a good practice to make branches that are specific to the issue you're currently working on. For instance, if you're working on an issue with a name like "NAME OF ISSUE #1234", from the master branch run the following code: `git checkout -b Issue#1234`. In doing so, you'll be making a branch that specifically identifies the issue at hand, and moves you right into it with the `checkout` flag. This keeps your main (master) repository clean and your personal workflow cruft out of sight when making a pull request.

### Finding issues to fix

After you've forked and cloned the repo, you can find issues to work on by heading over to our [issues list](https://github.com/Bunchhieng/hnreader/issues). We advise looking at the issues with the labels [help wanted](https://github.com/Bunchhieng/hnreader/issues?q=is%3Aissue+is%3Aopen+label%3A%22help+wanted%22) or [good first issue](https://github.com/Bunchhieng/hnreader/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22), as they will help you get familiar with the hnreader code.

### Rules of the discussions

Remember to be very clear and transparent when discussing any issue in the discussions boards. We ask that you keep the language to English and keep on track with the issue at hand. Lastly, please be respectful of our fellow contributors and keep an exemplary level of professionalism at all times.
