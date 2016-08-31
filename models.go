package main

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v1"
)

type Event struct {
	Event struct {
		Tag   string
		Date  string
		Venue string
	}
	Sponsors []string
	Title    string
}

type Talks struct {
	Talks []struct {
		Title     string
		Presenter struct {
			Name    string
			Tagline string
			Twitter string
		}
		Abstract string
	}
}

type Sponsors map[string]struct {
	Image string
	Name  string
	URL   string
}

var event Event
var talks Talks
var sponsors Sponsors

func loadTalks(filename string) error {
	cnt, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(cnt, &talks)
	if err != nil {
		return err
	}

	return nil
}

func loadEvent(filename string) error {
	cnt, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(cnt, &event)
	if err != nil {
		return err
	}

	return nil
}

func loadSponsors(filename string) error {
	cnt, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(cnt, &sponsors)
	if err != nil {
		return err
	}

	return nil
}
