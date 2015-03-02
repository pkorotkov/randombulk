package randombulk

import (
	"bufio"
	"io"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/dustin/randbo"
)

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
	fileBufferSize     int
	anyRandSource      io.Reader
	latinLetters       []byte
	latinLettersLength int
}

func NewRandomBulk() *RandomBulk {
	return &RandomBulk{
		64 * 1024 * 1024,
		randbo.New(),
		[]byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"),
		52,
	}
}

func (rb *RandomBulk) readASCII(bs []byte) {
	bsl := len(bs)
	for i := 0; i < bsl; i++ {
		bs[i] = rb.latinLetters[rand.Intn(rb.latinLettersLength)]
	}
}

func (rb *RandomBulk) readAny(bs []byte) {
	rb.anyRandSource.Read(bs)
}

func (rb *RandomBulk) DumpToFile(filePath string, minFileSize int64, isASCII bool, inclusions []Inclusion) (fileLength int64, err error) {
	var (
		f    *os.File
		file *bufio.Writer
	)
	if f, err = os.Create(filePath); err != nil {
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
		readBytes = rb.readASCII
	} else {
		readBytes = rb.readAny
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
