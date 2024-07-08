package main

import (
	"fmt"
	"math"
	"math/rand"
	"syscall/js"
)

type TanksGameState struct {
	Player       int
	WindX        int
	Turn         int
	BulletPaths  [][]Coordinates
	Hits         []Hit
	PlayersState map[int]PlayerState
}

type PlayerState struct {
	X           int
	Y           int
	CannonAngle int
	HP          int
}

type TanksDefaultSettings struct {
	Width int
	Y     int
	MaxHP int
}

type Coordinates struct {
	X float64
	Y float64
}

type Hit struct {
	Turn      int
	HitPlayer int
}

// type EventsValues int

var EVENTS = map[string]string{
	"UP":    "ArrowUp",
	"DOWN":  "ArrowDown",
	"LEFT":  "ArrowLeft",
	"RIGHT": "ArrowRight",
	"FIRE":  " ",
}

var DefaultGameState = TanksGameState{
	Player:      0,
	WindX:       120,
	Turn:        0,
	BulletPaths: [][]Coordinates{},
	Hits:        []Hit{},
}

// const getDefaultGameState = (
//
//	settings: TanksDefaultSettings
//
//	): TanksGameState => {
//	  const state = {
//	    ...DEFAULT_GAME_STATE,
//	    playersState: {
//	      0: getDefaultPlayerState(0, settings),
//	      1: getDefaultPlayerState(1, settings),
//	    },
//	  }
//	  return state as TanksGameState
//	}

// settings TanksDefaultSettings
func GetDefaultGameState(this js.Value, settings js.Value) TanksGameState {
	state := DefaultGameState
	state.PlayersState = map[int]PlayerState{
		0: getDefaultPlayerState(0, settings),
		1: getDefaultPlayerState(1, settings),
	}
	return state
}

// const getDefaultPlayerState = (
//
//	index: number,
//	settings: TanksDefaultSettings
//
//	) => {
//	  return {
//	    X: !index ? 100 : settings.WIDTH - 150,
//	    Y: settings.Y,
//	    cannonAngle: 0,
//	    HP: settings.MAX_HP,
//	  }
//	}
//
// func getDefaultPlayerState(index int, settings TanksDefaultSettings) PlayerState {
func getDefaultPlayerState(index int, settings js.Value) interface{} {
	x := 100
	width := settings.Get("WIDTH").Int()
	y := settings.Get("Y").Int()
	max_hp := settings.Get("MAX_HP").Int()
	if index != 0 {
		x = width - 150
	}
	return PlayerState{
		X:           x,
		Y:           y,
		CannonAngle: 0,
		HP:          max_hp,
	}
}

// const handlePlayersStateValueChange = (
//
//	prevState: TanksGameState,
//	key: keyof PlayerState,
//	multiplier = 1,
//	step = 5
//
//	) => {
//	  const currentValue = prevState.playersState[prevState.player][key]
//	  const nextValue = currentValue + step * multiplier
//	  const isNextValueOutOfCannonAngle =
//	    key === "cannonAngle" && (Math.abs(nextValue) > 90 || nextValue > 0)
//	  if (isNextValueOutOfCannonAngle) {
//	    return prevState
//	  }
//	  return {
//	    playersState: {
//	      ...prevState.playersState,
//	      [prevState.player]: {
//	        ...prevState.playersState[prevState.player],
//	        [key]: nextValue,
//	      },
//	    } as TanksGameState["playersState"],
//	  }
//	}
func handlePlayersStateValueChange(prevState TanksGameState, key string, multiplier, step int) TanksGameState {
	if multiplier == 0 {
		multiplier = 1
	}
	if step == 0 {
		step = 5
	}

	currentValue := prevState.PlayersState[prevState.Player]
	var nextValue int
	switch key {
	case "X":
		nextValue = currentValue.X + step*multiplier
	case "cannonAngle":
		nextValue = currentValue.CannonAngle + step*multiplier
	}

	isNextValueOutOfCannonAngle := key == "cannonAngle" && (math.Abs(float64(nextValue)) > 90 || nextValue > 0)
	if isNextValueOutOfCannonAngle {
		return prevState
	}

	newPlayersState := prevState.PlayersState
	if key == "X" {
		newPlayersState[prevState.Player] = PlayerState{X: nextValue, Y: currentValue.Y, CannonAngle: currentValue.CannonAngle, HP: currentValue.HP}
	} else if key == "cannonAngle" {
		newPlayersState[prevState.Player] = PlayerState{X: currentValue.X, Y: currentValue.Y, CannonAngle: nextValue, HP: currentValue.HP}
	}

	prevState.PlayersState = newPlayersState
	return prevState
}

