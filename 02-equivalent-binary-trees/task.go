package main

import (
	"golang.org/x/tour/tree"
)

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func WalkBfs(t *tree.Tree, ch chan int) {
	queue := []*tree.Tree{t}

	for len(queue) > 0 {
		point := queue[0]
		queue = queue[1:]

		if point == nil {
			continue
		}

		ch <- point.Value

		queue = append(queue, point.Left)
		queue = append(queue, point.Right)
	}

	close(ch)
}

// DFS
func Walk(t *tree.Tree, ch chan int) {
	defer close(ch)
	doWalk(t, ch)

}

func doWalk(t *tree.Tree, ch chan int) {
	if t == nil {
		return
	}
	doWalk(t.Left, ch)
	ch <- t.Value
	doWalk(t.Right, ch)

}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	ch1, ch2 := make(chan int), make(chan int)
	go Walk(t1, ch1)
	go Walk(t2, ch2)

	for {
		val0, ok0 := <-ch1
		val1, ok1 := <-ch2

		if ok0 != ok1 {
			return false
		}
		if val0 != val1 {
			return false
		}

		if !ok0 {
			return true
		}
	}

	return false
}
