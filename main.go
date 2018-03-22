package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"
)

const url = "https://vintagemonster.onefootball.com/api/teams/en/%v.json"

// Player is the struct to store necessary player information
// for Unmarshaled JSON
type Player struct {
	Country string `json:"country"`
	Name    string `json:"name"`
	Age     string `json:"age"`
}

// Players is used for storing results of total players
type Players []Player

// SingleTeam store result of each team in every API request
type SingleTeam struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Team struct {
			Name    string  `json:"name"`
			Players Players `json:"players"`
		} `json:"team"`
	} `json:"data"`
}

func isTeamDesired(name string) bool {
	var target = []string{"Germany", "England", "France", "Spain", "Manchester Utd", "Arsenal", "Chelsea", "Barcelona", "Real Madrid", "FC Bayern Munich"}
	for _, v := range target {
		if name == v {
			return true
		}
	}
	return false
}

func asyncHTTPGet(id int, jobs <-chan int, results chan<- Players) {
	for j := range jobs {
		var httpClient = &http.Client{
			Timeout: time.Second * 10,
		}
		resp, err := httpClient.Get(fmt.Sprintf(url, j))
		if err != nil {
			fmt.Printf("Team %d GET request error: %v\n", j, err.Error())
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Team %d read response body err: %v\n", j, err.Error())
		}
		var jsonResult = SingleTeam{}
		if err = json.Unmarshal(body, &jsonResult); err != nil {
			fmt.Printf("Team %d Unmarshl JSON response body err: %v\n", j, err.Error())
		}
		if jsonResult.Code != 0 && jsonResult.Code == 404 {
			if strings.Contains(jsonResult.Message, "not find team with id") {
				fmt.Printf("Team %d Not Found\n", j)
				results <- Players{}
			}
		}
		if isTeamDesired(jsonResult.Data.Team.Name) {
			results <- jsonResult.Data.Team.Players
		}
	}
}

func main() {

	jobs := make(chan int, 3)
	ch := make(chan Players, 10)

	var results Players
	var nonDuplicate = make(map[string]Player)

	for w := 1; w <= 10; w++ {
		go asyncHTTPGet(w, jobs, ch)
	}

	teamID := 1
	for {
		jobs <- teamID

		// Break for loop if result channel is full
		if len(ch) == cap(ch) {
			break
		}
		teamID++
	}
	close(jobs)

	for i := 1; i <= 10; i++ {
		res := <-ch
		for _, player := range res {
			// If the name of each player is not key in nonDuplicate
			// insert this player into nonDuplicate
			if _, ok := nonDuplicate[player.Name]; !ok {
				nonDuplicate[player.Name] = player
			}
		}
	}

	// Append non-duplicate players to final result slice
	for _, player := range nonDuplicate {
		results = append(results, player)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Name < results[j].Name
	})

	for _, value := range results {
		fmt.Printf("%s; %s; %s\n", value.Name, value.Age, value.Country)
	}
}
