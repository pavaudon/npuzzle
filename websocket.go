package main

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Event struct {
	Command string `json:"command"` //get_conf + set_conf + map + generate_map + solved + start + import_map + ok + error + timer
	Value   string `json:"value"`
}

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
	return true
}}

/*
	WsEvent gere les websocket pour l'interface web
*/
func (env *MyEnv) WsEvent(w http.ResponseWriter, r *http.Request) {
	var err error
	env.Conn, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer env.Conn.Close()
	mapS := MapToString(env.Map)
	env.SendEvent(Event{Command: "map", Value: mapS})
	for {
		_, message, err := env.Conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		evt := Event{}
		err = json.Unmarshal(message, &evt)
		if err != nil {
			log.Printf("ERROR TO DECRYPT MESSAGE FROM WEBSOCKET: %v, message %+v\n", err, string(message))
		}
		switch evt.Command {
		case "generate_map":
			err = env.GeneratePuzzle()
			if err != nil {
				env.PrintError(err)
				continue
			}
			mapS := MapToString(env.Map)
			env.SendEvent(Event{Command: "map", Value: mapS})
			break
		case "import_map":
			err := env.InsertTheMap(evt.Value)
			if err != nil {
				log.Printf("ERROR TO IMPORT: %v\n", err)
				env.PrintError(err)
				continue
			}
			mapS := MapToString(env.Map)
			env.SendEvent(Event{Command: "map", Value: mapS})
			break
		case "get_conf":
			confJSON, err := json.Marshal(env.Conf)
			if err != nil {
				log.Printf("ERROR TO MARSHAL THE CONF: %v\n", err)
				env.PrintError(err)
				continue
			}
			env.SendEvent(Event{Command: "set_conf", Value: string(confJSON)})
			break
		case "set_conf":
			if strings.Contains(evt.Value, "null") {
				env.PrintError(errors.New("WRONG CONF, can't be set"))
				continue
			}
			err := json.Unmarshal([]byte(evt.Value), env.Conf)
			if err != nil {
				log.Printf("ERROR TO MARSHAL THE CONF: %v\n", err)
				env.PrintError(err)
				continue
			}
			break
		case "start":
			go env.StartTheAlgo()
			break
		default:
			break
		}
		err = env.SendEvent(Event{Command: "ok"})
		if err != nil {
			log.Println("SendEvent ERROR:", err)
			break
		}
	}
}

/*
	SendEvent envoi un event aux sites connectés
*/
func (env *MyEnv) SendEvent(EL Event) error {
	env.WebSocketMutex.Lock()
	defer env.WebSocketMutex.Unlock()
	if env.Conn == nil {
		return errors.New("No websocket conn")
	}
	return env.Conn.WriteJSON(EL)
}

/*
	MapToString transforme le tableau d'int en une string JSON
*/
func MapToString(Map [][]uint8) (ret string) {
	type mapString struct {
		Map [][]string `json:"map"`
	}
	MS := mapString{}
	size := len(Map)
	MS.Map = make([][]string, size)
	for i := range MS.Map {
		MS.Map[i] = make([]string, size)
	}
	for i := range Map {
		for ii := range Map {
			MS.Map[i][ii] = strconv.FormatUint(uint64(Map[i][ii]), 10)
		}
	}
	retB, err := json.Marshal(MS)
	if err != nil {
		return ""
	}
	return string(retB)
}
