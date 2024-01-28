# tempmap

tempmap provides an interface for a generic, thread-safe map with values that expire after the provided amount of time

## Examples

```go
package main

import (
	"fmt"
	"time"

	"github.com/harmonherring/tempmap"
)

func main() {
	tempMap := tempmap.New[int, string]()
	defer tempMap.Close()

	tempMap.Put(0, "Hello, World!", 1*time.Second)

	if value, exists := tempMap.Get(0); exists {
		println(value)
	} else {
		fmt.Println("value does not exist")
	}

	time.Sleep(2 * time.Second)

	if value, exists := tempMap.Get(0); exists {
		println(value)
	} else {
		fmt.Println("value does not exist")
	}
}
```

## Future
- Evaluate [orcaman/concurrent-map](https://github.com/orcaman/concurrent-map) for better concurrency performance
