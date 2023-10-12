package main

import (
	"errors"
	"math"
)

/*
	findPlaceOnMap trouve ou est la valeur donnée en parametre
*/
func findPlaceOnMap(Map [][]uint8, val uint8) (x, y int, err error) {
	for x = range Map {
		for y = range Map[x] {
			if Map[x][y] == val {
				return
			}
		}
	}
	return 0, 0, errors.New("CAN'T FIND VALUE ON MAP")
}

/*
	selectHeuristic retourne la valeur H en fonction de l'heuristique
*/
func (myStar *Astar) selectHeuristic(env *MyEnv) (err error) {
	var h int

	switch myStar.SelectedHeuristic {
	case 0: //linear conflict
		h, err = myStar.LinearConflict(env)
		if err != nil {
			return
		}
		h2, err2 := myStar.Manhattan(env)
		if err2 != nil {
			return
		}
		h += h2
		break
	case 1: //Manhattan
		h, err = myStar.Manhattan(env)
		if err != nil {
			return
		}
		break
	case 2: //Misplaced Tiles
		h, err = myStar.MisplacedTile(env)
		if err != nil {
			return
		}
		h2, err2 := myStar.Manhattan(env)
		if err2 != nil {
			return
		}
		h += h2
		break
	case 3: //Diagonal conflict
		h, err = myStar.Manhattan(env)
		if err != nil {
			return
		}
		h2, err2 := myStar.DiagonalConflict(env)
		if err2 != nil {
			return
		}
		h += h2
		break
	case 4: //Corner Tiles
		h, err = myStar.Manhattan(env)
		if err != nil {
			return
		}
		h += myStar.CornerTiles(env)
		break
	case 5: //Combi LC (1) + DC (3) + CT
		h, err = myStar.Manhattan(env)
		if err != nil {
			return
		}
		h2, err2 := myStar.LinearConflict(env)
		if err2 != nil {
			return
		}
		h3, err3 := myStar.DiagonalConflict(env)
		if err3 != nil {
			return
		}
		h += myStar.CornerTiles(env) + h2 + h3
		break
	case 6: //Combi DC (3) + CT
		h, err = myStar.Manhattan(env)
		if err != nil {
			return
		}
		h3, err3 := myStar.DiagonalConflict(env)
		if err3 != nil {
			return
		}
		h += myStar.CornerTiles(env) + h3
		break
	}
	myStar.MapHeuristic = h
	return
}

/*
	SnailInversion retourne la valeur H Pour l'heuristique LinearConflict sur une Map type serpent
*/
func (env *MyEnv) SnailInversion(mymap [][]uint8) (ret int) {
	oneCount, oneLimite := 1, 1
	i, ii := 0, 0
	for oneCount < (env.Size * env.Size) {
		one := mymap[i][ii]
		j, jj := i, ii
		twoCount, twoLimite := oneCount, oneLimite
		for twoCount < (env.Size * env.Size) {
			if j >= twoLimite && jj == twoLimite-1 {
				j--
				if j == twoLimite {
					twoLimite++
				}
			} else if j == env.Size-twoLimite && jj > twoLimite-1 {
				jj--
			} else if jj == env.Size-twoLimite {
				j++
			} else if jj < env.Size {
				jj++
			}
			two := mymap[j][jj]
			if one != 0 && two != 0 && one > two {
				ret++
			}
			twoCount++
		}
		if i >= oneLimite && ii == oneLimite-1 {
			i--
			if i == oneLimite {
				oneLimite++
			}
		} else if i == env.Size-oneLimite && ii > oneLimite-1 {
			ii--
		} else if ii == env.Size-oneLimite {
			i++
		} else if ii < env.Size {
			ii++
		}
		oneCount++
	}
	return ret
}

/*
	TradiInversion retourne la valeur H Pour l'heuristique LinearConflict sur une Map type traditionnelle
*/
func (env *MyEnv) TradiInversion(mymap [][]uint8) (ret int) {
	dMap := make([]byte, env.Size*env.Size)
	for _, v := range mymap {
		dMap = append(dMap, v...)
	}
	for i := 0; i < len(dMap); i++ {
		for ii := i + 1; ii < len(dMap); ii++ {
			if dMap[i] != 0 && dMap[ii] != 0 && dMap[i] > dMap[ii] {
				ret++
			}
		}
	}
	return ret
}

/*
 	LinearConflict Cet Heuristique, en fonction de comment est orienté la map,
	compte le nombre de case inferieur entre sa position et la fin du puzzle
*/
func (myStar *Astar) LinearConflict(env *MyEnv) (conflicts int, err error) {
	switch *env.Conf.SortMode {
	case 0:
		conflicts = env.TradiInversion(myStar.CurrentMap)
	case 1:
		conflicts = env.SnailInversion(myStar.CurrentMap)
	default:
		return 0, errors.New("CHANGE THE SORT MODE FAILED")
	}
	return conflicts, nil
}

