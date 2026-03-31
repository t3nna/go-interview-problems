package main

import "sync"

type Donation struct {
	cond    *sync.Cond
	balance int
}

func NewDonation() *Donation {
	return &Donation{
		cond: sync.NewCond(&sync.Mutex{}),
	}
}

func (d *Donation) Add(amount int) {
}

func (d *Donation) WaitForGoal(goal int) int {
	return 0
}

func (d *Donation) Balance() int {
	return 0
}
