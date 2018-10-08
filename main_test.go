package main

import (
	"log"
	"reflect"
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

func TestGetDZoneStories(t *testing.T) {
	news, err := new(DZoneSource).Fetch(10)
	if err != nil {
		log.Fatal(err)
	}

	assert.NotNil(t, news)
	assert.Equal(t, 10, len(news), "They should be equal")
}

func TestGetDevToStories(t *testing.T) {
	news, err := new(DevToSource).Fetch(10)
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}

	assert.NotNil(t, news)
	assert.Equal(t, 10, len(news), "They should be equal")
}

func TestGetBrowserNameByOS(t *testing.T) {
	assertErrMsg := "They should be equal"

	os := "darwin"
	assert.Equal(t, "Firefox", getBrowserNameByOS("firefox", os), assertErrMsg)
	assert.Equal(t, "Firefox", getBrowserNameByOS("mozilla", os), assertErrMsg)
	assert.Equal(t, "Google Chrome", getBrowserNameByOS("chrome", os), assertErrMsg)
	assert.Equal(t, "Google Chrome", getBrowserNameByOS("google", os), assertErrMsg)
	assert.Equal(t, "Brave", getBrowserNameByOS("brave", os), assertErrMsg)

	os = "linux"
	assert.Equal(t, "firefox", getBrowserNameByOS("firefox", os), assertErrMsg)
	assert.Equal(t, "firefox", getBrowserNameByOS("mozilla", os), assertErrMsg)
	assert.Equal(t, "google-chrome", getBrowserNameByOS("chrome", os), assertErrMsg)
	assert.Equal(t, "google-chrome", getBrowserNameByOS("google", os), assertErrMsg)
	assert.Equal(t, "brave", getBrowserNameByOS("brave", os), assertErrMsg)

	os = "windows"
	assert.Equal(t, "firefox", getBrowserNameByOS("firefox", os), assertErrMsg)
	assert.Equal(t, "firefox", getBrowserNameByOS("mozilla", os), assertErrMsg)
	assert.Equal(t, "chrome", getBrowserNameByOS("chrome", os), assertErrMsg)
	assert.Equal(t, "chrome", getBrowserNameByOS("google", os), assertErrMsg)
	assert.Equal(t, "brave", getBrowserNameByOS("brave", os), assertErrMsg)
}
