package game

func (g *Game) Day(ch chan<- Signal) {
	select {
	case <-g.ctx.Done():
		return
	default:
		g.Lock()
		g.SetState(DayState)
		g.Unlock()

	}
}
