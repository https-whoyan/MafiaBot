package util

import "math/rand"

const (
	minNum = 1000
	maxNum = 9999
)

func GetRandomNumber() int {
	return minNum + rand.Intn(maxNum-minNum)
}
