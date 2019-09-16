package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"net/http"
	"yanjunshen/cryptoopt"
	"yanjunshen/garbledbloomfilter"
)

const testfilekey = "mykey44444444444"
const testmasterkey = "mykey23jdlkdleda"
const demodata = "https://www.monash.edu/study"
const seed = 324
const DEFAULTDownload = 100
const SERVERURL="http://127.0.0.1:3000/secretkey"


type FilterStruct struct {
	URL string  `json:"URL"`
	Filter []byte `json:Filter`
	Positionarray [][]uint `Positonarray`
}

func Gen_sk_CSK_x(downloadcount *int, SK []byte)([]uint64, []string){
	rand.Seed(seed)
	x_hat_array := make([]uint64,DEFAULTDownload)
	skx_array := make([]string, DEFAULTDownload)
	CSK_x_Array:= make([]string, DEFAULTDownload)
	max := new(big.Int)
	max.Exp(big.NewInt(2), big.NewInt(130), nil).Sub(max, big.NewInt(1))
	for i:=0;i<DEFAULTDownload;i++{
		x_hat_array[i] = rand.Uint64()
		bs := make([]byte, 8)
		binary.LittleEndian.PutUint64(bs, x_hat_array[i])
		encodedstring, err := aescrypto.Encrypt(bs, []byte(testmasterkey))
		if err != nil{
			fmt.Println("encode the conter failed!")
		}
		skx_array[i]=encodedstring
		sk_x := aescrypto.Stringtoaeskey(encodedstring)
		//we ensure the key size as 128
		if len(sk_x) != 16{
			panic("incorrect sk_x key size")
		}
		var C_sk_x string
		if i < *downloadcount {
			C_sk_x, err = aescrypto.Encrypt(SK, sk_x)
		}else{
			C_sk_x, err = aescrypto.Encrypt([]byte("ABORT"), sk_x)
		}
		CSK_x_Array[i] = C_sk_x
	}
	return x_hat_array, CSK_x_Array

}


func main(){

	var downloadcount = flag.Int("downloadcount", 10, "The download count")
	var filtersize = flag.Uint("Filter size", 1000, "set the filter size")
	var k = flag.Uint("Number of hash functions", 5,"number of the hash function")
	flag.Parse()


	SK := []byte(testfilekey)
	filter := garbledbloomfilter.New(20*(*filtersize), *k) // load of 20, 5 keys
	_, err := aescrypto.Encrypt([]byte(demodata), SK)
	if err != nil{
		fmt.Println("error in encrypt the file")
		return
	}

	_,CSK_x_Array :=Gen_sk_CSK_x(downloadcount, SK)
	positionsforeachcount := make([][]uint, len(CSK_x_Array))
	i := 0
	for _, each := range CSK_x_Array {
		_,locationsarray, err := filter.Add([]byte(each))
		if err != nil{
			panic("creating the bloom filter panic")
		}
		positionsforeachcount[i] = locationsarray
		i += 1
		break
	}
	log.Println([]byte(CSK_x_Array[0]))
	exportfilter, err := filter.Export()
	if err != nil{
		fmt.Println("error in export the bloom filter")
	}

	uploadedfilter := FilterStruct{
		demodata,
		exportfilter,
		positionsforeachcount,
	}

	uploadedbytes, err := json.Marshal(uploadedfilter)

	resp, err := http.Post(SERVERURL, "application/json", bytes.NewBuffer(uploadedbytes))

	if err != nil{
		fmt.Println(err)
	}

	if resp.StatusCode == 200{
		fmt.Println("upload the bloom filter to the server successfully")
	}
	return
}
