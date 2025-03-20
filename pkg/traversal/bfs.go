package traversal

import (
	"fmt"
	"log"
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
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}
	gofer := make(chan struct{}, 10)
	queue := make([]Node, 1)
	queue[0] = Node{Entity: start, Type: "person"}

	for len(queue) > 0 {
		l := len(queue)
		for i := 0; i < l; i++ {

			current := queue[0]
			queue = queue[1:]
			log.Println(current.Entity, " ", current.Type)
			if current.Entity == end {
				return &current
			}
			visited.Store(current.Entity, true)
			wg.Add(1)
			gofer <- struct{}{}
			go func(current Node) {
				defer func() {
					wg.Done()
					<-gofer
				}()
				if current.Type == "person" {
					var person models.Actor

					if err := util.GetByURL(current.Entity, &person); err != nil {
						visited.Store(current.Entity, true)
						log.Printf("Skipping inaccessible node: %s (%v)", current.Entity, err)

					}

					for _, credit := range person.Movies {
						node := Node{
							Entity: credit.URL,
							Type:   "movie",
							From:   &current,
							Link:   "Movie: " + credit.Name + " (" + credit.Role + ")",
						}

						if _, ok := visited.Load(credit.URL); !ok {
							mu.Lock()
							queue = append(queue, node)
							mu.Unlock()
						}
					}
				} else {
					var movie models.Movie

					if err := util.GetByURL(current.Entity, &movie); err != nil {
						visited.Store(current.Entity, true)
						log.Printf("Skipping inaccessible node: %s (%v)", current.Entity, err)

					}

					log.Println(movie.URL)
					for _, cast := range movie.Cast {

						node := Node{
							Entity: cast.URL,
							Type:   "person",
							From:   &current,
							Link:   "Cast: " + cast.Name + " (" + cast.Role + ")",
						}
						if _, ok := visited.Load(cast.URL); !ok {
							mu.Lock()
							queue = append(queue, node)
							mu.Unlock()
						}

					}
					for _, crew := range movie.Crew {

						node := Node{
							Entity: crew.URL,
							Type:   "person",
							From:   &current,
							Link:   "Crew: " + crew.Name + " (" + crew.Role + ")",
						}

						if _, ok := visited.Load(crew.URL); !ok {
							mu.Lock()
							queue = append(queue, node)
							mu.Unlock()
						}
					}
				}
			}(current)
		}
		wg.Wait()

	}
	return nil

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
