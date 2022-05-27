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

func getfacingMe(me PlayerState, enemy PlayerState) (result string) {
	if me.X==enemy.X {
		if me.Y>enemy.Y {
			if enemy.Direction=="S" {
				return "X"
			}
		} else {
			if enemy.Direction=="N" {
				return "X"
			}
		}
	} else if me.Y==enemy.Y {
		if me.X>enemy.X {
			if enemy.Direction=="E" {
				return "Y"
			}
		} else {
			if enemy.Direction=="W" {
				return "Y"
			}
		}
	}
	return "N"
}

func escape(input ArenaUpdate, me PlayerState, dir string) (res string){
	var commands = []string{"R", "L"}
	switch dir {
	case "X":
		switch me.Direction {
		case "N":
			if me.X==0 {
				return "R"
			} else if me.X==input.Arena.Dimensions[0] {
				return "L"
			}
			return commands[rand2.Intn(2)]
		case "S":
			if me.X==0 {
				return "L"
			} else if me.X==input.Arena.Dimensions[0] {
				return "R"
			}
			return commands[rand2.Intn(2)]
		}
		break
	case "Y":
		switch me.Direction {
		case "E":
			if me.Y==0 {
				return "R"
			} else if me.Y==input.Arena.Dimensions[1] {
				return "L"
			}
			return commands[rand2.Intn(2)]
		case "W":
			if me.Y==0 {
				return "L"
			} else if me.Y==input.Arena.Dimensions[1] {
				return "R"
			}
			return commands[rand2.Intn(2)]
		}
		break
	}
	return "F"
}

func play(input ArenaUpdate) (response string) {
	log.Printf("IN: %#v", input)

	dir := []string{"X", "Y"}

	var me =  input.Arena.State[input.Links.Self.Href]
	delete(input.Arena.State, input.Links.Self.Href)

	for _, v := range input.Arena.State {
		if getfacingMe(v, me)!="N" {
			return "T"
		} 
	}

	for _, v := range input.Arena.State { 
		var enemyFacingMe = getfacingMe(me, v)
		if enemyFacingMe!="N" {
			return escape(input, me, enemyFacingMe)
		} 
	}

	return escape(input, me, dir[rand2.Intn(2)])
}
