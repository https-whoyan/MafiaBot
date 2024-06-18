package game

type GameLogger interface {
	InitNewGame(*Game) error
}
