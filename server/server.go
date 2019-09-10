package main

import (
	"flag"
	"fmt"
	"math/rand"
	"yanjunshen/cryptoopt"
	"yanjunshen/garbledbloomfilter"
)


const testmasterkey = "mykey23jdlkdleda"
const demourl = "https://www.monash.edu/study"
const seed = 324
func main(){


	var n uint
	n = 1000
	filter := garbledbloomfilter.New(20*n, 5) // load of 20, 5 keys
	filter.Add([]byte("hello"))
	value, err := filter.Get([]byte("hello"))
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println(string(value))



	var downloadcount = flag.Int("downloadcount", 100, "The download count")

	flag.Parse()
	fmt.Println(*downloadcount)

	rand.Seed(seed)
	for i:=0;i<10;i++{
		fmt.Println(rand.Int())
	}


	output, err := aescrypto.Encrypt([]byte(demourl),[]byte(testmasterkey))
	if err != nil{
		fmt.Println(err)
		return
	}
	fmt.Println(output)

	rawhello, err := aescrypto.Decrypt(output, []byte(testmasterkey))
	fmt.Println(rawhello)

}