/* MisplacedTile Cet Heuristique nous dit seulement combien de cases ne sont pas a leur place */
func (myStar *Astar) MisplacedTile(env *MyEnv) (dist int, err error) {
	for i := range myStar.CurrentMap {
		for ii := range myStar.CurrentMap[i] {
			if myStar.CurrentMap[i][ii] != 0 && myStar.CurrentMap[i][ii] != env.MapSorted[i][ii] {
				dist++
			}
		}
	}
	return dist, nil
}

/* Manhattan Cet Heuristique calcule les distances pour chaque piece, et les aditionne */
func (myStar *Astar) Manhattan(env *MyEnv) (dist int, err error) {
	for i := range myStar.CurrentMap {
		for ii := range myStar.CurrentMap[i] {
			if myStar.CurrentMap[i][ii] != 0 {
				srcX, srcY, err1 := findPlaceOnMap(myStar.CurrentMap, myStar.CurrentMap[i][ii])
				if err1 != nil {
					return
				}
				dstX, dstY, err2 := findPlaceOnMap(env.MapSorted, myStar.CurrentMap[i][ii])
				if err2 != nil {
					return
				}
				dist += int(math.Abs(float64(dstX)-float64(srcX)) + math.Abs(float64(dstY)-float64(srcY)))
			}
		}
	}
	return
}

/*
	DiagonalConflict Cet Heuristique verifie la map pour voir si les destinations sont en diagonale
*/
func (myStar *Astar) DiagonalConflict(env *MyEnv) (dist int, err error) {
	for i := range myStar.CurrentMap {
		for ii := range env.Map[i] {
			if myStar.CurrentMap[i][ii] != 0 {
				srcX, srcY, err1 := findPlaceOnMap(myStar.CurrentMap, myStar.CurrentMap[i][ii])
				if err1 != nil {
					return
				}
				dstX, dstY, err2 := findPlaceOnMap(env.MapSorted, myStar.CurrentMap[i][ii])
				if err2 != nil {
					return
				}
				if myStar.CurrentMap[dstX][dstY] == 0 || (srcX == dstX && srcY == dstY) {
					continue
				}
				if i-1 == dstX && ii-1 == dstY || // Coin supperieur gauche
					i-1 == dstX && ii+1 == dstY || //Coin supperieur droit
					i+1 == dstX && ii-1 == dstY || //Coin inferieur gauche
					i+1 == dstX && ii+1 == dstY { // Coin inferieur droit
					dist += 2
				}
			}
		}
	}
	return
}

/*
	DiagonalConflict Cet Heuristique Ajoute des coups si les coins sont pas placé mais que les cases a coté le sont.
*/
func (myStar *Astar) CornerTiles(env *MyEnv) (dist int) {
	var a, b, c, d bool
	if myStar.CurrentMap[0][0] != 0 && myStar.CurrentMap[0][0] != env.MapSorted[0][0] { //Coin supperieur gauche
		if myStar.CurrentMap[0][1] == env.MapSorted[0][1] {
			dist += 2
			a = true
		}
		if myStar.CurrentMap[1][0] == env.MapSorted[1][0] {
			dist += 2
			d = true
		}
	}
	if myStar.CurrentMap[0][env.Size-1] != 0 && myStar.CurrentMap[0][env.Size-1] != env.MapSorted[0][env.Size-1] { //Coin supperieur droit
		if myStar.CurrentMap[0][env.Size-2] == env.MapSorted[0][env.Size-2] && (!a || env.Size > 3) { 
			dist += 2
		}
		if myStar.CurrentMap[1][env.Size-1] == env.MapSorted[1][env.Size-1] {
			dist += 2
			b = true
		}
	}
	if myStar.CurrentMap[env.Size-1][0] != 0 && myStar.CurrentMap[env.Size-1][0] != env.MapSorted[env.Size-1][0] { //Coin inferieur gauche
		if myStar.CurrentMap[env.Size-2][0] == env.MapSorted[env.Size-2][0] && (!d || env.Size > 3) { 
			dist += 2
		}
		if myStar.CurrentMap[env.Size-1][1] == env.MapSorted[env.Size-1][1] {
			dist += 2
			c = true
		}
	}
	if myStar.CurrentMap[env.Size-1][env.Size-1] != 0 && myStar.CurrentMap[env.Size-1][env.Size-1] != env.MapSorted[env.Size-1][env.Size-1] { // Coin inferieur droit
		if myStar.CurrentMap[env.Size-2][env.Size-1] == env.MapSorted[env.Size-2][env.Size-1] && (!b || env.Size > 3) { 
			dist += 2
		}
		if myStar.CurrentMap[env.Size-1][env.Size-2] == env.MapSorted[env.Size-1][env.Size-2] && (!c || env.Size > 3) {
			dist += 2
		}
	}
	return
}
