# randombulk
A library to generate random binary files

## Installation
    go get -u -v github.com/pkorotkov/randombulk
    
## Usage sample
```go
import (
    "fmt"
    
    _ "github.com/pkorotkov/randombulk"
)

func main() {
    bulk := NewRandomBulk()
    incls := []Inclusion{
        NewInclusionFromString("D0zWholeBr7q4W9", Frequencies.Sometimes),
        NewInclusionFromString("Kw87_uiX2Y", Frequencies.Rarely),
        NewInclusionFromString("RWo-45vZl", Frequencies.Sometimes),
    }
    fl, _ := bulk.DumpToFile("data.bin", 100*1024*1024, false, incls)
    fmt.Printf("File's byte length: %d\n", fl)
}
```
