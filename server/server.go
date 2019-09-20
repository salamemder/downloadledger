package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-martini/martini"
	"io/ioutil"
	"log"
	"net/http"
	"yanjunshen/garbledbloomfilter"
)

type FilterStruct struct {
	URL string  `json:"URL"`
	Filter []byte `json:Filter`
	Positionarray [][]uint `Positonarray`
}

type Datapack struct{
	filter garbledbloomfilter.GarbledBloomFilter
	Positionarray [][]uint
	counter uint
}

var Filterdic map[string]Datapack


func init(){
	log.SetPrefix("Server")
}

// AddAlbum creates the posted album.
func AddFilter(r *http.Request) (int) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	var filter garbledbloomfilter.ExportedFilter
	var recevedfilter  FilterStruct

	err = json.Unmarshal(body, &recevedfilter)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(recevedfilter.Filter, &filter)
	if err != nil {
		panic(err)
	}


	url := recevedfilter.URL
	positionlist := recevedfilter.Positionarray

	bloomfilter, err := garbledbloomfilter.Import(filter)
	if err != nil{
		return http.StatusBadRequest
	}

	data := Datapack{
		*bloomfilter,
		positionlist,
		uint(0),
	}

	Filterdic[url] = data

	return http.StatusAccepted
}

func QueryDowload(r *http.Request) (int, []byte) {
	qs := r.URL.Query()
	url := qs.Get("url")

	if val, ok := Filterdic[url]; ok {
		counter := val.counter
		retdata, err := val.filter.GetByCnt(val.Positionarray[counter])
		log.Println(retdata)
		if err != nil{
			return http.StatusBadRequest, nil
		}else{
			return http.StatusAccepted, retdata
		}
		val.counter +=1
		fmt.Println(val)
		Filterdic[url] = val
	}else{
		return http.StatusBadRequest,nil
	}


	return http.StatusAccepted, nil
}


func main() {

	Filterdic = make(map[string]Datapack)


	m := martini.Classic()
	m.Get("/", func() string {
		return "Hello world!"
	})

	m.Get("/query", QueryDowload)

	m.Post("/secretkey", AddFilter)
	m.Run()
}