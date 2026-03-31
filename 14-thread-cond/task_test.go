package main

import (
	"testing"
	"time"
)

func TestDonationWaitForGoalBlocksUntilReached(t *testing.T) {
	donation := NewDonation()

	done := make(chan int, 1)
	go func() {
		done <- donation.WaitForGoal(3)
	}()

	assertNoResult(t, done)
	donation.Add(1)
	assertNoResult(t, done)
	donation.Add(1)
	assertNoResult(t, done)
	donation.Add(1)

	got := waitForResult(t, done)
	if got != 3 {
		t.Errorf("Expected reached balance to be 3, got: %d", got)
	}

	if donation.Balance() != 3 {
		t.Errorf("Expected total balance to be 3, got: %d", donation.Balance())
	}
}

func TestDonationWaitForGoalAlreadyReached(t *testing.T) {
	donation := NewDonation()
	donation.Add(5)

	done := make(chan int, 1)
	go func() {
		done <- donation.WaitForGoal(3)
	}()

	got := waitForResult(t, done)
	if got != 5 {
		t.Errorf("Expected reached balance to be 5, got: %d", got)
	}
}

func TestDonationWakesAllWaiters(t *testing.T) {
	donation := NewDonation()

	const waiters = 5
	const goal = 3

	start := make(chan struct{})
	done := make(chan int, waiters)

	for range waiters {
		go func() {
			<-start
			done <- donation.WaitForGoal(goal)
		}()
	}

	close(start)
	time.Sleep(20 * time.Millisecond)

	for range goal {
		donation.Add(1)
	}

	for i := range waiters {
		got := waitForResult(t, done)
		if got < goal {
			t.Errorf("Waiter #%d expected reached balance to be at least %d, got: %d", i+1, goal, got)
		}
	}
}

func assertNoResult(t *testing.T, ch <-chan int) {
	t.Helper()
	select {
	case got := <-ch:
		t.Fatalf("Expected waiter to still be blocked, got result: %d", got)
	case <-time.After(40 * time.Millisecond):
	}
}

func waitForResult(t *testing.T, ch <-chan int) int {
	t.Helper()
	select {
	case got := <-ch:
		return got
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Timed out waiting for waiter result")
		return 0
	}
}