//	const handleWindChange = () => {
//	  return generateRandomIntInRange(200, -200)
//	}
func handleWindChange() int {
	return rand.Intn(400) - 200
}

//	const getDegreeToRadians = (degree: number) => {
//	  const CONSTANT = 0.0174532925
//	  return CONSTANT * Math.abs(degree)
//	}
func getDegreeToRadians(degree float64) float64 {
	constant := 0.0174532925
	return constant * math.Abs(degree)
}

//	const getXInT = (t: number, alfa: number, v0: number) => {
//	  // v0t cos(alfa)
//	  return v0 * t * Math.cos(getDegreeToRadians(alfa))
//	}
func getXInT(t, alfa, v0 float64) float64 {
	return v0 * t * math.Cos(getDegreeToRadians(alfa))
}

//	const getYInT = (t: number, alfa: number, v0: number) => {
//	  // v0t sin(alfa) - g(t^2)/2
//	  const G = 9.8 // gravity
//	  return v0 * t * Math.sin(getDegreeToRadians(alfa)) - (G * Math.pow(t, 2)) / 2
//	}
func getYInT(t, alfa, v0 float64) float64 {
	const g = 9.8
	return v0*t*math.Sin(getDegreeToRadians(alfa)) - (g*math.Pow(t, 2))/2
}

// const getBulletPath = (wind: number, angle: number): Coordinates[] => {
//   // console.log("getBulletPath", { angle })
//   const BULLET_MASS = 10 // kg
//   const INTERVAL = 0.01 // s
//   const V0 = 100
//   const FLOOR_LVL = -100

//   let t = 0 + INTERVAL
//   const result = [
//     {
//       X: 0,
//       Y: 0,
//     },
//   ]
//   while (true) {
//     const nextResult = {
//       X: getXInT(t, angle, V0),
//       Y: getYInT(t, angle, V0),
//     }
//     result.push(nextResult)
//     t += INTERVAL

//	    if (nextResult.Y <= FLOOR_LVL) {
//	      // console.log({ result })
//	      return result
//	    }
//	  }
//	}
func getBulletPath(wind, angle int) []Coordinates {
	const (
		BULLET_MASS = 10
		INTERVAL    = 0.01
		V0          = 100
		FLOOR_LVL   = -100
	)

	var t float64 = INTERVAL
	result := []Coordinates{{X: 0, Y: 0}}

	for {
		nextResult := Coordinates{
			X: getXInT(t, float64(angle), float64(V0)),
			Y: getYInT(t, float64(angle), float64(V0)),
		}
		result = append(result, nextResult)
		t += INTERVAL

		if nextResult.Y <= FLOOR_LVL {
			return result
		}
	}
}

// const getHits = (
//   gameState: TanksGameState,
//   nextPath: Coordinates[]
// ): Hit | null => {
//   const TANK_WIDTH = 150
//   const TANK_HEIGHT = 100
//   const PLAYERS_COUNT = 2
//   const currentPlayerPosition = {
//     X: gameState.playersState[gameState.player].X + TANK_WIDTH,
//     Y: gameState.playersState[gameState.player].Y,
//   }

//   const hitBoundriesP0X = [
//     gameState.playersState[0].X - TANK_WIDTH,
//     gameState.playersState[0].X + TANK_WIDTH,
//   ]
//   const hitBoundriesP1X = [
//     gameState.playersState[1].X - TANK_WIDTH,
//     gameState.playersState[1].X + TANK_WIDTH,
//   ]
//   const hitBoundriesP0Y = [
//     gameState.playersState[0].Y,
//     gameState.playersState[0].Y + TANK_HEIGHT,
//   ]
//   const hitBoundriesP1Y = [
//     gameState.playersState[1].Y,
//     gameState.playersState[1].Y + TANK_HEIGHT,
//   ]

//   const pXs = [hitBoundriesP0X, hitBoundriesP1X]
//   const pYs = [hitBoundriesP0Y, hitBoundriesP1Y]

