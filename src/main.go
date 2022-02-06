package main

import (
	"fmt"
	"sync"
	"time"
)

type Limit struct {
	count             float64
	durationInSeconds int
}

type Limiter struct {
	mu        sync.Mutex
	limit     Limit
	counter   float64
	startTime time.Time
	endTime   time.Time
}

func (lim *Limiter) GetLimit() Limit {
	lim.mu.Lock()
	defer lim.mu.Unlock()
	return lim.limit
}

func (lim *Limiter) GetCount() float64 {
	lim.mu.Lock()
	defer lim.mu.Unlock()
	return lim.counter
}

func (lim *Limiter) SetCount(count float64) {
	lim.mu.Lock()
	lim.counter = count
	lim.mu.Unlock()
}

func (lim *Limiter) IncreaseCount() {
	lim.mu.Lock()
	lim.counter++
	lim.mu.Unlock()
}

func NewLimiter(r Limit) *Limiter {
	return &Limiter{
		limit:     r,
		counter:   0,
		startTime: time.Now(),
		endTime: time.Now().Local().Add(
			time.Duration(r.durationInSeconds * 1e9),
		),
	}
}

func (lim *Limiter) Allow() bool {
	isAllowed := false

	if lim.endTime.Equal(time.Now()) || lim.endTime.Before(time.Now()) {
		lim.SetCount(0)
		lim.startTime = time.Now()
		lim.endTime = time.Now().Local().Add(
			time.Duration(lim.GetLimit().durationInSeconds * 1e9),
		)
	}

	lim.IncreaseCount()

	if lim.GetCount() <= lim.GetLimit().count {
		isAllowed = true
	}

	return isAllowed
}

func main() {
	limiter := NewLimiter(
		Limit{
			2,
			1,
		},
	)

	fmt.Println("Fixed window limiter settings:")
	fmt.Println("requests limit -", limiter.GetLimit().count)
	fmt.Println("duration in seconds -", limiter.GetLimit().durationInSeconds)
	fmt.Println()

	fmt.Println("is allowed request:", limiter.Allow())
	fmt.Println("is allowed request:", limiter.Allow())
	fmt.Println("is allowed request:", limiter.Allow())

	fmt.Println()
	fmt.Println("sleep one second")
	time.Sleep(time.Second)
	fmt.Println()

	fmt.Println("is allowed request:", limiter.Allow())
	fmt.Println("is allowed request:", limiter.Allow())

	fmt.Println()
	fmt.Println("sleep one second")
	time.Sleep(time.Second)
	fmt.Println()

	fmt.Println("is allowed request:", limiter.Allow())
	fmt.Println("is allowed request:", limiter.Allow())
	fmt.Println("is allowed request:", limiter.Allow())
}
