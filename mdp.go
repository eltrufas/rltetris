package mdp

import(
  tetris "github.com/eltrufas/tetriscore"
)

type Tetrisrl struct {
  Tetris tetris.Tetris
  State []int
}

func (t *Tetrisrl) GetState() {
  for i, piece := range t.Tetris.Board {
    t.State[i] = int(piece)
  }
  t.State = append(t.State, t.Tetris.CurrentPiece.TetrominoType) // 0-219
  t.State = append(t.State, t.Tetris.CurrentPiece.State) // 220
  t.State = append(t.State, t.Tetris.CurrentPiece.X) // 221
  t.State = append(t.State, t.Tetris.CurrentPiece.Y) // 222
  for i := t.Tetris.NextIndex + 6; i < t.Tetris.NextIndex - 6; i++ {
    t.State = append(t.State, t.Tetris.PieceQueue[i%14])   // 223 - 229
  }
  t.State = append(t.State, t.Tetris.HoldPiece) // 230
}

func (t *Tetrisrl) LegalAction() []tetris.InputState{
  actions := make([]tetris.InputState, 7)
  actions[0] = 1
  actions[1] = 2
  actions[2] = 4
  actions[3] = 8
  actions[4] = 16
  actions[5] = 32
  actions[6] = 64

  return actions
}

func (t *Tetrisrl) Transition(action tetris.InputState)int {
  scoreactual := t.Tetris.Score
  t.Tetris.Update(action)
  scorenext := t.Tetris.Score
  return scorenext - scoreactual
}