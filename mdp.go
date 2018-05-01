package rltetris

import (
	"fmt"
	tetris "github.com/eltrufas/tetriscore"
	"math/rand"
)

type Tetrisrl struct {
	Tetris     *tetris.Tetris
	lastAction uint32
	State      []bool
}

func appendOneHot(arr []byte, min, max, val int) []byte {
	if val < min || val > max {
		panic(fmt.Sprintf("Valor fuera de rango, min: %v, max: %v, val: %v", min, max, val))
	}

	for i := min; i < val; i++ {
		arr = append(arr, 0)
	}

	arr = append(arr, 1)

	for i := val; i < max; i++ {
		arr = append(arr, 0)
	}

	return arr
}
func GetByteState(t *tetris.Tetris) (state []byte) {
	state = make([]byte, 0, 300)
	for _, val := range t.Board {
		if val != tetris.Empty {
			state = append(state, 1)
		} else {
			state = append(state, 0)
		}
	}

	state = appendOneHot(state, -3, 9, t.CurrentPiece.X)
	state = appendOneHot(state, 0, 21, t.CurrentPiece.Y)
	state = appendOneHot(state, 0, 3, t.CurrentPiece.State)
	state = appendOneHot(state, 0, 6, t.CurrentPiece.TetrominoType)

	for i := 0; i < 6; i++ {
		state = appendOneHot(state, 0, 6, t.PieceQueue[(t.NextIndex+i)%14])
	}

	if !t.Held {
		state = appendOneHot(state, 0, 7, 7)
	} else {
		state = appendOneHot(state, 0, 7, t.HoldPiece)
	}

	return
}
func (t *Tetrisrl) GetByteState() (state []byte) {
	state = GetByteState(t.Tetris)
	return
}

func (t *Tetrisrl) GetState() []bool {
	t.State = make([]bool, 0)
	for i, _ := range t.Tetris.Board {
		t.State = append(t.State, t.Tetris.Board[i] != tetris.Empty)
	} // 0-219

	// pieza actual
	pieces := make([]bool, 7)
	pieces[t.Tetris.CurrentPiece.TetrominoType] = true
	rotation := make([]bool, 4)
	rotation[t.Tetris.CurrentPiece.State] = true
	x := make([]bool, 11)
	idx := t.Tetris.CurrentPiece.X
	if idx < 0 {
		idx = 0
	}
	x[idx] = true
	y := make([]bool, 22)

	idx = t.Tetris.CurrentPiece.Y
	if idx < 0 {
		idx = 0
	}
	y[idx] = true

	t.State = append(t.State, pieces...)   // 220-226
	t.State = append(t.State, rotation...) // 227-231
	t.State = append(t.State, x...)        // 232-242
	t.State = append(t.State, y...)        // 243-265

	next := make([]bool, 7)
	for i := t.Tetris.NextIndex + 6; i < t.Tetris.NextIndex-6; i++ {
		next[t.Tetris.PieceQueue[i%14]] = true
		t.State = append(t.State, next...)
		next[t.Tetris.PieceQueue[i%14]] = false
	} // 266-308

	hold := make([]bool, 8)
	if t.Tetris.Held == false {
		t.State = append(t.State, hold...) // 309-316
	} else {
		hold[t.Tetris.HoldPiece] = true
		t.State = append(t.State, hold...) // 309-316
	}
	level := make([]bool, 15)
	level[t.Tetris.Level] = true
	t.State = append(t.State, level...) // 317-332

	for i := 0; i < 10; i++ {
		height := 0
		for j := 0; j < 22; j++ {
			if t.Tetris.Board[i+j*10] != tetris.Empty {
				height = j
				break
			}
		}

		if height > 19 {
			height = 19
		}

		gaps := 0
		for j := height; j < 22; j++ {
			if t.Tetris.Board[j*10+i] == tetris.Empty {
				gaps++
			}
		}

		gapsVec := make([]bool, 23)
		gapsVec[gaps] = true
		t.State = append(t.State, gapsVec...)

		heightVec := make([]bool, 20)
		heightVec[height] = true
		t.State = append(t.State, heightVec...)
	} // 333-533
	ghost := t.Tetris.GhostPiece()
	x = make([]bool, 11)
	idx = ghost.X
	if idx < 0 {
		idx = 0
	}
	x[idx] = true
	y = make([]bool, 22)
	idx = ghost.Y
	if idx < 0 {
		idx = 0
	}
	y[idx] = true
	t.State = append(t.State, x...)
	t.State = append(t.State, y...)

	for _, val := range t.Tetris.It {
		t.State = append(t.State, val != 0)
	}

	return t.State
}

func (t *Tetrisrl) LegalAction() []uint32 {
	actions := make([]uint32, 8)
	actions[0] = 1
	actions[1] = 2
	actions[2] = 4
	actions[3] = 8
	actions[4] = 16
	actions[5] = 32
	actions[6] = 64
	//actions[7] = 0

	for i := range actions {
		j := rand.Intn(i + 1)
		actions[i], actions[j] = actions[j], actions[i]
	}

	return actions
}

func (t *Tetrisrl) Transition(action uint32) float64 {
	reward := t.Tetris.Step(tetris.InputState(action))

	return float64(reward)
}

func (t *Tetrisrl) Terminal() bool {
	return t.Tetris.FlagLoss
}

func CreateTetris() Tetrisrl {
	var t Tetrisrl
	t.Tetris = tetris.CreateTetris()
	t.State = t.GetState()
	return t
}