//   for (let index = 0; index < nextPath.length; index++) {
//     const xP0 = nextPath[index].X + currentPlayerPosition.X
//     const xP1 =
//       -1 * nextPath[index].X + currentPlayerPosition.X - TANK_WIDTH * 2
//     const y = nextPath[index].Y * -1
//     let xHitPlayer = null
//     let yHit = false

//     for (let j = 0; j < PLAYERS_COUNT; j++) {
//       const playerX = pXs[j]

//       if (
//         (xP0 > playerX[0] && xP0 < playerX[1]) ||
//         (xP1 > playerX[0] && xP1 < playerX[1])
//       ) {
//         xHitPlayer = j
//       }
//     }
//     for (let j = 0; j < PLAYERS_COUNT; j++) {
//       const playerY = pYs[j]

//       if (y > playerY[0] && y < playerY[1]) {
//         yHit = true
//       }
//     }

//	    if (xHitPlayer !== null && yHit) {
//	      return {
//	        turn: gameState.turn,
//	        hitPlayer: xHitPlayer as Hit["hitPlayer"],
//	      }
//	    }
//	  }
//	  return null
//	}
func getHits(gameState TanksGameState, nextPath []Coordinates) *Hit {
	const (
		TANK_WIDTH    = 150
		TANK_HEIGHT   = 100
		PLAYERS_COUNT = 2
	)

	currentPlayerPosition := Coordinates{
		X: float64(gameState.PlayersState[gameState.Player].X + TANK_WIDTH),
		Y: float64(gameState.PlayersState[gameState.Player].Y),
	}

	hitBoundariesP0X := [2]int{
		gameState.PlayersState[0].X - TANK_WIDTH,
		gameState.PlayersState[0].X + TANK_WIDTH,
	}
	hitBoundariesP1X := [2]int{
		gameState.PlayersState[1].X - TANK_WIDTH,
		gameState.PlayersState[1].X + TANK_WIDTH,
	}
	hitBoundariesP0Y := [2]int{
		gameState.PlayersState[0].Y,
		gameState.PlayersState[0].Y + TANK_HEIGHT,
	}
	hitBoundariesP1Y := [2]int{
		gameState.PlayersState[1].Y,
		gameState.PlayersState[1].Y + TANK_HEIGHT,
	}

	pXs := [2][2]int{hitBoundariesP0X, hitBoundariesP1X}
	pYs := [2][2]int{hitBoundariesP0Y, hitBoundariesP1Y}

	for _, coord := range nextPath {
		xP0 := coord.X + currentPlayerPosition.X
		xP1 := -coord.X + currentPlayerPosition.X - float64(TANK_WIDTH*2)
		y := -coord.Y
		var xHitPlayer *int
		yHit := false

		for j := 0; j < PLAYERS_COUNT; j++ {
			playerX := pXs[j]

			if (xP0 > float64(playerX[0]) && xP0 < float64(playerX[1])) || (xP1 > float64(playerX[0]) && xP1 < float64(playerX[1])) {
				xHitPlayer = &j
			}
		}
		for j := 0; j < PLAYERS_COUNT; j++ {
			playerY := pYs[j]

			if y > float64(playerY[0]) && y < float64(playerY[1]) {
				yHit = true
			}
		}

		if xHitPlayer != nil && yHit {
			return &Hit{
				Turn:      gameState.Turn,
				HitPlayer: *xHitPlayer,
			}
		}
	}
	return nil
}

//	const updatePlayersStateHealthByHit = (gameState: TanksGameState, hit: Hit) => {
//	  const player = hit.hitPlayer
//	  const updatedHP = gameState.playersState[player].HP - 1
//	  return {
//	    ...gameState.playersState,
//	    [player]: {
//	      ...gameState.playersState[player],
//	      HP: updatedHP,
//	    },
//	  }
//	}
func updatePlayersStateHealthByHit(gameState TanksGameState, hit Hit) map[int]PlayerState {
	player := hit.HitPlayer
	updatedHP := gameState.PlayersState[player].HP - 1
	newPlayersState := gameState.PlayersState
	newPlayersState[player] = PlayerState{
		X:           newPlayersState[player].X,
		Y:           newPlayersState[player].Y,
		CannonAngle: newPlayersState[player].CannonAngle,
		HP:          updatedHP,
	}
	return newPlayersState
}

