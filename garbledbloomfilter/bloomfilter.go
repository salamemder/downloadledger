package garbledbloomfilter

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/spaolacci/murmur3"
	"math/big"
)

const SIZE = 128

type GarbledBloomFilter struct {
	m uint
	k uint
	b [][]byte
}

type ExportedFilter struct {
	M uint     `json:"M"`
	K uint     `json:"K"`
	B []string `json:Filter`
}

// New creates a new Bloom filter with _m_ bits and _k_ hashing functions
// We force _m_ and _k_ to be at least one to avoid panics.
func New(m uint, k uint) *GarbledBloomFilter {
	storage := make([][]byte, m)
	return &GarbledBloomFilter{max(1, m), max(1, k), storage}
}

// baseHashes returns the four hash values of data that are used to create k
// hashes
func baseHashes(data []byte) [4]uint64 {
	a1 := []byte{1} // to grab another bit of data
	hasher := murmur3.New128()
	hasher.Write(data) // #nosec
	v1, v2 := hasher.Sum128()
	hasher.Write(a1) // #nosec
	v3, v4 := hasher.Sum128()
	return [4]uint64{
		v1, v2, v3, v4,
	}
}

// location returns the ith hashed location using the four base hash values
func location(h [4]uint64, i uint) uint64 {
	ii := uint64(i)
	return h[ii%2] + ii*h[2+(((ii+(ii%2))%4)/2)]
}

// location returns the ith hashed location using the four base hash values
func (f *GarbledBloomFilter) location(h [4]uint64, i uint) uint {
	return uint(location(h, i) % uint64(f.m))
}

func (f *GarbledBloomFilter) Export() ([]byte, error) {

	k := f.k
	m := f.m
	stringarray := make([]string, m)
	for i := uint(0); i < m; i++ {
		str := base64.StdEncoding.EncodeToString(f.b[i])
		stringarray[i] = str
	}

	res1D := &ExportedFilter{
		m,
		k,
		stringarray,
	}
	ret, err := json.Marshal(res1D)
	if err != nil{
		return nil, err
	}
	return ret, nil
}

// Add data to the Bloom Filter. Returns the filter (allows chaining)
func (f *GarbledBloomFilter) Add(data []byte) (*GarbledBloomFilter, error) {
	localtionmap := make(map[uint]*big.Int)
	empty := 0
	h := baseHashes(data)
	max := new(big.Int)
	max.Exp(big.NewInt(2), big.NewInt(130), nil).Sub(max, big.NewInt(1))
	lastelement := new(big.Int).SetBytes(data)
	for i := uint(0); i < f.k; i++ {
		location := f.location(h, i)
		if len(f.b[location]) == 0 {
			empty += 1
			localtionmap[location] = big.NewInt(-1)

		} else {
			localtionmap[location] = new(big.Int).SetBytes(f.b[location])

		}
	}
	if empty == 0 {
		return nil, fmt.Errorf("no position for this key")
	}
	pos := uint(0)
	if empty == 1 {
		for key, val := range (localtionmap) {
			if val == big.NewInt(-1) {
				pos = key
			} else {
				lastelement = new(big.Int).Xor(val, lastelement)
			}
		}
		localtionmap[pos] = lastelement

	} else {
		handled := false
		for key, val := range (localtionmap) {
			if val.Cmp(big.NewInt(-1)) == 0 && handled == false {
				pos = key
				handled = true
			} else {
				n, err := rand.Int(rand.Reader, max)
				if err != nil {
					panic("error in put into the bloom filter")
				}
				localtionmap[key] = n
				lastelement = lastelement.Xor(lastelement, n)
			}
		}
		localtionmap[pos] = lastelement
	}

	//verify
	//org:= big.NewInt(0)
	//for _,val := range(localtionmap){
	//	//fmt.Println(val)
	//	org = org.Xor(org, val)
	//}
	//
	//fmt.Println(org.Bytes())
	//fmt.Println("org", data)

	for key, value := range localtionmap {
		f.b[key] = value.Bytes()
	}
	return f, nil
}

// Add data to the Bloom Filter. Returns the filter (allows chaining)
func (f *GarbledBloomFilter) Get(data []byte) ([]byte, error) {
	h := baseHashes(data)

	location := f.location(h, uint(0))
	retval := big.NewInt(0)

	for i := uint(0); i < f.k; i++ {
		location = f.location(h, i)
		storedval := f.b[location]
		if len(storedval) == 0 {
			return nil, fmt.Errorf("no such key")
		}
		thisval := new(big.Int).SetBytes(storedval)

		if thisval.Cmp(big.NewInt(-1)) == 0 {

		}
		retval = new(big.Int).Xor(retval, thisval)
	}
	return retval.Bytes(), nil
}

func max(x, y uint) uint {
	if x > y {
		return x
	}
	return y
}
