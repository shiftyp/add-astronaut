package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"
)

type reset func()

type astronaut struct {
	Id int64 `json:"id"`
	Color string `json:"color"`
	Power string `json:"power"`
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func main() {
	maxCount := 50
	mutex := &sync.Mutex{}
	astronauts := []astronaut{}
	colorPattern := regexp.MustCompile(`^#[A-Fa-f0-9]{6}$`)
	powerPattern := regexp.MustCompile(`^(\x{1F4A5}|\x{1F496}|\x{1F4A7}|\x{1F525}|\x{2B50}|\x{1F48E})$`)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			var count int
			var err error
			var values url.Values

			values, err = url.ParseQuery(r.URL.RawQuery)
			
			if err == nil && values.Get("count") != "" {
				count, err = strconv.Atoi(values.Get("count"));
			}

			if count <= 0 || err != nil {
				count = maxCount;
			}

			astronautsJson, err := json.Marshal(astronauts[0:min(len(astronauts), count)]);
			
			if err != nil {
				log.Fatal(err)
			}

			fmt.Fprintf(w, string(astronautsJson))
		case "POST":
			decoder := json.NewDecoder(r.Body)

			var astro astronaut

			err := decoder.Decode(&astro)

			if err != nil {
				log.Print(err)
				http.Error(w, "Error", 500)
			}

			match := colorPattern.MatchString(astro.Color) && powerPattern.MatchString(astro.Power)

			if match == false {
				log.Printf("Bad Astronaut %+v", astro)
				http.Error(w, "Bad Astronaut", 400)
			} else {
				mutex.Lock()
				astro.Id = time.Now().UnixNano()
				astronauts = append([]astronaut{astro}, astronauts[0:min(len(astronauts), maxCount - 1)]...)
				mutex.Unlock()
			}
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	fmt.Printf("Astronauts launching on :%s 🚀\n", port)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)

	if err != nil {
		panic(err)
	}
}
