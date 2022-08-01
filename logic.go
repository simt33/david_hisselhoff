package main

// This file can be a nice home for your Battlesnake logic and related helper functions.
//
// We have started this for you, with a function to help remove the 'neck' direction
// from the list of possible moves!

import (
	"log"
  "fmt"
  
)

var hazardCoords []Coord

// This function is called when you register your Battlesnake on play.battlesnake.com
// See https://docs.battlesnake.com/guides/getting-started#step-4-register-your-battlesnake
// It controls your Battlesnake appearance and author permissions.
// For customization options, see https://docs.battlesnake.com/references/personalization
// TIP: If you open your Battlesnake URL in browser you should see this data.
func info() BattlesnakeInfoResponse {
	log.Println("INFO")
	return BattlesnakeInfoResponse{
		APIVersion: "1",
		Author:     "",        // TODO: Your Battlesnake username
		Color:      "#CDFF03", // TODO: Personalize
		Head:       "shades", // TODO: Personalize
		Tail:       "curled", // TODO: Personalize
	}
}

// This function is called everytime your Battlesnake is entered into a game.
// The provided GameState contains information about the game that's about to be played.
// It's purely for informational purposes, you don't have to make any decisions here.
func start(state GameState) {
  log.Printf("START")
	if len(state.Board.Hazards) > 0 {
    hazardCoords = BuildHazardCoords(state)
  }
}

