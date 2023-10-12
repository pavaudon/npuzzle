package main

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"time"
)

/*
	PrintMoves imprime les mouvements
*/
func (curr *Astar) PrintMoves() {
	if curr.MyDad != nil {
		curr.MyDad.PrintMoves()
	} else {
		return
	}
	log.Println("Move ", curr.Cost, " | 0 == Map[", curr.MyMoveX, "][", curr.MyMoveY, "]")
	return
}

/*
	GetMoves retourne les mouvements pour le site
*/
func (curr *Astar) GetMoves() (ret []string) {
	if curr.MyDad != nil {
		ret = curr.MyDad.GetMoves()
	} else {
		return
	}
	ret = append(ret, fmt.Sprint(curr.MyMoveX, "x", curr.MyMoveY))
	return
}

/*
	AddMove ajoute un enfant dans open seulement s'il n'est pas deja dans open ou closed
*/
func (res *Resolv) AddMove(X, Y, ZeroX, ZeroY int, curr *Astar) (err error) {
	MapMove, _ := res.newAstar(curr.CurrentMap, false, curr.SelectedHeuristic)

	err = MapMove.makeAMove(X, Y, ZeroX, ZeroY)
	if err != nil {
		return err
	}
	MapMove.setId()
	err = MapMove.selectHeuristic(res.env)
	if err != nil {
		return err
	}

	if !res.AlreadyInClosed(MapMove) && !res.AlreadyInOpen(MapMove) {
		MapMove.Cost = curr.Cost + 1
		MapMove.MyMoveX, MapMove.MyMoveY = X, Y
		MapMove.MyDad = curr
		_ = res.AddElem(*MapMove, true)
		if len(res.Open) > res.SizeComplexity {
			res.SizeComplexity = len(res.Open)
		}
	}
	return nil
}

/*
	NewMoves me dit quels sont les possibles enfants de curr
*/
func (res *Resolv) NewMoves(curr *Astar) (err error) {
	ZeroX, ZeroY, err := findPlaceOnMap(curr.CurrentMap, 0)
	if err != nil {
		return
	}
	if ZeroX != 0 {
		err = res.AddMove(ZeroX-1, ZeroY, ZeroX, ZeroY, curr)
		if err != nil {
			return err
		}
	}
	if ZeroX != res.env.Size-1 {
		err = res.AddMove(ZeroX+1, ZeroY, ZeroX, ZeroY, curr)
		if err != nil {
			return err
		}
	}
	if ZeroY != 0 {
		err = res.AddMove(ZeroX, ZeroY-1, ZeroX, ZeroY, curr)
		if err != nil {
			return err
		}
	}
	if ZeroY != res.env.Size-1 {
		err = res.AddMove(ZeroX, ZeroY+1, ZeroX, ZeroY, curr)
		if err != nil {
			return err
		}
	}
	return nil
}

/*
	search prend a chaque fois une instance de open pour chercher ses enfants et voir si la map est resolue
*/
func (res *Resolv) search() (err error) {
	i := 0
	for len(res.Open) > 0 && !IsClosed(res.env.Thread) {
		elem := res.Open[0]
		if reflect.DeepEqual(elem.CurrentMap, res.env.MapSorted) {
			res.env.ThreadMutex.Lock()
			defer res.env.ThreadMutex.Unlock()
			if IsClosed(res.env.Thread) {
				return
			}
			close(res.env.Thread)
			res.env.Thread = nil
			if *res.env.Conf.UI {

				solved := &Solved{
					Cost:           elem.Cost,
					Heuristic:      fmt.Sprint(getHeuristicName(elem.SelectedHeuristic), "(", elem.SelectedHeuristic, ")"),
					ComplexityTime: len(res.Closed) + 1,
					ComplexitySize: res.SizeComplexity,
					Moves:          elem.GetMoves(),
				}
				solvedJSON, err := json.Marshal(solved)
				if err != nil {
					log.Printf("ERROR TO MARSHAL THE CONF: %v\n", err)
				}
				res.env.SendEvent(Event{Command: "solved", Value: string(solvedJSON)})
			} else {
				log.Println("SUCCESS WITH", elem.Cost, "moves")
				log.Println("SUCCESS WITH THE HEURISTIC:", getHeuristicName(elem.SelectedHeuristic), "(", elem.SelectedHeuristic, ")")
				log.Println("The complexity in time is", len(res.Closed)+1)
				log.Println("The complexity in size is", res.SizeComplexity)
				elem.PrintMoves()
				PrintTheMap(elem.CurrentMap)
			}
			log.Println("SUCCESS WITH", elem.Cost, "moves")
			close(res.env.Finish)
			res.env.Finish = nil
			return nil
		}

		i++
		res.AddElem(elem, false)
		res.DelFirstElem(true)
		err = res.NewMoves(&elem)
		if err != nil {
			return err
		}
		if len(res.Open) > 1 {
			err = res.sortOpenAstar()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

/*
	Astar est le premier maillon de l'algo avant de boucler dans search() (crée le first avant de lancer la boucle)
*/
func (env *MyEnv) Astar(H uint) {
	defer func() {
		if !IsClosed(env.Thread) {
			close(env.Thread)
			env.Thread = nil
		}
	}()
	res := &Resolv{env: env, IdOpen: map[string]bool{}, IdClosed: map[string]bool{}}
	first, err := res.newAstar(res.env.Map, true, H)
	if err != nil {
		log.Println("[ERROR]", err)
		return
	}
	err = res.NewMoves(first)
	if err != nil {
		log.Println("[ERROR]", err)
		return
	}
	if len(res.Open) > 1 {
		err = res.sortOpenAstar()
		if err != nil {
			log.Println("[ERROR]", err)
			return
		}
	}
	res.SizeComplexity = 1
	err = res.search()
	if err != nil {
		log.Println("[ERROR]", err)
		return
	}
	return
}

/*
	StartTheAlgo permet de faire les premieres verifications, et de lancer plusieurs threads si multithread
*/
func (env *MyEnv) StartTheAlgo() {

	if err := env.isSolvable(); err != nil {
		env.PrintError(err)
		return
	}
	start := time.Now()
	if !reflect.DeepEqual(env.Map, env.MapSorted) {
		env.Finish = make(chan struct{})
		env.Thread = make(chan struct{})
		if *env.Conf.Heuristic < 7 {
			go env.Astar(*env.Conf.Heuristic)
		} else {
			for i := 0; i < 6; i++ {
				log.Println("Start thread with : ", getHeuristicName(uint(i)))
				go env.Astar(uint(i))
			}
		}
		<-env.Finish
	} else if *env.Conf.UI {
		solved := &Solved{}
		solvedJSON, err := json.Marshal(solved)
		if err != nil {
			log.Printf("ERROR TO MARSHAL THE CONF: %v\n", err)
		}
		env.SendEvent(Event{Command: "solved", Value: string(solvedJSON)})
	} else {
		log.Println("### ALREADY SOLVED ###")
		log.Println("SUCCESS WITH", 0, "moves")
		log.Println("SUCCESS WITH THE HEURISTIC:", "NO HEURISTIC")
		log.Println("The complexity in time is", 0)
		log.Println("The complexity in size is", 0)
		PrintTheMap(env.Map)
	}
	if *env.Conf.AlgoTimer {
		fmt.Printf("Algo DURATION TO RESOLVE : %dms\n", time.Now().Sub(start).Milliseconds())
		if *env.Conf.UI {
			env.SendEvent(Event{Command: "timer", Value: fmt.Sprint(time.Now().Sub(start).Milliseconds(), " ms")})
		}
	}
	return
}
