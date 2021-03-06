# randombulk
A library to generate random binary files

[![GoDoc](https://godoc.org/github.com/pkorotkov/randombulk?status.svg)](https://godoc.org/github.com/pkorotkov/randombulk)

## Installation
    go get -u -v github.com/pkorotkov/randombulk
    
## Usage sample
```go
import (
    "fmt"
    
    . "github.com/pkorotkov/randombulk"
)

func main() {
    bulk := NewRandomBulk()
    incls := []Inclusion{
        NewInclusionFromString("D0zWholeBr7q4W9", Frequencies.Sometimes),
        NewInclusionFromString("Kw87_uiX2Y", Frequencies.Rarely),
        NewInclusionFromBytes([]byte("RWo-45vZl"), Frequencies.Sometimes),
    }
    fl, _ := bulk.DumpToFile("data.bin", 100*1024*1024, false, incls)
    fmt.Printf("File's byte length: %d\n", fl)
}
```