// This function is called when a game your Battlesnake was in has ended.
// It's purely for informational purposes, you don't have to make any decisions here.
func end(state GameState) {
	//log.Printf("%s END\n\n", state.Game.ID)
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func BuildPossibleMoves (coord Coord, state GameState) map[string]Coord {
  
  possibleMoves := map[string]Coord{
		"up":    Coord{X: coord.X, Y: coord.Y + 1},
		"down":  Coord{X: coord.X, Y: coord.Y - 1},
		"left":  Coord{X: coord.X - 1, Y: coord.Y},
		"right": Coord{X: coord.X + 1, Y: coord.Y},
	}
  
  return possibleMoves
  
}

func MakeChecks (move string, coord Coord, state GameState, occupiedCoords []Coord) bool {

return CheckPath(move, coord, state, 0, occupiedCoords)
}


func CheckPath (move string, coord Coord, state GameState, depth int, occupiedCoords []Coord) bool {
  //A recursive function that scans ahead 4 spaces to check for obstacles
  fmt.Printf("Scanning ahead (%s): %i space(s)", move, depth)
  
  if depth == 5 {
    return true
  }

  possibleMoves := BuildPossibleMoves (coord, state)
 
  for move, testCoord := range possibleMoves {
    
    safeSpace := CheckSafety(move, testCoord, state, occupiedCoords)
    
    if safeSpace == true {
      newOccupiedCoords := append(occupiedCoords, testCoord)
      pathFound := CheckPath(move, testCoord, state, depth + 1, newOccupiedCoords)
      if pathFound == true {
        return true
      }
    }
  }
  return false
}


func BuildOccupiedCoords (state GameState) []Coord {

  var occupiedCoords []Coord 
  
  dangerNoodles := state.Board.Snakes

  for i := 0; i< len(dangerNoodles); i++ {
    dangerNoodleLength := int(dangerNoodles[i].Length)
    dangerNoodleBody := dangerNoodles[i].Body

    // Every other snek's body
    for j := 0; j< dangerNoodleLength; j++ {
      occupiedCoords = append(occupiedCoords, dangerNoodleBody[j])
    }
	  
    // All moves around every other snek's head
    if dangerNoodles[i].ID != state.You.ID {
    	snekNextMoves := BuildPossibleMoves (dangerNoodles[i].Head, state)
    
    	for _, coord := range snekNextMoves {
      	   occupiedCoords = append(occupiedCoords, coord)
    }
    
    // ... And their actual head
    occupiedCoords = append(occupiedCoords, dangerNoodles[i].Head)
  } 
  }
  return occupiedCoords
}

func BuildHazardCoords (state GameState) []Coord {
  
  hazards := state.Board.Hazards

  for i := 0; i< len(hazards); i++ {
    hazardCoords = append(hazardCoords, hazards[i])
  }

  log.Printf("Built Hazard Coords")
  return hazardCoords
}

func CheckSafety (move string, testCoord Coord, state GameState, occupiedCoords []Coord) bool {
  //Checks safety of the new coordinate
  fmt.Printf("Checking safety: %s\n", move)
  
  x := testCoord.X
  y := testCoord.Y  
  
  //Check for walls
  if x < 0 || x >= state.Board.Width || y < 0 || y >= state.Board.Height {
    //fmt.Printf("OUT OF BOUNDS\n")
    return false
  }
  
  //Check for body and other sneks  
  for i := 1; i < len(occupiedCoords) ; i++ {
    //fmt.Printf("comparing to occupied Coord %i: %i\n", occupiedCoords[i].X, occupiedCoords[i].Y)
    
    if occupiedCoords[i].X == x && occupiedCoords[i].Y == y {
      fmt.Printf("SNEK FOUND %s\n", move)
      return false
    }
  }
  
  return true
}


func HungrySnek(preferredMoves map[string]Coord, state GameState, occupiedCoords []Coord) string {
  //Returns string value of possibleMove which minimises distance to closest food.

  if len(preferredMoves) <= 0 {
    return ""
  }
  
  var availableMoves []string
    
  for move, _ := range preferredMoves {
    availableMoves = append(availableMoves, move)
  }
  log.Printf("available moves: %v", availableMoves)
  
  allFood := state.Board.Food
  bestMove := ""
  var bestMoveCoord Coord
  
  if len(allFood) > 0 {
    
    myHead := state.You.Body[0]
    log.Printf("NOM NOM NOM", move)
    closestFoodX := 0
    closestFoodY := 0
    closestFoodDistance := 1000
  
    //Get Closest Food  
    for i := 0; i< len(allFood); i++ {
        diffX := Abs(myHead.X - allFood[i].X)
        diffY := Abs(myHead.Y - allFood[i].Y)
        distance := diffX + diffY
      
        if distance <= closestFoodDistance {
          closestFoodDistance = distance
          closestFoodX = allFood[i].X
          closestFoodY = allFood[i].Y
        }
    
    }
    //Get best move towards the food
    bestMoveDistance := 1000
    
    for move, coord := range preferredMoves{
      diffX := Abs(coord.X - closestFoodX)
      diffY := Abs(coord.Y - closestFoodY)
      distance := diffX + diffY
  
      if distance < bestMoveDistance{
        bestMove = move
        bestMoveDistance = distance
        bestMoveCoord = coord
      }
      
    }
  } else {

    log.Printf("NO FOOD!!")
    
    bestMove = availableMoves[0]
    bestMoveCoord = preferredMoves[bestMove]
  }
  
  //Recursively check path for this plan
  pathAhead := MakeChecks(bestMove, bestMoveCoord, state, occupiedCoords)
  if pathAhead == false{
    delete(preferredMoves, bestMove)
    bestMove = HungrySnek(preferredMoves, state, occupiedCoords)
  }
  
  return bestMove
  
}

// This function is called on every turn of a game. Use the provided GameState to decide
// where to move -- valid moves are "up", "down", "left", or "right".
// We've provided some code and comments to get you started.
func move(state GameState) BattlesnakeMoveResponse {
	
  possibleMoves := BuildPossibleMoves(state.You.Body[0], state)  
  occupiedCoords := BuildOccupiedCoords(state)
  
  if len(state.Board.Hazards) > 0 {
    for i := 0; i<len(hazardCoords); i++ {
        occupiedCoords = append(occupiedCoords, hazardCoords[i])
      }
  }
  
  for move, testCoord := range possibleMoves {
    valid := CheckSafety(move, testCoord, state, occupiedCoords)
    if valid == false {
      delete(possibleMoves, move)
    } else {
      log.Printf("%s is valid", move)
    }
  }
  
	var nextMove string

	if len(possibleMoves) == 0 {
		nextMove = "down"
		log.Printf("%s MOVE %d: No safe moves detected! Moving %s\n", state.Game.ID, state.Turn, nextMove)
	} else {
    
  //FIND Preferred Move which minimises distance to food and return nextMove as string (key of the possible move)
		
    preferredMoves := make(map[string]Coord)

    for k,v := range possibleMoves {
    preferredMoves[k] = v
    }

    occupiedCoords = append(occupiedCoords, state.You.Head)
    
    nextMove = HungrySnek(preferredMoves, state, occupiedCoords)
		log.Printf("MOVE %d: %s\n", state.Turn, nextMove)

    if nextMove == "" {

      var availableMoves []string
      
      for move, _ := range possibleMoves {
    
        availableMoves = append(availableMoves, move)
      }
      nextMove = availableMoves[0]
    }
	}
	return BattlesnakeMoveResponse{
		Move: nextMove,
	}
}
