package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/djimenez/iconv-go"
	_ "github.com/djimenez/iconv-go"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func KonulariGettir(konuId int) map[int]int {
	url := fmt.Sprintf("http://islamilimleri.com/Ktphn/Kitablar/05/001/Turkce/%02d/000.htm", konuId)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)

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
	return m
}

func ElemanleriGettir() {
	resp, err := http.Get("http://islamilimleri.com/Ktphn/Kitablar/05/001/Turkce/04/060.htm")
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)

	}

	data, _ := ioutil.ReadAll(resp.Body)
	yeni := make([]byte, len(data)*2)
	iconv.Convert(data, yeni, "windows-1254", "utf-8")
	yepyeni := bytes.TrimFunc(yeni, func(r rune) bool {
		return r == 0
	})
	var sb strings.Builder
	var yepYeniRune = []rune(string(yepyeni))
	pkapanmali := false

	for i, harf := range yepYeniRune {

		if harf == '<' {
			if yepYeniRune[i+1] == 'p' {
				pkapanmali = true
			}
			if yepYeniRune[i+1] == '/' && yepYeniRune[i+2] == 'p' && yepYeniRune[i+3] == '>' {
				pkapanmali = false
			}

			if pkapanmali && yepYeniRune[i+1] == '/' && yepYeniRune[i+2] == 't' && yepYeniRune[i+3] == 'd' && yepYeniRune[i+4] == '>' {
				sb.WriteString("</p>")
				pkapanmali = false
			}
		}

		sb.WriteRune(harf)

	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(sb.String()))
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("#icerik div table tr td[valign=top] p").Each(func(i int, selection *goquery.Selection) {
		fmt.Println(i, "----")
		fmt.Println(selection.Html())
		fmt.Println("---")
	})

	/*metin:=doc.Find("td[valign=top]").Text()
	fmt.Println(metin)
	doc.Find("#icerik").Each(func(i int, selection *goquery.Selection) {
		fmt.Println(selection.Text())
	})
	metin=doc.Find("#icerik").Text()
	fmt.Println(metin)*/
}

func verileriCek(id int, babNo int) Hadis {
	url := fmt.Sprintf("http://islamilimleri.com/Ktphn/Kitablar/05/001/Turkce/%02d/%03d.htm", id, babNo)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)

	}

	data, _ := ioutil.ReadAll(resp.Body)
	yeni := make([]byte, len(data)*2)
	iconv.Convert(data, yeni, "windows-1254", "utf-8")
	yepyeni := bytes.TrimFunc(yeni, func(r rune) bool {
		return r == 0
	})
	var sb strings.Builder
	var yepYeniRune = []rune(string(yepyeni))
	pkapanmali := false

	for i, harf := range yepYeniRune {

		if harf == '<' {
			if yepYeniRune[i+1] == 'p' {
				pkapanmali = true
			}
			if yepYeniRune[i+1] == '/' && yepYeniRune[i+2] == 'p' && yepYeniRune[i+3] == '>' {
				pkapanmali = false
			}

			if pkapanmali && yepYeniRune[i+1] == '/' && yepYeniRune[i+2] == 't' && yepYeniRune[i+3] == 'd' && yepYeniRune[i+4] == '>' {
				sb.WriteString("</p>")
				pkapanmali = false
			}
		}

		sb.WriteRune(harf)

	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(sb.String()))
	if err != nil {
		log.Fatal(err)
	}

	BaslıkBul := func(name string) string {
		secim := ""
		doc.Find("select[name=" + name + "] option").Each(func(i int, selection *goquery.Selection) {
			if _, varMi := selection.Attr("selected"); varMi {
				secim = selection.Text()
			}
		})
		return secim
	}
	Baslık := BaslıkBul("CD71")
	var metin []string
	var ilkHtml string
	var html []string

	doc.Find("td[valign=top] b ").Each(func(i int, selection *goquery.Selection) {
		ilkHtml, _ = selection.Html()
		ilkHtml = strings.ReplaceAll(ilkHtml, "\u003c", "<")
		ilkHtml = strings.ReplaceAll(ilkHtml, "\u003e", ">")
		ilkHtml = strings.ReplaceAll(ilkHtml, "\u0026#34", "")
		ilkHtml = strings.ReplaceAll(ilkHtml, "\n", "")
		ilkHtml = strings.ReplaceAll(ilkHtml, "\t", "")
		html = append(html, ilkHtml)
	})
	doc.Find("td[valign=top] Br ").Each(func(i int, selection *goquery.Selection) {
		ilkHtml, _ = selection.Html()
		ilkHtml = strings.ReplaceAll(ilkHtml, "\\u0026#34", "")
		ilkHtml = strings.ReplaceAll(ilkHtml, "\n", "")
		ilkHtml = strings.ReplaceAll(ilkHtml, "\t", "")
		html = append(html, ilkHtml)
	})
	doc.Find("#icerik div table tr td[valign=top] p").Each(func(i int, selection *goquery.Selection) {
		MetinTrim := strings.TrimSpace(selection.Text())
		MetinTrim = strings.ReplaceAll(MetinTrim, "\n", "")
		MetinTrim = strings.ReplaceAll(MetinTrim, "\t", "")
		MetinTrim = strings.ReplaceAll(MetinTrim, "\"", "")
		metin = append(metin, MetinTrim)

	})

	hadis := Hadis{
		Konu:    Baslık,
		Numara:  babNo,
		Metin:   metin,
		HtmlTag: html,
	}
	return hadis
}

/*func SayfaninMetniniGetir(konuId int) []Hadis {

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
			Numara: deger.Numara,
			Metin:  icerik,
		})
	}
	return Hadisler
}*/

func main() {
	hadisler := make([]Hadis, 0)
	for i := 1; i < 99; i++ {
		konuId := KonulariGettir(i)
		fmt.Println(konuId)
		for _, val := range konuId {
			hadisler = append(hadisler, verileriCek(i, val))
		}
	}

	if data, err := json.MarshalIndent(hadisler, " ", " "); err != nil {
		log.Fatal(err)
	} else {
		ioutil.WriteFile("hadisler.json", data, 0644)
	}

	fmt.Println(len(hadisler))
}
