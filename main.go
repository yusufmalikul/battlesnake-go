package main

// Welcome to
// __________         __    __  .__                               __
// \______   \_____ _/  |__/  |_|  |   ____   ______ ____ _____  |  | __ ____
//  |    |  _/\__  \\   __\   __\  | _/ __ \ /  ___//    \\__  \ |  |/ // __ \
//  |    |   \ / __ \|  |  |  | |  |_\  ___/ \___ \|   |  \/ __ \|    <\  ___/
//  |________/(______/__|  |__| |____/\_____>______>___|__(______/__|__\\_____>
//
// This file can be a nice home for your Battlesnake logic and helper functions.
//
// To get you started we've included code to prevent your Battlesnake from moving backwards.
// For more info see docs.battlesnake.com

import (
	"fmt"
	"log"
)

// info is called when you create your Battlesnake on play.battlesnake.com
// and controls your Battlesnake's appearance
// TIP: If you open your Battlesnake URL in a browser you should see this data
func info() BattlesnakeInfoResponse {
	log.Println("INFO")

	return BattlesnakeInfoResponse{
		APIVersion: "1",
		Author:     "yusufmalikul", // Your Battlesnake username
		Color:      "#0040ff",      // Choose color
		Head:       "earmuffs",     // Choose head
		Tail:       "bolt",         // Choose tail
	}
}

// start is called when your Battlesnake begins a game
func start(state GameState) {
	log.Println("GAME START")
}

// end is called when your Battlesnake finishes a game
func end(state GameState) {
	log.Printf("GAME OVER\n\n")
}

