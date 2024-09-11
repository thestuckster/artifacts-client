package utils

import "math"

type MapCoOrds struct {
	X, Y int
}

// FindClosestMapTile uses Manhattan distance [md = (|x1 - x2|) + (|y1 - y2|)]  to determine the closest map tile to your
// character. If you pass in an empty array of target coords it will return your current position. If it cannot find a
// closest tile, for whatever reason, it will return you to the origin (0,0)
func FindClosestMapTile(originX, originY int, targetCoOrds []MapCoOrds) (destX, destY int) {

	if len(targetCoOrds) == 0 {
		return originX, originY
	}

	shortest := 1000
	for _, coords := range targetCoOrds {
		temp := int((math.Abs(float64(originX) - float64(coords.X))) + (math.Abs(float64(originY) - float64(coords.Y))))
		if temp < shortest {
			shortest = temp
			destX = coords.X
			destY = coords.Y
		}
	}

	return destX, destY
}
