package main

import (
	"bytes"
	"crypto/sha256"
	"download/cryptoopt"
	"download/garbledbloomfilter"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
)

const testfilekey = "mykey44444444444"      //this is for encrypt the file
const testmasterkey = "mykey23jdlkdleda"    //this is for encrypt the random sequence
const demodata = "https://www.monash.edu/study"
const seed = 324
const DEFAULTDownload = 100
const SERVERURL="http://127.0.0.1:3000/secretkey"
const DEMOURL="http://test.test.com/secretkey"

type FilterStruct struct {
	URL string  `json:"URL"`
	Filter []byte `json:Filter`
	Positionarray [][]uint `Positonarray`
}

type SentJson struct{

	Args []string `json:"Args"`
}

func Gen_sk_CSK_x(downloadcount *int, SK []byte)([]string, []string){
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


		skx_array[i]=base64.StdEncoding.EncodeToString([]byte(encodedstring))
		encbs := aescrypto.Stringtoaeskey(string(bs))

		//we ensure the key size as 128

		var C_sk_x string
		if i < *downloadcount {
			C_sk_x, err = aescrypto.Encrypt(SK, encbs)
		}else{
			C_sk_x, err = aescrypto.Encrypt([]byte("ABORT"), encbs)
		}
		CSK_x_Array[i] = C_sk_x
	}
	return skx_array, CSK_x_Array

}


func main(){

	var downloadcount = flag.Int("downloadcount", 10, "The download count")
	var filtersize = flag.Uint("Filter size", 1000, "set the filter size")
	var k = flag.Uint("Number of hash functions", 5,"number of the hash function")
	flag.Parse()


	SK := []byte(testfilekey)
	filter := garbledbloomfilter.New(20*(*filtersize), *k) // filtersize, 5 keys by default
	_, err := aescrypto.Encrypt([]byte(demodata), SK)
	if err != nil{
		fmt.Println("error in encrypt the file")
		return
	}

	decrypoolkey,CSK_x_Array :=Gen_sk_CSK_x(downloadcount, SK)
	positionsforeachcount := make([][]uint, len(CSK_x_Array))
	i := 0
	for _, each := range CSK_x_Array {
		_,locationsarray, err := filter.Add([]byte(each), i)
		if err != nil{
			panic("creating the bloom filter panic")
		}
		positionsforeachcount[i] = locationsarray
		i += 1
	}

	//dotset(decrypoolkey, CSK_x_Array, filter, positionsforeachcount)

	exportfilter, err := filter.Export()
	if err != nil{
		fmt.Println("error in export the bloom filter")
	}

	uploadedfilter := FilterStruct{
		DEMOURL,
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

	sendledger,_  := json.Marshal(decrypoolkey)
	sendledgerbase64 := base64.StdEncoding.EncodeToString(sendledger)



	url := uploadedfilter.URL

	h := sha256.New()
	h.Write([]byte(url))
	outhash  := fmt.Sprintf("%x", h.Sum(nil))

	sent := make([]string,3)
	sent[0]= "upload"
	sent[1] = outhash
	sent[2] = sendledgerbase64

	sentpack :=SentJson{
		sent,
	}

	sentjson,_ := json.Marshal(sentpack)

	cmd := exec.Command("peer", "chaincode","invoke", "-n", "mycc", "-c", string(sentjson), "-C", "myc")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}

	resultstr := stderr.String()

	formatedstr := strings.Split(resultstr,`\n`)
	targetstr := formatedstr[len(formatedstr)-1]
	retflag := strings.Index(targetstr, "successful")
	if retflag > 0{
		fmt.Println(targetstr)
	}else{
		fmt.Println("upload failed")
	}


	return
}

func dotset(decrypoolkey []string,CSK_x_Array []string,filter *garbledbloomfilter.GarbledBloomFilter, positionsforeachcount [][]uint ){
	for i:=0;i<len(decrypoolkey);i++ {

		//log.Println(decrypoolkey[i])
		//log.Println([]byte(CSK_x_Array[i]))
		encryptedkeyServer:=CSK_x_Array[i]

		location := positionsforeachcount[i]

		retdata, _ := filter.GetByCnt(location)

		if encryptedkeyServer == string(retdata){
			fmt.Println("ok")
		}else{
			fmt.Println("error in ", i)
			fmt.Println(len([]byte(encryptedkeyServer)))
			fmt.Println([]byte(encryptedkeyServer))
			fmt.Println(retdata)
		}
}
}


func stringtobigintexample() {
	t1 := "119 78 112 55 100 56 51 55 99 43 83 66 112 85 103 73 66 82 81 112 111 85 111 122 88 90 52 104 98 98 78 99 109 81 68 73 114 99 119 65 105 73 52 61"

	t1array := strings.Split(t1, " ")

	//test := big.NewInt(32323)

	digit := arraytodig(t1array)

	r1, _ := new(big.Int).SetString("677140963341764673917977438998999713218", 10)
	r2, _ := new(big.Int).SetString("587821017083724015184665925913920905766", 10)
	r3, _ := new(big.Int).SetString("1171503186354600147576549943126075296647", 10)
	r4, _ := new(big.Int).SetString("325197258411314705963805618973758782760", 10)

	out1 := new(big.Int).Xor(digit, r1)
	out2 := new(big.Int).Xor(out1, r2)
	out3 := new(big.Int).Xor(out2, r3)
	out4 := new(big.Int).Xor(out3, r4)

	final := new(big.Int).Xor(out4, r1)
	final = new(big.Int).Xor(final, r2)
	final = new(big.Int).Xor(final, r3)
	final = new(big.Int).Xor(final, r4)
	fmt.Println(final)
	fmt.Println(digit)
}

func arraytodig(array []string) *big.Int{

	testnum := make([]uint8, len(array))
	for i,val := range array{

		tesm,_ := strconv.Atoi(val)
		testnum[i] = uint8(tesm)

	}

	return new(big.Int).SetBytes(testnum)


}