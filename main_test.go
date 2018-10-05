package main

import (
	"log"
	"reflect"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	actual := Init()
	expected := &App{Name: AppName, Version: AppVersion, Description: AppDescription, Author: AppAuthor, Email: AppEmail}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Test failed, expected: '%s', got:  '%s'", expected, actual)
	}
}

func TestGetHNStories(t *testing.T) {
	news, err := new(HackerNewsSource).Fetch(10)
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}

	assert.NotNil(t, news)
	assert.Equal(t, 30, len(news), "They should be equal")
}

func TestGetRedditStories(t *testing.T) {
	news, err := new(RedditSource).Fetch(10)
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}

	assert.NotNil(t, news)
	assert.Equal(t, 10, len(news), "They should be equal")
}

func TestGetLobstersStories(t *testing.T) {
	news, err := new(LobstersSource).Fetch(10)
	if err != nil {
		log.Fatal(err)
	}

	assert.NotNil(t, news)
	assert.Equal(t, 10, len(news), "They should be equal")
}

func TestGetBrowserNameByOS(t *testing.T) {
	// For now, getBrowserName only assumes the logic below
	// applies to the darwin OS; add subtests as more os values
	// are added.
	if runtime.GOOS == "darwin" {
		t.Run("Validate Darwin OS", func(t *testing.T) {
			browserNameMap := map[string][]string{
				"Google Chrome": []string{"google", "chrome"},
				"Firefox":       []string{"mozilla", "firefox"},
				"Brave":         []string{"brave"},
			}

			// test every possible browser string value
			for normalizedStr, browserStrSlice := range browserNameMap {
				for _, browserStr := range browserStrSlice {
					assert.Equal(t, getBrowserNameByOS(browserStr), normalizedStr)
				}
			}
		})
	}
}
