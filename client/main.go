package main

import (
	"bytes"
	"crypto/sha256"
	aescrypto "download/cryptoopt"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"
)
const DEMOURL="http://test.test.com/secretkey"

type SentJson struct{

	Args []string `json:"Args"`
}


const SERVERURL="http://127.0.0.1:3000/query"


func GetDatafromServer(resourceurl string) []byte{

	URLDATA := SERVERURL+"?"+"url="+resourceurl
	resp, err := http.Get(URLDATA)
	if err != nil{
		log.Println("error in query the server")
	}
	if resp.StatusCode != 202{
		log.Println("error in processing the request")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body

}


func main(){
	encryptedkeyServer := GetDatafromServer(DEMOURL)

	h := sha256.New()
	h.Write([]byte(DEMOURL))
	outhash  := fmt.Sprintf("%x", h.Sum(nil))

	sent := make([]string,2)
	sent[0]= "downloadquery"
	sent[1] = outhash

	sentpack :=SentJson{
		sent,
	}

	sentjson,_ := json.Marshal(sentpack)

	cmd := exec.Command("peer", "chaincode","invoke", "-n", "mycc", "-c", string(sentjson), "-C", "myc")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}

	resultstr := stderr.String()

	formatedstr := strings.Split(resultstr,`\n`)
	targetstr := formatedstr[len(formatedstr)-1]
	finalstr := strings.SplitAfter(targetstr, "payload:")

	decodekey := finalstr[len(finalstr)-1]
	decodekeybyte := []byte(decodekey)
	length := len(decodekeybyte)
	decodekey = string(decodekeybyte[1:length-3])

	sk_x := aescrypto.Stringtoaeskey(decodekey)


	fmt.Println(encryptedkeyServer)
	fmt.Println(decodekey)


	decrypeted,err := aescrypto.Decrypt(string(encryptedkeyServer),sk_x)

	if err != nil{
		fmt.Println(err)
	}

	fmt.Printf("mast key %s\n",  decrypeted)

	//log.Printf("Command finished with error: %v", output)
	//log.Printf("Command finished with error: %v", err)


}
