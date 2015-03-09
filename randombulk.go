package randombulk

import (
	"bufio"
	"math"
	"math/rand"
	"os"
	"time"
)

var latinLetters =[]byte{
	65, 66, 67, 68, 69, 70, 71, 72, 73, 74,
	75, 76, 77, 78, 79, 80, 81, 82, 83, 84,
	85, 86, 87, 88, 89, 90, 97, 98, 99, 100,
	101, 102, 103, 104, 105, 106, 107, 108, 109, 110,
	111, 112, 113, 114, 115, 116, 117, 118, 119, 120,
	121, 122,
}

var Frequencies struct {
	Frequently string
	Sometimes  string
	Rarely     string
}

func init() {
	Frequencies.Frequently = "frequently"
	Frequencies.Sometimes = "sometimes"
	Frequencies.Rarely = "rarely"
	rand.Seed(time.Now().UTC().UnixNano())
}

type Inclusion struct {
	byteData   []byte
	stringData string
	length     int64
	b1, b2     int
}

func mapFrequencyToBoundaries(f string) (int, int) {
	switch f {
	case Frequencies.Frequently:
		return 10000, 9000
	case Frequencies.Sometimes:
		return 10000, 5000
	case Frequencies.Rarely:
		return 10000, 2000
	default:
		return 10000, 5000
	}
}

func NewInclusionFromString(data string, frequency string) Inclusion {
	b1, b2 := mapFrequencyToBoundaries(frequency)
	return Inclusion{[]byte(data), data, int64(len(data)), b1, b2}
}

func NewInclusionFromBytes(data []byte, frequency string) Inclusion {
	b1, b2 := mapFrequencyToBoundaries(frequency)
	return Inclusion{data, string(data), int64(len(data)), b1, b2}
}

type RandomBulk struct {
	fileBufferSize int
}

func NewRandomBulk() *RandomBulk {
	return &RandomBulk{64 * 1024 * 1024}
}

func readASCIIBytes(bs []byte) {
	bsl := len(bs)
	for i := 0; i < bsl; i++ {
		bs[i] = latinLetters[rand.Intn(52)]
	}
}

func readAnyBytes(bs []byte) {
	var (
		l    int   = len(bs)
		ind  int   = 0
		ri64 int64 = 0
	)
	for {
		ri64 = rand.Int63()
		for i := 0; i < 8; i++ {
			bs[ind] = byte(ri64)
			if l--; l == 0 {
				return
			}
			ind++
			ri64 >>= 8
		}
	}
}

func (rb *RandomBulk) DumpToFile(filePath string, minFileSize int64, isASCII bool, inclusions []Inclusion) (fileLength int64, err error) {
	var (
		f    *os.File
		file *bufio.Writer
	)
	if f, err = os.Create(filePath); err != nil {
		fileLength = -1
		return
	}
	file = bufio.NewWriterSize(f, rb.fileBufferSize)
	defer func() {
		file.Flush()
		f.Close()
	}()
	var (
		stublen int64  = int64(128)
		stub    []byte = make([]byte, stublen)
		ncycles int    = int(math.Ceil(float64(minFileSize) / float64(stublen)))
	)
	var readBytes func([]byte)
	if isASCII {
		readBytes = readASCIIBytes
	} else {
		readBytes = readAnyBytes
	}
	for i := 0; i < ncycles; i++ {
		readBytes(stub)
		if _, err = file.Write(stub); err != nil {
			fileLength = -1
			return
		}
		fileLength += stublen
		for _, incl := range inclusions {
			if rand.Intn(incl.b1) < incl.b2 {
				if _, err = file.Write(incl.byteData); err != nil {
					fileLength = -1
					return
				}
				fileLength += incl.length
			}
		}
	}
	return
}
