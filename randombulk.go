package randombulk

import (
	"bufio"
	"io"
	"math"
	"math/rand"
	"os"
	"time"
)

var (
	allBytes = []byte{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
		10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
		20, 21, 22, 23, 24, 25, 26, 27, 28, 29,
		30, 31, 32, 33, 34, 35, 36, 37, 38, 39,
		40, 41, 42, 43, 44, 45, 46, 47, 48, 49,
		50, 51, 52, 53, 54, 55, 56, 57, 58, 59,
		60, 61, 62, 63, 64, 65, 66, 67, 68, 69,
		70, 71, 72, 73, 74, 75, 76, 77, 78, 79,
		80, 81, 82, 83, 84, 85, 86, 87, 88, 89,
		90, 91, 92, 93, 94, 95, 96, 97, 98, 99,
		100, 101, 102, 103, 104, 105, 106, 107, 108, 109,
		110, 111, 112, 113, 114, 115, 116, 117, 118, 119,
		120, 121, 122, 123, 124, 125, 126, 127, 128, 129,
		130, 131, 132, 133, 134, 135, 136, 137, 138, 139,
		140, 141, 142, 143, 144, 145, 146, 147, 148, 149,
		150, 151, 152, 153, 154, 155, 156, 157, 158, 159,
		160, 161, 162, 163, 164, 165, 166, 167, 168, 169,
		170, 171, 172, 173, 174, 175, 176, 177, 178, 179,
		180, 181, 182, 183, 184, 185, 186, 187, 188, 189,
		190, 191, 192, 193, 194, 195, 196, 197, 198, 199,
		200, 201, 202, 203, 204, 205, 206, 207, 208, 209,
		210, 211, 212, 213, 214, 215, 216, 217, 218, 219,
		220, 221, 222, 223, 224, 225, 226, 227, 228, 229,
		230, 231, 232, 233, 234, 235, 236, 237, 238, 239,
		240, 241, 242, 243, 244, 245, 246, 247, 248, 249,
		250, 251, 252, 253, 254, 255
	}
	latinLetters = append(allBytes[65:91], allBytes[97:123]...)
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
	fileBufferSize int
}

func NewRandomBulk() *RandomBulk {
	return &RandomBulk{64 * 1024 * 1024}
}

func (rb *RandomBulk) readASCIIBytes(bs []byte) {
	bsl := len(bs)
	for i := 0; i < bsl; i++ {
		bs[i] = latinLetters[rand.Intn(52)]
	}
}

func (rb *RandomBulk) readAnyBytes(bs []byte) {
	bsl := len(bs)
	for i := 0; i < bsl; i++ {
		bs[i] = allBytes[rand.Intn(256)]
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
		readBytes = rb.readASCIIBytes
	} else {
		readBytes = rb.readAnyBytes
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
