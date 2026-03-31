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
	d.cond.L.Lock()
	d.balance += amount
	d.cond.L.Unlock()
	d.cond.Broadcast()
}

func (d *Donation) WaitForGoal(goal int) int {
	// This approach doesn't make loop busy
	d.cond.L.Lock()
	defer d.cond.L.Unlock()

	for d.balance < goal {
		d.cond.Wait()
	}

	return d.balance
}

func (d *Donation) Balance() int {
	d.cond.L.Lock()
	defer d.cond.L.Unlock()

	return d.balance
}
