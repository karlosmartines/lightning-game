package game

import (
	"math/rand"
	"time"
)

func Play(evenBet bool) bool {
	rand.Seed(time.Now().Unix())
	evenWin := rand.Intn(36)%2 == 0
	return evenWin == evenBet
}
