# bencode

`bencode` is a Go language binding for encoding and decoding data in the bencode format that is used by the BitTorrent peer-to-peer file sharing protocol.

## Features

- **Encode**: Encode data to bencode format.
- **Decode**: Decode bencode data to Go data types.

## Installation

To include `bencode` in your project, use Go modules:

```bash
go get github.com/codescalersinternships/bencode-eyadhussein
```

## Usage

### Encoding

```go
package main

import (
    "fmt"
    "log"
    bencode "github.com/codescalersinternships/bencode-eyadhussein/pkg"
)

func main() {
    data := map[string]any{
        "string": "Hello, World!",
        "integer": 42,
        "list": []any{}{
            "foo",
            "bar",
            "baz",
        },
        "dict": map[string]any{
            "key": "value",
        },
    }

    encoded, err := bencode.Encode(data)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(encoded)
}
```

### Decoding

```go
package main

import (
    "fmt"
    "log"
    bencode "github.com/codescalersinternships/bencode-eyadhussein/pkg"
)

func main() {
    data := "d6:string13:Hello, World!7:integeri42e4:listl3:foo3:bar3:baze4:dictd3:key5:valueee"

    decoded, err := bencode.Decode(data)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(decoded)
}
```

## Linting

```bash
make lint
```

## Testing

```bash
make test
```
