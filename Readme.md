# Degrees of Separation CLI

A Go CLI application that calculates the shortest path between two actors/movie professionals using data from [Moviebuff](https://www.moviebuff.com/). The solution uses concurrent BFS (Breadth-First Search) to efficiently find connections through movies and crew members.

## Features

- **Concurrent BFS implementation in Graph** with worker pooling
- **Rate limiting** to handle API throttling  (currently 1000 requests per second with a 1000 as burst)
- **Smart retry mechanism** for failed requests (some entries are not accessible)
- **Path reconstruction** showing movie connections 
- **Efficient memory management** with sync.Map and compact storage

## Installation

```go build -o degrees```

```./degrees <start-person> <end-person>```

```./degrees -a amitabh-bachchan -b robert-de-niro```