// const handleEvent = async (
//
//	event: EventsValues,
//	prevState: TanksGameState
//
//	): Promise<TanksGameState> => {
//	  const nextState = await new Promise<TanksGameState>((res, rej) => {
//	    let updated = {}
//	    let player = prevState.player
//	    switch (event) {
//	      case EVENTS.UP:
//	        updated = handlePlayersStateValueChange(prevState, "cannonAngle", -1)
//	        break
//	      case EVENTS.DOWN:
//	        updated = handlePlayersStateValueChange(prevState, "cannonAngle")
//	        break
//	      case EVENTS.LEFT:
//	        updated = handlePlayersStateValueChange(prevState, "X", -1)
//	        break
//	      case EVENTS.RIGHT:
//	        updated = handlePlayersStateValueChange(prevState, "X")
//	        break
//	      case EVENTS.FIRE:
//	        player = prevState.player === 0 ? 1 : 0
//	        const nextPath = getBulletPath(
//	          prevState.windX,
//	          prevState.playersState[prevState.player].cannonAngle
//	        )
//	        const hit = getHits(prevState, nextPath)
//	        console.log({ hit })
//	        let playersState = prevState.playersState
//	        if (hit) {
//	          playersState = updatePlayersStateHealthByHit(prevState, hit)
//	        }
//	        updated = {
//	          playersState,
//	          windX: handleWindChange(),
//	          bulletPaths: [...prevState.bulletPaths, nextPath],
//	          hits: !!hit ? [...prevState.hits, hit] : prevState.hits,
//	        }
//	        console.log({ updated })
//	        break
//	      default:
//	        res(prevState)
//	    }
//	    res({
//	      ...prevState,
//	      ...updated,
//	      player,
//	    })
//	  })
//	  return nextState
//	}
func HandleEvent(this js.Value, jsArgs []js.Value) TanksGameState {
	event := jsArgs[0]
	prevState := jsArgs[1]
	var updated map[string]interface{}
	player := prevState.Player

	switch event {
	case EVENTS["UP"]:
		prevState = handlePlayersStateValueChange(prevState, "cannonAngle", -1, 5)
	case EVENTS["DOWN"]:
		prevState = handlePlayersStateValueChange(prevState, "cannonAngle", 1, 5)
	case EVENTS["LEFT"]:
		prevState = handlePlayersStateValueChange(prevState, "X", -1, 5)
	case EVENTS["RIGHT"]:
		prevState = handlePlayersStateValueChange(prevState, "X", 1, 5)
	case EVENTS["FIRE"]:
		player = 1 - prevState.Player
		nextPath := getBulletPath(prevState.WindX, prevState.PlayersState[prevState.Player].CannonAngle)
		hit := getHits(prevState, nextPath)

		playersState := prevState.PlayersState
		if hit != nil {
			playersState = updatePlayersStateHealthByHit(prevState, *hit)
		}
		updated = map[string]interface{}{
			"playersState": playersState,
			"windX":        handleWindChange(),
			"bulletPaths":  append(prevState.BulletPaths, nextPath),
			"hits": func() []Hit {
				if hit != nil {
					return append(prevState.Hits, *hit)
				}
				return prevState.Hits
			}(),
		}
	default:
		return prevState
	}

	prevState.Player = player
	for k, v := range updated {
		switch k {
		case "playersState":
			prevState.PlayersState = v.(map[int]PlayerState)
		case "windX":
			prevState.WindX = v.(int)
		case "bulletPaths":
			prevState.BulletPaths = v.([][]Coordinates)
		case "hits":
			prevState.Hits = v.([]Hit)
		}
	}
	return prevState
}

func Test() string {
	fmt.Println("Hello WASM")
	return "Hello WASM"
}

func main() {
	// settings := TanksDefaultSettings{Width: 800, Y: 300, MaxHP: 100}
	// gameState := getDefaultGameState(settings)
	// fmt.Printf("Initial Game State: %+v\n", gameState)
	// gameState = handleEvent(UP, gameState)
	// fmt.Printf("Game State after UP event: %+v\n", gameState)

	// Create a channel to keep the Go program alive
	done := make(chan struct{}, 0)

	// Expose the Go function `fibonacciSum` to JavaScript
	js.Global().Set("handleEvent", js.FuncOf(HandleEvent))
	js.Global().Set("getDefaultGameState", js.FuncOf(GetDefaultGameState))

	// Block the program from exiting
	<-done
}
