# bencode

Bencode implementation in Go

Inspired by [marksamman/bencode](https://github.com/marksamman/bencode)

## Install
```
$ go get github.com/impos/bencode
```

## Usage

### Encode
bencode.Encode takes a map[string]interface{} as argument and returns ([]byte, error). Example:

```go
package main

import (
	"fmt"
	"log"

	"github.com/impos/bencode"
)

func main() {
	dict := map[string]interface{}{}
	dict["int key"] = 123456
	dict["string key"] = "hello world"
	dict["bytes key"] = []byte("hello world")
	dict["list key"] = []interface{}{1, "str", 2}
	dict["dict key"] = map[string]interface{}{
		"1": 1,
		"2": "2",
	}
	
	result, err := bencode.Encode(dict)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("bencode encoded dict: %s\n", result)
}
```

### Decode
bencode.Decode takes an io.Reader as argument and returns (map[string]interface{}, error). Example:
```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/impos/bencode"
)

func main() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	dict, err := bencode.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("string: %s\n", dict["string key"].(string))
	fmt.Printf("int: %d\n", dict["int key"].(int64))
}
```

