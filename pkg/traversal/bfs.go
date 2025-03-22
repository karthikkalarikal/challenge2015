package traversal

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/karthikkalarikal/Qube/models"
	"github.com/karthikkalarikal/Qube/pkg/util"
)

type Node struct {
	Entity string
	Type   string
	From   *Node
	Link   string
}

func NewNode(start, target string) {
	node := bfs(start, target)
	printPath(node)
}

func bfs(start, end string) *Node {
	var visited sync.Map
	var once sync.Once
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	type safeQueue struct {
		mu    sync.Mutex
		cond  *sync.Cond
		items []Node
	}

	queue := &safeQueue{items: make([]Node, 0)}
	queue.cond = sync.NewCond(&queue.mu)
	result := make(chan *Node, 1)

	queue.mu.Lock()
	queue.items = append(queue.items, Node{Entity: start, Type: "person"})
	queue.cond.Signal()
	queue.mu.Unlock()

	worker := func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				queue.mu.Lock()
				// Wait for work or cancellation
				for len(queue.items) == 0 {
					if ctx.Err() != nil {
						queue.mu.Unlock()
						return
					}
					queue.cond.Wait()
				}
				// Dequeue node
				current := queue.items[0]
				queue.items = queue.items[1:]
				queue.mu.Unlock()

				// Skip if already processed
				if _, loaded := visited.LoadOrStore(current.Entity, true); loaded {
					continue
				}

				// Early termination check
				if current.Entity == end {
					once.Do(func() {
						result <- &current
						cancel()
					})
					return
				}

				// Process node
				if current.Type == "person" {
					var person models.Actor
					if err := util.GetByURL(current.Entity, &person); err != nil {
						if checkErrorForbidden(err) {
							visited.Delete(current.Entity)
						}
						log.Printf("Retryable error on %s: %v", current.Entity, err)
						continue
					}

					for _, credit := range person.Movies {
						child := Node{
							Entity: credit.URL,
							Type:   "movie",
							From:   &current,
							Link:   "Movie: " + credit.Name + " (" + credit.Role + ")",
						}
						queue.mu.Lock()
						queue.items = append(queue.items, child)
						queue.cond.Signal()
						queue.mu.Unlock()
					}
				} else {
					var movie models.Movie
					if err := util.GetByURL(current.Entity, &movie); err != nil {
						if checkErrorForbidden(err) {
							visited.Delete(current.Entity)
						}
						log.Printf("Retryable error on %s: %v", current.Entity, err) //since some of the links are not accessible, the retry logic is specific.
						continue
					}

					for _, cast := range movie.Cast {
						child := Node{
							Entity: cast.URL,
							Type:   "person",
							From:   &current,
							Link:   "Cast: " + cast.Name + " (" + cast.Role + ")",
						}
						queue.mu.Lock()
						queue.items = append(queue.items, child)
						queue.cond.Signal()
						queue.mu.Unlock()
					}
					for _, crew := range movie.Crew {
						child := Node{
							Entity: crew.URL,
							Type:   "person",
							From:   &current,
							Link:   "Crew: " + crew.Name + " (" + crew.Role + ")",
						}
						queue.mu.Lock()
						queue.items = append(queue.items, child)
						queue.cond.Signal()
						queue.mu.Unlock()
					}
				}
			}
		}
	}

	// Start workers
	numWorkers := 30
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go worker()
	}

	go func() {
		wg.Wait()
		once.Do(func() { close(result) })
	}()

	select {
	case res := <-result:
		return res
	case <-ctx.Done():
		return nil
	}
}

func printPath(node *Node) {
	if node == nil {
		fmt.Println("No path found.")
		return
	}

	var path []*Node
	for node != nil {
		path = append([]*Node{node}, path...)
		node = node.From
	}
	for i, n := range path {
		if i == 0 {
			fmt.Printf("%d. Start: %s\n", i+1, n.Entity)
		} else {
			fmt.Printf("%d. %s\n", i+1, n.Link)
		}
	}
}

func checkErrorForbidden(err error) bool {
	return strings.Contains(err.Error(), "access denied for URL:")
}