// move is called on every turn and returns your next move
// Valid moves are "up", "down", "left", or "right"
// See https://docs.battlesnake.com/api/example-move for available data
func move(state GameState) BattlesnakeMoveResponse {

	isMoveSafe := map[string]bool{
		"up":    true,
		"down":  true,
		"left":  true,
		"right": true,
	}

	// We've included code to prevent your Battlesnake from moving backwards
	myHead := state.You.Body[0] // Coordinates of your head
	myNeck := state.You.Body[1] // Coordinates of your "neck"

	if myNeck.X < myHead.X { // Neck is left of head, don't move left
		isMoveSafe["left"] = false

	} else if myNeck.X > myHead.X { // Neck is right of head, don't move right
		isMoveSafe["right"] = false

	} else if myNeck.Y < myHead.Y { // Neck is below head, don't move down
		isMoveSafe["down"] = false

	} else if myNeck.Y > myHead.Y { // Neck is above head, don't move up
		isMoveSafe["up"] = false
	}

	// Step 1 - Prevent your Battlesnake from moving out of bounds
	boardWidth := state.Board.Width - 1
	boardHeight := state.Board.Height - 1

	if isMoveSafe["up"] && myHead.Y == boardHeight {
		isMoveSafe["up"] = false
	}

	if isMoveSafe["down"] && myHead.Y == 0 {
		isMoveSafe["down"] = false
	}

	if isMoveSafe["right"] && myHead.X == boardWidth {
		isMoveSafe["right"] = false
	}

	if isMoveSafe["left"] && myHead.X == 0 {
		isMoveSafe["left"] = false
	}

	log.Println("myHead.Y", myHead.Y, "boardHeight", boardHeight)
	log.Println("myHead.X", myHead.X, "boardWidth", boardWidth)

	// Step 2 - Prevent your Battlesnake from colliding with itself
	// log.Println("start debug body")
	mybody := state.You.Body
	for _, v := range mybody {
		// log.Println("body X:", v.X, "body Y:", v.Y)
		if isMoveSafe["up"] && myHead.Y+1 == v.Y && myHead.X == v.X {
			isMoveSafe["up"] = false
		}

		if isMoveSafe["down"] && myHead.Y-1 == v.Y && myHead.X == v.X {
			isMoveSafe["down"] = false
		}

		if isMoveSafe["right"] && myHead.X+1 == v.X && myHead.Y == v.Y {
			isMoveSafe["right"] = false
		}

		if isMoveSafe["left"] && myHead.X-1 == v.X && myHead.Y == v.Y {
			isMoveSafe["left"] = false
		}
	}
	// log.Println("end debug body")

	// Step 3 - Prevent your Battlesnake from colliding with other Battlesnakes
	opponents := state.Board.Snakes
	for _, v := range opponents {
		for _, snake := range v.Body {
			// log.Println("body X:", snake.X, "body Y:", snake.Y)
			if isMoveSafe["up"] && myHead.Y+1 == snake.Y && myHead.X == snake.X {
				isMoveSafe["up"] = false
			}

			if isMoveSafe["down"] && myHead.Y-1 == snake.Y && myHead.X == snake.X {
				isMoveSafe["down"] = false
			}

			if isMoveSafe["right"] && myHead.X+1 == snake.X && myHead.Y == snake.Y {
				isMoveSafe["right"] = false
			}

			if isMoveSafe["left"] && myHead.X-1 == snake.X && myHead.Y == snake.Y {
				isMoveSafe["left"] = false
			}
		}
	}

	// Prevent entering closed grid with no way to escape
	// we will use flood fill algorithm
	// First we mark all filled grid (filled by our snake body and other snake body)
	moveSafeGridSize := map[string]int{
		"up":    0,
		"down":  0,
		"left":  0,
		"right": 0,
	}

	if isMoveSafe["up"] {
		total := fill(prepareGrid(boardWidth+1, boardHeight+1, opponents), boardWidth, boardHeight, myHead.X, myHead.Y+1, "", "@")
		moveSafeGridSize["up"] = total
	}

	if isMoveSafe["down"] {
		total := fill(prepareGrid(boardWidth+1, boardHeight+1, opponents), boardWidth, boardHeight, myHead.X, myHead.Y-1, "", "@")
		moveSafeGridSize["down"] = total
	}

	if isMoveSafe["right"] {
		total := fill(prepareGrid(boardWidth+1, boardHeight+1, opponents), boardWidth, boardHeight, myHead.X+1, myHead.Y, "", "@")
		moveSafeGridSize["right"] = total
	}

	if isMoveSafe["left"] {
		total := fill(prepareGrid(boardWidth+1, boardHeight+1, opponents), boardWidth, boardHeight, myHead.X-1, myHead.Y, "", "@")
		moveSafeGridSize["left"] = total
	}

	fmt.Println("isMoveSafe", isMoveSafe)
	fmt.Println("moveSafeGridSize", moveSafeGridSize)

	// Are there any safe moves left?
	safeMoves := []string{}
	for move, isSafe := range isMoveSafe {
		if isSafe {
			safeMoves = append(safeMoves, move)
		}
	}
	log.Println("safeMoves:", safeMoves)

	if len(safeMoves) == 0 {
		log.Printf("MOVE %d: No safe moves detected! Moving down\n", state.Turn)
		return BattlesnakeMoveResponse{Move: "down"}
	}

	// Choose a random move from the safe ones
	// nextMove := safeMoves[rand.Intn(len(safeMoves))]

	// Choose best move with largest room
	nextMove := "up"
	largestRoom := moveSafeGridSize["up"]
	for i := range moveSafeGridSize {
		if moveSafeGridSize[i] > largestRoom {
			largestRoom = moveSafeGridSize[i]
			nextMove = i
		}
	}

	// TODO: Step 4 - Move towards food instead of random, to regain health and survive longer
	// food := state.Board.Food

	log.Printf("MOVE %d: %s\n", state.Turn, nextMove)
	return BattlesnakeMoveResponse{Move: nextMove}
}

func fill(grid [][]string, boardWidth, boardHeight, x int, y int, old string, new string) int {

	total := 0

	if x < 0 {
		return 0
	}

	if y < 0 {
		return 0
	}

	if x > boardWidth {
		return 0
	}

	if y > boardHeight {
		return 0
	}

	if grid[x][y] == new || grid[x][y] != old {
		return 0
	}

	grid[x][y] = "X"
	total++

	total += fill(grid, boardWidth, boardHeight, x+1, y, old, new)
	total += fill(grid, boardWidth, boardHeight, x-1, y, old, new)
	total += fill(grid, boardWidth, boardHeight, x, y+1, old, new)
	total += fill(grid, boardWidth, boardHeight, x, y-1, old, new)

	return total
}

func prepareGrid(width int, height int, opponents []Battlesnake) [][]string {
	markedGrid := make([][]string, width)
	for i := range markedGrid {
		markedGrid[i] = make([]string, height)
	}
	for _, v := range opponents {
		// fmt.Println("debug grid")
		for _, snake := range v.Body {
			// fmt.Println("X", snake.X, "Y", snake.Y)
			markedGrid[snake.X][snake.Y] = "X"
		}
	}
	return markedGrid
}

func main() {
	RunServer()
}
