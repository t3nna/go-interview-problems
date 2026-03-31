# Thread Cond
Implement a donation goal tracker using `sync.Cond`.

Requirements:

- `NewDonation()` creates a tracker with a balance of `0`.
- `Add(amount int)` increases the balance and notifies all waiting goroutines.
- `WaitForGoal(goal int) int` blocks until the balance is at least `goal`, then returns the reached balance.
- `Balance() int` returns the current balance safely.

The solution must avoid busy loops and support multiple listeners waiting for donation goals concurrently.

## Tags
`Concurrency`
