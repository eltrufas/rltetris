package rltetris
import "math/rand"

func GetGreedyAction(weights map[uint32][]float64, state []bool, actions []uint32) uint32 {
  maxA := actions[0]
  maxQ := Q(state, maxA, weights[maxA])
  for _, a := range actions[1:] {
    q := Q(state, a, weights[a])
    if q > maxQ {
      maxA = a
      maxQ = q
    }
  }
	return maxA
}

func GetEGreedyAction(w map[uint32][]float64, state []bool, actions []uint32) uint32 {
  if rand.Float64() > 0.9 {
    i := rand.Intn(len(actions))
    return actions[i]
  } else {
    return GetGreedyAction(w, state, actions)
  }
}

func Q(state []bool, action uint32, weights []float64) float64 {
	var acc float64
	for i := 0; i < len(state); i++ {
    if state[i] {
      acc += weights[i]
    }
	}

	return acc
}

func Sarsa(w map[uint32][]float64, episodes int, alpha, discount float64) {
	for i := 0; i < episodes; i++ {
		game := CreateTetris()
		s := game.GetState()
		a := GetEGreedyAction(w, s, game.LegalAction())
		for !game.Terminal() {
			r := float64(game.Transition(a))
			sPrime := game.GetState()
			if game.Terminal() {
				actionWeights := w[a]
				change := alpha * (r - Q(s, a, actionWeights))
				for j := 0; j < len(w); j++ {
          if s[j] {
            actionWeights[j] += change
          }
				}
				break
			}
			aPrime := GetEGreedyAction(w, sPrime, game.LegalAction())
			actionWeights := w[a]
			change := alpha * (r + discount*Q(sPrime, aPrime, actionWeights))
			for j := 0; j < len(s); j++ {
        if s[j] {
          actionWeights[j] += change
        }
			}

			s = sPrime
			a = aPrime
		}
	}
}
