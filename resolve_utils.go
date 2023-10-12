package main

import (
	"errors"
	"strconv"
)

/*
	IsClosed verifie si un channel est ouvert ou non, pour ne pas le fermer 2 fois de suite.
*/
func IsClosed(ch <-chan struct{}) bool {
	if ch == nil {
		return true
	}
	select {
	case <-ch:
		return true
	default:
	}
	return false
}

func (res *Resolv) AlreadyInOpen(ast *Astar) bool {
	if _, ok := res.IdOpen[ast.Id]; ok {
		return true
	}
	return false
}

/*
	AlreadyInClosed Verifie si Closed contient deja la meme carte.
*/
func (res *Resolv) AlreadyInClosed(ast *Astar) bool {
	if _, ok := res.IdClosed[ast.Id]; ok {
		return true
	}
	return false
}

/*
	GetLCForSnail convert a snail map to a standard map to check if map is solvable, return the heuristic value like a traditionnal map.
*/
func (env *MyEnv) GetLCForSnail(mymap [][]uint8) (int, error) {
	copyMap := make([][]uint8, env.Size)
	for i := range env.Map {
		copyMap[i] = make([]uint8, env.Size)
	}
	for i := range copyMap {
		for ii := range copyMap[i] {
			copyMap[i][ii] = mymap[i][ii]
		}
	}
	//Reorder the map to check if solvable
	//create a tmpMap with a new order
	x0, y0, err2 := findPlaceOnMap(env.MapSorted, 0)
	if err2 != nil {
		return 0, err2
	}
	tmpMap := make([][]uint8, env.Size)
	for i := range tmpMap {
		tmpMap[i] = make([]uint8, env.Size)
	}
	count := uint8(1)
	i, ii := 0, 0
	for int(count) < (env.Size * env.Size) {
		if !(i == x0 && ii == y0) {
			tmpMap[i][ii] = count
			count++
		}
		ii++
		if ii == env.Size {
			ii = 0
			i++
		}
	}
	//assign the new value to map
	for i := range tmpMap {
		for ii := range tmpMap[i] {
			x, y, err2 := findPlaceOnMap(copyMap, env.MapSorted[i][ii])
			if err2 != nil {
				return 0, err2
			}

			copyMap[x][y] = -tmpMap[i][ii]
		}
	}
	for i := range copyMap {
		for ii := range env.Map[i] {
			copyMap[i][ii] = -copyMap[i][ii]
		}
	}
	return env.TradiInversion(copyMap), nil
}

/*
	isSolvable Verifie si la map est correcte.
*/
func (env *MyEnv) isSolvable() (err error) {
	if env.Size < 3 || env.Size > 7 {
		return errors.New("Map is too small or too big, it needs to be between 3 and 6")
	}
	LinearConflict := 0
	X0, _, err := findPlaceOnMap(env.Map, 0)
	if err != nil {
		return
	}
	switch *env.Conf.SortMode {
	case 0:
		LinearConflict = env.TradiInversion(env.Map)
		X0++ //si le 0 est sur une ligne pair
	case 1:
		LinearConflict, err = env.GetLCForSnail(env.Map)
		if err != nil {
			return
		}
		X02, _, err2 := findPlaceOnMap(env.MapSorted, 0)
		if err2 != nil {
			return err2
		}
		if X02%2 != 0 {
			X0++
		}
	default:
		return errors.New("Can't solve this sort mode")
	}

	if env.Size%2 != 0 && LinearConflict%2 == 0 {
		return
	}
	if env.Size%2 == 0 &&
		((LinearConflict+X0)%2 == 0) {
		return
	}
	return errors.New("Map is not solvable")
}

/*
	makeAMove Deplace le 0 a la destination, et augmente le cout de 1
*/
func (ast *Astar) makeAMove(x, y, zeroX, zeroY int) (err error) {
	if (zeroX-x == 0 && zeroY-y != -1 && zeroY-y != 1) || // Same colomn
		(zeroY-y == 0 && zeroX-x != -1 && zeroX-x != 1) || // Same line
		(zeroY-y != 0 && zeroX-x != 0) { // Other Move
		return errors.New("The move is too big")
	}
	ast.MyMoveX = x
	ast.MyMoveY = y
	ast.Cost += 1
	val := ast.CurrentMap[x][y]
	ast.CurrentMap[x][y] = 0
	ast.CurrentMap[zeroX][zeroY] = val
	return
}

func (ast *Astar) setId() {
	for i := range ast.CurrentMap {
		for ii := range ast.CurrentMap[i] {
			ast.Id += strconv.Itoa(int(ast.CurrentMap[i][ii])) + "/"
		}
	}
	return
}

/*
	getHeuristicName retourne le nom de l'heuristique en fonction de sa valeur.
*/
func getHeuristicName(H uint) (ret string) {
	switch H {
	case 0:
		ret = "linear conflict"
	case 1:
		ret = "Manhattan"
	case 2:
		ret = "Misplaced Tiles"
	case 3: 
		ret = "Diagonal conflict"
	case 4:
		ret = "Corner Tiles"
	case 5:
		ret = "Combi LC (1) + DC (3) + CT (4)"
	case 6:
		ret = "Combi DC (3) + CT (4)"
	}
	return
}
