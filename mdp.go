package rltetris

import (
	tetris "github.com/eltrufas/tetriscore"
//  "fmt"
  "math/rand"
)

type Tetrisrl struct {
	Tetris *tetris.Tetris
  lastAction uint32
	State  []bool
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

	t.State = append(t.State, pieces...)         // 220-226
	t.State = append(t.State, rotation...)       // 227-231
	t.State = append(t.State, x...)              // 232-242
	t.State = append(t.State, y...)              // 243-265

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

  for i:= 0; i< 10; i++ {
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
      if t.Tetris.Board[j * 10 + i] == tetris.Empty {
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

  for _, val := range(t.Tetris.It) {
    t.State = append(t.State, val != 0)
  }

	return t.State
}

func (t *Tetrisrl) LegalAction() []uint32 {
	actions := make([]uint32, 7)
	actions[0] = 1
	actions[1] = 2
	actions[2] = 4
	actions[3] = 8
	actions[4] = 16
	actions[5] = 32
	actions[6] = 64

  for i := range actions {
    j := rand.Intn(i + 1)
    actions[i], actions[j] = actions[j], actions[i]
  }

	return actions
}

func (t *Tetrisrl) Transition(action uint32) float64 {
	scoreactual := t.Tetris.Score
	t.Tetris.Update(tetris.InputState(action))
	scorenext := t.Tetris.Score
  reward := float64(scorenext - scoreactual)
  if (action == 0 && t.Tetris.CurrentPiece.X <= 0) || (action == 1 && t.Tetris.CurrentPiece.X >= 9){
    reward -= 20
  }
  if action == t.lastAction {
    reward -= 20
  }

  t.lastAction = action
	return reward
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
