package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"swapp_go/model"
	"swapp_go/utils"
)

var (
	BASE_URL = "http://swapi.co/api/"
	PEOPLE   = "people/"
	SEARCH   = "?search="
)

func main() {

	//Esta función está para timear el tiempo de ejecución.
	//defer utils.TimeTrack(time.Now(), "main")

	if len(os.Args) < 3 {
		fmt.Println("Error, ingrese dos personajes.")
		os.Exit(0)
	}

	arg1 := os.Args[1]
	arg2 := os.Args[2]
	c := make(chan model.People, 2)

	go getCharacter(arg1, c)
	go getCharacter(arg2, c)
	character1 := <-c
	character2 := <-c

	filmsBoth := utils.Intersect(character1.FilmsURL, character2.FilmsURL)

	if len(filmsBoth) == 0 {
		fmt.Println("Los personajes ingresados no tienen ninguna película en común.")
		os.Exit(0)
	}

	URLS := make(chan string, len(filmsBoth))
	results := make(chan model.Film, len(filmsBoth))

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

func getCharacter(name string, c chan model.People) {
	resp, err := http.Get(BASE_URL + PEOPLE + SEARCH + name)
	if err != nil {
		fmt.Println("Something went wrong.")
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var ret model.PeopleContainer
	json.Unmarshal(body, &ret)

	if ret.Count != 1 {
		fmt.Println("Alguno de los personajes ingresados no es válido.")
		os.Exit(0)
	}
	c <- ret.Results[0]
}

func getFilm(URL string) model.Film {
	resp, err := http.Get(URL)
	if err != nil {
		fmt.Println("Something went wrong.")
		os.Exit(1)
	}
	defer resp.Body.Close()
	var film model.Film
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &film)

	return film
}

func filmFetcher(URLS <-chan string, results chan<- model.Film) {
	for url := range URLS {
		results <- getFilm(url)
	}
}
