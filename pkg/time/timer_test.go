package time

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	t.Run("Test1", func(t *testing.T) {
		ch := make(chan int)
		done := make(chan struct{})

		fmt.Println("Запустил timer")
		go Timer(ch, done, 5)
		go func() {
			ch <- 5
			done <- struct{}{}
		}()

		val := <-ch

		assert.Equal(t, 5, val)
	})

	t.Run("Test2", func(t *testing.T) {
		ch := make(chan int)
		done := make(chan struct{})

		go Timer(ch, done, 5)
		go func() {
			time.Sleep(4 * time.Second)
			ch <- 5
			done <- struct{}{}
		}()
		val := <-ch
		assert.Equal(t, 5, val)
	})

	t.Run("Test3", func(t *testing.T) {
		ch := make(chan int)
		done := make(chan struct{})

		go Timer(ch, done, 5)
		go func() {
			time.Sleep(6 * time.Second)
			ch <- 5
			done <- struct{}{}
		}()
		val := <-ch
		assert.Equal(t, -1, val)
	})

	t.Run("Test4", func(t *testing.T) {
		ch := make(chan int)
		done := make(chan struct{})

		go Timer(ch, done, 2)
		go func() {
			time.Sleep(5 * time.Second)
			ch <- 5
			done <- struct{}{}
		}()
		val := <-ch
		assert.Equal(t, -1, val)
	})
}

func TestFakeTimer(t *testing.T) {
	t.Run("Test1", func(t *testing.T) {
		ch := make(chan int)

		startTime := time.Now()
		go FakeTimer(ch)
		val := <-ch
		endTime := time.Now()
		log.Print(
			"Fake timer runs: ",
			float64(endTime.Sub(startTime).Milliseconds())/1000.0,
			" seconds")
		assert.Equal(t, -1, val)
	})
}
