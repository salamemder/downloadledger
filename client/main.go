package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	aescrypto "download/cryptoopt"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)
const DEMOURL="http://test.test.com/secretkey"

type SentJson struct{

	Args []string `json:"Args"`
}


const SERVERURL="http://127.0.0.1:3000/query"
const testmasterkey = "mykey23jdlkdleda"    //this is for encrypt the random sequence


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


func querybloclchain(url string) string{

	h := sha256.New()
	h.Write([]byte(url))
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
	decodekey  = strings.TrimSuffix(decodekey, "\n")

	lengthstring := len(decodekey)-2
	decodekey = decodekey[1:lengthstring]
	return decodekey

}

func ShowMsg(){

	fmt.Printf("*********************\n")
	fmt.Println("1: test key decryption.\n2: test concurrency.\n3: quit program\n")
	fmt.Printf("*********************\n")


}

func Testkeydecryption() string{


	decodekey := querybloclchain(DEMOURL)





	return decodekey

}

func main(){


	previouskey:= ""
        counter :=0
	for {
		ShowMsg()
		fmt.Print("")
		reader := bufio.NewReader(os.Stdin)
		input,_ := reader.ReadString('\n')

		switch input {
		case "1\n":
			for {
				decodekey := querybloclchain(DEMOURL)
				if decodekey == previouskey {
					time.Sleep(time.Second)
					continue
				} else {
					previouskey = decodekey
					encryptedkeyServer := GetDatafromServer(DEMOURL)

					decoded, err := base64.StdEncoding.DecodeString(decodekey)

					if err != nil {
						fmt.Println(err)
					}

					encodedstring, err := aescrypto.Decrypt(string(decoded), []byte(testmasterkey))
					if err != nil {
						fmt.Println(err)
					}

					sk_x := aescrypto.Stringtoaeskey(encodedstring)

					decrypted, err := aescrypto.Decrypt(string(encryptedkeyServer), sk_x)
					counter += 1
					if err != nil {
						fmt.Println("Failed to decrypt data, inconsistent download count")
					}

					if counter == 2{
					    fmt.Println("Fail to decrypt data, inconsistent download count")
	 				    break
					}

					if decrypted == "ABORT"{
					    fmt.Println("Execeed download limit")
					    break
					}
					
					if decrypted == "This is the Demo"{

					    fmt.Println("This is the Demo data")
					}else{
					    fmt.Printf("Decrypted data error")
					}
					break
				}
			}
		case "2\n":

			var timestart time.Time
			var pre string
			for {
				decodekey := querybloclchain(DEMOURL)
				if decodekey == previouskey{
					time.Sleep(time.Second)
				}else{
					fmt.Printf("get data from server\n")
					GetDatafromServer(DEMOURL)
					pre = decodekey
					timestart = time.Now()
					break
				}

			}
			var delta time.Duration
			for {
				decodekey := querybloclchain(DEMOURL)
				if decodekey == pre{
					time.Sleep(time.Millisecond*10)
				}else{
					fmt.Printf("get data from server\n")

					GetDatafromServer(DEMOURL)
					delta = time.Since(timestart)
					break
				}

			}
			fmt.Println(delta, "with +-10 ms")


		case "3\n":
			break
		default:
			fmt.Println("invalid input")
			continue
		}
		if input == "3\n"{
			break
		}

	}
	return






}
