package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	_ "github.com/djimenez/iconv-go"
	"github.com/muhammedaliakkaya/islamilimler/model"
	"github.com/muhammedaliakkaya/islamilimler/servis"
	"io/ioutil"
	"time"
)

func main() {
	var sayfalar = [][]int{
		{1, 98},
		{2, 57},
		{3, 46},
		{4, 42},
		{5, 52},
		{6, 38},
		{7, 60},
		{9, 24},
		{19, 20},
	}
	var hadisler []model.Hadis
	var baslik string
	var konuId []int
	for i := 0; i < len(sayfalar); i++ {
		babSayisi := sayfalar[i][1]
		kitapNo := sayfalar[i][0]
		baslik = servis.Baslikgettir(kitapNo)
		time.Sleep(1 * time.Second)
		for j := 1; j <= babSayisi; j++ {
			konuId = servis.KonulariGettir(kitapNo, j)
			for _, val := range konuId {
				hadisler = append(hadisler, servis.VerileriCek(j, val, kitapNo, baslik))
			}
		}
	}

	hadislerJson, _ := json.MarshalIndent(&hadisler, "", " ")

	hadislerJson = bytes.Replace(hadislerJson, []byte("\\u003c"), []byte("<"), -1)
	hadislerJson = bytes.Replace(hadislerJson, []byte("\\u003e"), []byte(">"), -1)
	hadislerJson = bytes.Replace(hadislerJson, []byte("\\u0026"), []byte("&"), -1)

	ioutil.WriteFile("ArapcaHadisler.json", hadislerJson, 0644)

	fmt.Println(len(hadisler))
}
