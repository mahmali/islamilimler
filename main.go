package main

import (
	"encoding/json"
	_ "encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func SayfaninMetniniGetir(konuId int, sayfaNu int) Hadis {

	url := fmt.Sprintf("http://islamilimleri.com/Ktphn/Kitablar/05/001/Turkce/%02d/%03d.htm", konuId, sayfaNu)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	icerik, err := doc.Find("#icerik").Html()

	secimiBul := func(name string) string {
		secim := ""
		doc.Find("select[name=" + name + "] option").Each(func(i int, selection *goquery.Selection) {
			if _, varMi := selection.Attr("selected"); varMi {
				secim = selection.Text()
			}
		})
		return secim
	}
	konu := secimiBul("CD71")
	strhadis := secimiBul("CD39")
	no, _ := strconv.Atoi(strhadis)

	return Hadis{
		Konu:   konu,
		Numara: no,
		Metin:  icerik,
	}
}

func main() {
	hadisler := make([]Hadis, 0)
	for hno := 1; hno < 10; hno++ {
		hadis := SayfaninMetniniGetir(2, hno)
		hadisler = append(hadisler, hadis)
	}
	if data, err := json.MarshalIndent(hadisler, " ", " "); err != nil {
		log.Fatal(err)
	} else {

		ioutil.WriteFile("hadisler.json", data, 0644)
	}

	fmt.Println(len(hadisler))
}
