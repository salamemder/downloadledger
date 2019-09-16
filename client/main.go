package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

const demodata = "https://www.monash.edu/study"

const SERVERURL="http://127.0.0.1:3000/query"


func GetDatafromServer() []byte{

	URLDATA := SERVERURL+"?"+"url="+demodata
	resp, err := http.Get(URLDATA)
	if err == nil{
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
	encryptedkeyServer := GetDatafromServer()
	log.Println(encryptedkeyServer)



}
