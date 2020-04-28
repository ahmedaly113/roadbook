package gps


import(
	"time"
	"sync"
	"math/rand"
	"github.com/ahmedaly113/roadbook/model"
)

func Simulator(state *model.Model,stateMu *sync.Mutex) {
	t := time.NewTicker(250 * time.Millisecond)
	for _ = range t.C {
		stateMu.Lock()
		// meander Heading
		state.Heading += rand.Float32() - 0.5
		if state.Heading < 0 || state.Heading > 359.9 {
			state.Heading = 0
		}

		// meander speed
		state.Speed += rand.Float32()*3 - 1.5
		if state.Speed > 70 {
			state.Speed = 70
		}
		if state.Speed < 15 {
			state.Speed = 15
		}

		// advance waypoint occasionally
		if rand.Intn(10) == 0 {
			state.Idx++

			if state.Idx >= len(state.Book)-4 {
				state.Idx = 0
			}

			// reset distance
			state.Distance = state.Book[state.Idx].Distance
		}
		state.Distance *= 0.95 // deplete some distance remaining

		stateMu.Unlock()
	}
}
