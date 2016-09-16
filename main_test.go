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

func TestGetStories(t *testing.T) {
	news := make(map[int]string)
	news, err := GetStories(10)
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}

	assert.NotNil(t, news)
	assert.Equal(t, 30, len(news), "They should be equal")
}
