package main

import (
	"encoding/json"
	"fmt"
	"log"
	rand2 "math/rand"
	"net/http"
	"os"
)

func main() {
	port := "8080"
	if v := os.Getenv("PORT"); v != "" {
		port = v
	}
	http.HandleFunc("/", handler)

	log.Printf("starting server on port :%s", port)
	err := http.ListenAndServe(":"+port, nil)
	log.Fatalf("http listen error: %v", err)
}

func handler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		fmt.Fprint(w, "Let the battle begin!")
		return
	}

	var v ArenaUpdate
	defer req.Body.Close()
	d := json.NewDecoder(req.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&v); err != nil {
		log.Printf("WARN: failed to decode ArenaUpdate in response body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := play(v)
	fmt.Fprint(w, resp)
}

func getRandRL() (response string) {
	var commands = []string{"R", "L"}
	return commands[rand2.Intn(2)]
}

func getfacingMe(me PlayerState, enemy PlayerState) (result string) {
	if me.X == enemy.X {
		if me.Y > enemy.Y {
			if enemy.Direction == "S" {
				return "X"
			}
		} else {
			if enemy.Direction == "N" {
				return "X"
			}
		}
	} else if me.Y == enemy.Y {
		if me.X > enemy.X {
			if enemy.Direction == "E" {
				return "Y"
			}
		} else {
			if enemy.Direction == "W" {
				return "Y"
			}
		}
	}
	return "N"
}

func getDir(input ArenaUpdate, me PlayerState, defaultResponse string) (res string) {
	switch me.Direction {
	case "N":
		if me.Y == 0 {
			if me.X == 0 {
				return "R"
			} else if me.X == input.Arena.Dimensions[0]-1 {
				return "L"
			}
		}
		break
	case "S":
		if me.Y == input.Arena.Dimensions[1]-1 {
			if me.X == 0 {
				return "L"
			} else if me.X == input.Arena.Dimensions[0]-1 {
				return "R"
			}
		}
		break
	case "E":
		if me.X == input.Arena.Dimensions[1]-1 {
			if me.Y == 0 {
				return "R"
			} else if me.Y == input.Arena.Dimensions[1]-1 {
				return "L"
			}
		}
		break
	case "W":
		if me.X == 0 {
			if me.Y == 0 {
				return "L"
			} else if me.Y == input.Arena.Dimensions[1]-1 {
				return "R"
			}
		}
		break
	}
	return defaultResponse
}

func escape(input ArenaUpdate, me PlayerState, dir string) (res string) {
	switch dir {
	case "X":
		switch me.Direction {
		case "N":
		case "S":
			return getDir(input, me, getRandRL())
		}
		break
	case "Y":
		switch me.Direction {
		case "E":
		case "W":
			return getDir(input, me, getRandRL())
		}
		break
	}
	return getDir(input, me, "F")
}

func play(input ArenaUpdate) (response string) {
	var mapN = make(map[string]PlayerState)
	var mapE = make(map[string]PlayerState)
	var mapS = make(map[string]PlayerState)
	var mapW = make(map[string]PlayerState)

	var me = input.Arena.State[input.Links.Self.Href]
	delete(input.Arena.State, input.Links.Self.Href)

	if me.WasHit {
		for _, v := range input.Arena.State {
			var enemyFacingMe = getfacingMe(me, v)
			if enemyFacingMe != "N" {
				return escape(input, me, enemyFacingMe)
			}
		}
	}

	for k, v := range input.Arena.State {
		if getfacingMe(v, me) != "N" {
			return "T"
		}
		if v.X == me.X {
			if v.Y < me.Y {
				mapN[k] = v
			} else {
				mapS[k] = v
			}
		} else if v.Y == me.Y {
			if v.X < me.X {
				mapW[k] = v
			} else {
				mapE[k] = v
			}
		}
	}

	for _, v := range mapN {
		var enemyFacingMe = getfacingMe(me, v)
		if enemyFacingMe != "N" {
			return escape(input, me, enemyFacingMe)
		}
	}
	for _, v := range mapS {
		var enemyFacingMe = getfacingMe(me, v)
		if enemyFacingMe != "N" {
			return escape(input, me, enemyFacingMe)
		}
	}
	for _, v := range mapW {
		var enemyFacingMe = getfacingMe(me, v)
		if enemyFacingMe != "N" {
			return escape(input, me, enemyFacingMe)
		}
	}
	for _, v := range mapE {
		var enemyFacingMe = getfacingMe(me, v)
		if enemyFacingMe != "N" {
			return escape(input, me, enemyFacingMe)
		}
	}

	switch me.Direction {
	case "N":
		if len(mapW) > len(mapE) {
			return getDir(input, me, "L")
		} else {
			return getDir(input, me, "R")
		}
	case "E":
		if len(mapN) > len(mapS) {
			return getDir(input, me, "L")
		} else {
			return getDir(input, me, "R")
		}
	case "S":
		if len(mapE) > len(mapW) {
			return getDir(input, me, "L")
		} else {
			return getDir(input, me, "R")
		}
	case "W":
		if len(mapS) > len(mapN) {
			return getDir(input, me, "L")
		} else {
			return getDir(input, me, "R")
		}
	}

	return getDir(input, me, getRandRL())
}
