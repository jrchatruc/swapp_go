package main

import (
	"os"
	"os/exec"
	"reflect"
	"swapp_go/model"
	"swapp_go/utils"
	"testing"
)

func TestGetCharacter(t *testing.T) {
	c := make(chan model.People)
	go getCharacter("luke", c)

	if character := <-c; character.Name != "Luke Skywalker" {
		t.Errorf("getCharacter failed, expected %v, got %v", "Luke Skywalker", character.Name)
	}

	go getCharacter("leia", c)

	if character := <-c; character.Name != "Leia Organa" {
		t.Errorf("getCharacter failed, expected %v, got %v", "Leia Organa", character.Name)
	}
}

func TestGetFilm(t *testing.T) {
	if film := getFilm("https://swapi.co/api/films/1/"); film.Title != "A New Hope" {
		t.Errorf("getFilm failed, expected %v, got %v", "A New Hope", film.Title)
	}
}

func TestIntersect(t *testing.T) {
	a := []string{"a", "ab", "abc"}
	b := []string{"a", "abc", "abcd"}
	c := []string{"a", "abc"}

	if inter := utils.Intersect(a, b); !reflect.DeepEqual(inter, c) {
		t.Errorf("intersect failed, expected %v, got %v", c, inter)
	}

	a = []string{"ab", "abc", "a"}
	b = []string{"a", "ab", "abc"}
	c = []string{"ab", "abc", "a"}

	if inter := utils.Intersect(a, b); !reflect.DeepEqual(inter, c) {
		t.Errorf("intersect failed, expected %v, got %v", c, inter)
	}
}

func TestGetFilmCrash(t *testing.T) {
	if os.Getenv("BE_CRASH") == "1" {
		getFilm("")
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestGetFilmCrash")
	cmd.Env = append(os.Environ(), "BE_CRASH=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("TestGetFilmCrash ran with err %v, want exit status 1", err)
}
