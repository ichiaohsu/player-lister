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

type Player struct {
	Country string `json:"country"`
	Name    string `json:"name"`
	Age     string `json:"age"`
}

type Players []Player

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

func isTeamDesired(name string, target []string) bool {
	for _, v := range target {
		if name == v {
			return true
		}
	}
	return false
}

func main() {

	var desiredTeams = []string{"Germany", "England", "France", "Spain", "Manchester Utd", "Arsenal", "Chelsea", "Barcelona", "Real Madrid", "FC Bayern Munich"}
	var httpClient = &http.Client{
		Timeout: time.Second * 10,
	}
	var result Players
	teamID := 1
	teamCount := 0

	for {
		response, err := httpClient.Get(fmt.Sprintf(url, teamID))
		if err != nil {
			panic(err)
		}
		defer response.Body.Close()

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			panic(err)
		}
		var jsonResult = SingleTeam{}
		if err = json.Unmarshal(body, &jsonResult); err != nil {
			panic(err)
		}
		if jsonResult.Code != 0 && jsonResult.Code == 404 {
			if strings.Contains(jsonResult.Message, "not find team with id") {
				break
			}
			panic(err)
		}
		if isTeamDesired(jsonResult.Data.Team.Name, desiredTeams) {
			result = append(result, jsonResult.Data.Team.Players...)
			teamCount++
		}
		if teamCount == len(desiredTeams) {
			break
		}
		teamID++
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	for _, value := range result {
		fmt.Printf("%s; %s; %s\n", value.Name, value.Age, value.Country)
	}
}
