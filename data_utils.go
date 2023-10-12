package main

import (
	"errors"
	"sort"
)

type Astar struct {
	CurrentMap        [][]uint8
	Id                string
	Cost              int64
	MapHeuristic      int
	SelectedHeuristic uint
	MyMoveX           int
	MyMoveY           int
	MyDad             *Astar
}

/* newAstar duplicate a map and return a new instance*/
func (res *Resolv) newAstar(Map [][]uint8, first bool, H uint) (ast *Astar, err error) {
	ast = &Astar{}
	ast.CurrentMap = make([][]uint8, res.env.Size)
	for i := range ast.CurrentMap {
		ast.CurrentMap[i] = make([]uint8, res.env.Size)
		copy(ast.CurrentMap[i], Map[i])
	}
	ast.SelectedHeuristic = H
	return
}

/* AddElem Ajoute un element a la fin de CLOSED ou De OPEN*/
func (res *Resolv) AddElem(elem Astar, isOpen bool) (ret int) {
	if isOpen {
		res.Open = append(res.Open, elem)
		res.IdOpen[elem.Id] = true
		return len(res.Open)
	} else {
		res.Closed = append(res.Closed, elem)
		res.IdClosed[elem.Id] = true
		return len(res.Closed)
	}
}

/* DelFirstElem Supprime le premier element de CLOSED ou De OPEN*/
func (res *Resolv) DelFirstElem(isOpen bool) {
	if isOpen {
		delete(res.IdOpen, res.Open[0].Id)
		res.Open = append(res.Open[:0], res.Open[1:]...)
	} else {
		delete(res.IdClosed, res.Closed[0].Id)
		res.Closed = append(res.Closed[:0], res.Closed[1:]...)
	}
}

/* sortOpenAstar Effectue le tri de Open pour avoir le plus petit heuristique.*/
func (res *Resolv) sortOpenAstar() (err error) {
	if len(res.Open) == 0 {
		return errors.New("OPEN IS EMPTY")
	}
	sort.Slice(res.Open, func(i, j int) bool {
		return res.Open[i].MapHeuristic < res.Open[j].MapHeuristic
	})
	return
}
