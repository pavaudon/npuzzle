package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Conf struct {
	UI            *bool
	File          *string
	GenSize       *uint
	Solvable      *bool
	Iteration     *uint
	StartAfterGen *bool
	SortMode      *int
	Heuristic     *uint
	AlgoTimer     *bool
	GlobalTimer   *bool
}

type Resolv struct {
	env            *MyEnv
	Open           []Astar
	IdOpen         map[string]bool
	Closed         []Astar
	IdClosed       map[string]bool
	StartH         int
	SizeComplexity int
}

type MyEnv struct {
	Map            [][]uint8
	Found          bool
	MapSorted      [][]uint8
	Size           int
	SubthreadCount int
	Finish         chan struct{}
	Thread         chan struct{}
	ThreadMutex    sync.Mutex
	Conn           *websocket.Conn
	WebSocketMutex sync.Mutex
	Conf           *Conf
}

type Solved struct {
	Cost           int64
	Heuristic      string
	ComplexityTime int
	ComplexitySize int
	Moves          []string
}

// -o npuzzle
// --ui=false --mode=1 --file=examples\ordered\3.1.txt
// -file "examples/ordered/4.1.txt"
func main() {
	start := time.Now()
	conf := &Conf{
		UI:            flag.Bool("ui", true, "put true if you want a website to show our resolution\n"),
		File:          flag.String("file", "", "enter a file to resolv\n"),
		GenSize:       flag.Uint("size", 4, "Enter a size of the map to generate a new one randomly\n"),
		Solvable:      flag.Bool("solvable", true, "Enter if you want a solvalbe map or not\n"),
		Iteration:     flag.Uint("iteration", 0, "Number of passes on the generate npuzzle\nif 0, randomly pick a number between (size * size) and (size ^ size * size)"),
		StartAfterGen: flag.Bool("sag", false, "Start After random map Generation ?\n"),
		SortMode:      flag.Int("mode", 0, "mode 0 : traditional (white in bot-right) / mode 1 : snail\n"),
		Heuristic: flag.Uint("heuristic", 7, "H 0 : linear conflict\nH 1 : Manhattan\n"+
			"H 2 : Misplaced Tiles (le plus performant sur map 3x3)\nH 3 : Diagonal conflict\n"+
			"H 4 : Corner Tiles (le plus performant sur les grosses map (>3) et pas que)\n"+
			"H 5 : Combi LC (1) + DC (3) + CT (4)\nH 6 : Combi + DC (3) + CT (4)\nH 7 : MULTITHREAD WITH ALL HEURISTIC\n"),
		AlgoTimer:   flag.Bool("a-timer", true, "Print time to resolv (time taken from the beginning to the end of the resolution)\n"),
		GlobalTimer: flag.Bool("g-timer", true, "Print time to resolv (time taken at the launch and at the end of the software)\n"),
	}
	env := &MyEnv{Conf: conf}
	flag.Parse()

	if len(*env.Conf.File) > 0 {
		err := env.ParseFile(*env.Conf.File)
		if err != nil {
			log.Println("Fail to parse File, exit. Error:", err)
			return
		}
		if !env.ValidMap() {
			log.Println("[ERROR IN FILE] Bad values")
			return
		}
	} else {
		err := env.GeneratePuzzle()
		if err != nil {
			log.Println("ERROR TO GENERATE A MAP:", err)
			return
		}
		if !*env.Conf.StartAfterGen && !*env.Conf.UI {
			return
		}
	}
	if *env.Conf.UI {
		http.HandleFunc("/event_ws", env.WsEvent)
		handler := http.FileServer(http.Dir("./html"))
		http.Handle("/", handler)
		log.Println("Web UI on ", "http://127.0.0.1:8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
	} else { //RESOLVE HERE
		env.StartTheAlgo()
	}
	if *env.Conf.GlobalTimer {
		fmt.Printf("TOTAL DURATION TO RESOLVE : %dms\n", time.Now().Sub(start).Milliseconds())
	}

}

/*
PrintTheMap impirme la MAP correctement
*/
func PrintTheMap(Map [][]uint8) {
	fmt.Printf("%+v\n", strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(fmt.Sprintf("%v", Map), "] [", "\n"), "[", ""), "]", ""))
}

func (env *MyEnv) PrintError(err error) {
	log.Println("[ERROR]", err)
	if *env.Conf.UI {
		err = env.SendEvent(Event{Command: "error", Value: err.Error()})
		if err != nil {
			log.Println("[ERROR] Can't Write to the websocket :", err)
		}
	}
}
