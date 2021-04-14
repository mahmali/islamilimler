package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/djimenez/iconv-go"
	_ "github.com/djimenez/iconv-go"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
)

func KonulariGettir(konuId int) map[int]int {
	url := fmt.Sprintf("http://islamilimleri.com/Ktphn/Kitablar/05/001/Turkce/%02d/000.htm", konuId)

	fmt.Println(url)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)
		return map[int]int{}
	}

	//charset := detectContentCharset(resp.Body)
	utfBody, err := iconv.NewReader(resp.Body, "windows-1254", "utf-8")
	if err != nil {
		// handler error
	}

	doc, err := goquery.NewDocumentFromReader(utfBody)
	if err != nil {
		log.Fatal(err)
	}
	m := make(map[int]int)
	doc.Find("select[name=CD39] option").Each(func(i int, selection *goquery.Selection) {
		fmt.Println(selection.Text())
		if selection.Text() != "BAB" {

			deger, err := strconv.Atoi(selection.Text())
			if err != nil {
				fmt.Println(err)
			}
			m[i] = deger
		}
	})
	//select listeki tüm o 001 002 vs geziyok
	for val, _ := range m {
		verileriCek(fmt.Sprintf("http://islamilimleri.com/Ktphn/Kitablar/05/001/Turkce/%02d/%03d.htm", konuId, val), konuId, val)
	}
	return m
}

func verileriCek(url string, konuId, babNo int) {
	dir, _ := os.Getwd()
	yol := path.Join(dir, "sonuc", strconv.Itoa(konuId))
	os.MkdirAll(yol, 0777)
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)
		return
	}

	utfBody, err := iconv.NewReader(resp.Body, "windows-1254", "utf-8")
	if err != nil {
		fmt.Println(err.Error())
	}

	doc, err := goquery.NewDocumentFromReader(utfBody)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("td[valign=top]").Each(func(i int, selection *goquery.Selection) {
		dosyaAdi := fmt.Sprintf("bab-%d.html", babNo)
		file, _ := os.Create(path.Join(yol, dosyaAdi))
		file.WriteString(selection.Text())
	})
}

func SayfaninMetniniGetir(konuId int) []Hadis {

	url := fmt.Sprintf("http://islamilimleri.com/Ktphn/Kitablar/05/001/Turkce/%02d/000.htm", konuId)

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
	icerik := doc.Find("#icerik").Text()

	secimiBul := func(name string) string {
		secim := ""
		doc.Find("select[name=" + name + "] option").Each(func(i int, selection *goquery.Selection) {
			if _, varMi := selection.Attr("selected"); varMi {
				secim = selection.Text()
			}
		})
		return secim
	}
	var Hadisler []Hadis
	konu := secimiBul("CD71")
	konular := KonulariGettir(konuId)
	for _, deger := range konular {
		Hadisler = append(Hadisler, Hadis{
			Konu:   konu,
			Numara: deger,
			Metin:  icerik,
		})
	}
	return Hadisler
}

func main() {
	hadisler := make([]Hadis, 0)
	for i := 1; i < 98; i++ {
		KonulariGettir(i)
	}
	/*for hno := 1; hno < 98; hno++ {
		Thadis := SayfaninMetniniGetir(hno)
		hadisler = append(hadisler, Thadis...)
	}*/
	if data, err := json.MarshalIndent(hadisler, " ", " "); err != nil {
		log.Fatal(err)
	} else {
		ioutil.WriteFile("hadisler.json", data, 0644)
	}

	fmt.Println(len(hadisler))
}
