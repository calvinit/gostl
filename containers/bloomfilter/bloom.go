package bloom

import (
	"bytes"
	"encoding/binary"
	"github.com/liyue201/gostl/algorithm/hash"
	"github.com/liyue201/gostl/containers/bitmap"
	"math"
)

const Salt = "g9hmj2fhgr"

type BloomFilter struct {
	m uint64
	k uint64
	b *bitmap.Bitmap
}

// New new a BloomFilter with m bits and k hash functions
func New(m, k uint64) *BloomFilter {
	return &BloomFilter{
		m: m,
		k: k,
		b: bitmap.New(m),
	}
}

// New new a BloomFilter with n and fp.
// n is the capacity of the BloomFilter
// fp is the tolerated error rate of the BloomFilter
func NewWithEstimates(n uint64, fp float64) *BloomFilter {
	m, k := EstimateParameters(n, fp)
	return New(m, k)
}

//NewFromData new a BloomFilter by data passed, the data was generated by function 'Data()'
func NewFromData(data []byte) *BloomFilter {
	b := &BloomFilter{}
	reader := bytes.NewReader(data)
	binary.Read(reader, binary.LittleEndian, &b.m)
	binary.Read(reader, binary.LittleEndian, &b.k)
	b.b = bitmap.NewFromData(data[8+8:])
	return b
}

func EstimateParameters(n uint64, p float64) (m uint64, k uint64) {
	m = uint64(math.Ceil(-1 * float64(n) * math.Log(p) / (math.Ln2 * math.Ln2)))
	k = uint64(math.Ceil(math.Ln2 * float64(m) / float64(n)))
	return
}

// Add add a value to the BloomFilter
func (this *BloomFilter) Add(val string) {
	hashs := hash.GenHashInts([]byte(Salt+val), int(this.k))
	for i := uint64(0); i < this.k; i++ {
		this.b.Set(hashs[i] % this.m)
	}
}

// Contains returns true if value passed is (high probability) in the BloomFilter, or false if not.
func (this *BloomFilter) Contains(val string) bool {
	hashs := hash.GenHashInts([]byte(Salt+val), int(this.k))
	for i := uint64(0); i < this.k; i++ {
		if !this.b.IsSet(hashs[i] % this.m) {
			return false
		}
	}
	return true
}

// Contains returns the data of BloomFilter, it can bee used to new a BloomFilter by using function 'NewFromData' .
func (this *BloomFilter) Data() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, this.m)
	binary.Write(buf, binary.LittleEndian, this.k)
	buf.Write(this.b.Data())
	return buf.Bytes()
}