package mdp

func GetEGreedyAction(w map[uint32][]float64, state []int, actions []uint32) uint32 {

	return 0
}

func Q(state []int, action uint32, weights []float64) float64 {
	var acc float64
	for i := 0; i < len(state); i++ {
		acc += float64(state[i]) * weights[i]
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
					actionWeights[j] += change * float64(s[j])
				}
				break
			}
			aPrime := GetEGreedyAction(w, sPrime, game.LegalAction())
			actionWeights := w[a]
			change := alpha * (r + discount*Q(sPrime, aPrime, actionWeights))
			for j := 0; j < len(s); j++ {
				actionWeights[j] += change * float64(s[j])
			}

			s = sPrime
			a = aPrime
		}
	}
}
