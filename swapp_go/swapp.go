package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	BASE_URL = "http://swapi.co/api/"
	PEOPLE   = "people/"
	SEARCH   = "?search="
)

func main() {

	//Esta función está para timear el tiempo de ejecución.
	defer timeTrack(time.Now(), "main")

	if len(os.Args) < 3 {
		fmt.Println("Error, ingrese dos personajes.")
		os.Exit(0)
	}

	arg1 := os.Args[1]
	arg2 := os.Args[2]
	c := make(chan People, 2)

	go getCharacter(arg1, c)
	go getCharacter(arg2, c)
	character1 := <-c
	character2 := <-c

	filmsBoth := intersect(character1.FilmsURL, character2.FilmsURL)

	if len(filmsBoth) == 0 {
		fmt.Println("Los personajes ingresados no tienen ninguna película en común.")
		os.Exit(0)
	}

	URLS := make(chan string, len(filmsBoth))
	results := make(chan Film, len(filmsBoth))

	for _, _ = range filmsBoth {
		go filmFetcher(URLS, results)
	}

	for _, url := range filmsBoth {
		URLS <- url
	}
	close(URLS)

	for _, _ = range filmsBoth {
		fmt.Println(<-results)
	}

}

func getCharacter(name string, c chan People) {
	resp, err := http.Get(BASE_URL + PEOPLE + SEARCH + name)
	if err != nil {
		fmt.Println("Something went wrong.")
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var ret PeopleQuerySet
	json.Unmarshal(body, &ret)

	if ret.Count != 1 {
		fmt.Println("Alguno de los personajes ingresados no es válido.")
		os.Exit(0)
	}
	c <- ret.Results[0]
}

func getFilms(URL string) Film {
	resp, err := http.Get(URL)
	if err != nil {
		fmt.Println("Something went wrong.")
		os.Exit(1)
	}
	defer resp.Body.Close()
	var film Film
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &film)

	return film
}

func filmFetcher(URLS <-chan string, results chan<- Film) {
	for url := range URLS {
		results <- getFilms(url)
	}
}

//Utils
func intersect(list1 []string, list2 []string) []string {
	ret := []string{}

	for _, s := range list1 {
		if contains(list2, s) {
			ret = append(ret, s)
		}
	}
	return ret
}

func contains(a []string, s string) bool {
	for _, i := range a {
		if i == s {
			return true
		}
	}
	return false
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

//End Utils

//Model
type People struct {
	Name     string   `json:"name"`
	FilmsURL []string `json:"films"`
}

func (p People) String() string {
	return p.Name
}

type PeopleQuerySet struct {
	Count   int      `json:"count"`
	Results []People `json:"results"`
}

type Film struct {
	Title string `json:"title"`
}

func (f Film) String() string {
	return f.Title
}

//End Model
