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

func KonulariGettir(konuId int) []int {
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
	var m []int
	doc.Find("select[name=CD39] option").Each(func(i int, selection *goquery.Selection) {
		fmt.Println(selection.Text())
		if selection.Text() != "BAB" {
			deger, err := strconv.Atoi(selection.Text())
			if err != nil {
				fmt.Println(err)
			}
			m = append(m, deger)
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
	fmt.Println(url)
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
			if yepYeniRune[i+1] == 'b' {
				sb.WriteString("|b|")
			}
			if yepYeniRune[i+1] == '/' && yepYeniRune[i+2] == 'b' && yepYeniRune[i+3] == '>' {
				sb.WriteString("|/b|")
			}
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

	metinGetir := func(aranacak string) {
		doc.Find(fmt.Sprintf("#icerik div table tr td[valign=top] %s", aranacak)).Each(func(i int, selection *goquery.Selection) {
			hammetin := selection.Text()
			MetinTrim := strings.TrimSpace(hammetin)
			MetinTrim = strings.ReplaceAll(MetinTrim, "\n", "")
			MetinTrim = strings.ReplaceAll(MetinTrim, "\t", "")
			MetinTrim = strings.ReplaceAll(MetinTrim, "...", "")
			MetinTrim = strings.ReplaceAll(MetinTrim, "…", "")
			MetinTrim = strings.ReplaceAll(MetinTrim, ">..", "")
			MetinTrim = strings.ReplaceAll(MetinTrim, "..", "...")
			MetinTrim = strings.ReplaceAll(MetinTrim, "‏", "")
			MetinTrim = strings.ReplaceAll(MetinTrim, "­", "")
			MetinTrim = strings.ReplaceAll(MetinTrim, " ", "")
			MetinTrim = strings.ReplaceAll(MetinTrim, "''", "\"")
			MetinTrim = strings.ReplaceAll(MetinTrim, "_", "")
			MetinTrim = strings.ReplaceAll(MetinTrim, "        ", " ")
			MetinTrim = strings.ReplaceAll(MetinTrim, "(=", "(")
			MetinTrim = strings.ReplaceAll(MetinTrim, ">=", ">")
			MetinTrim = strings.ReplaceAll(MetinTrim, "( =", "(")
			MetinTrim = strings.ReplaceAll(MetinTrim, ",", "")
			MetinTrim = strings.ReplaceAll(MetinTrim, "„", "")
			MetinTrim = strings.ReplaceAll(MetinTrim, "|b|", "<b>")
			MetinTrim = strings.ReplaceAll(MetinTrim, "|/b|", "</b>")
			metin = append(metin, MetinTrim)
		})

	}

	metinGetir("h3")
	metinGetir("p")
	hadis := Hadis{
		Kitap:  "SAHÎH-İ BUHÂRÎ",
		Konu:   Baslık,
		Numara: babNo,
		Metin:  metin,
	}
	return hadis
}

func main() {
	var hadisler []Hadis
	for i := 1; i < 99; i++ {
		konuId := KonulariGettir(i)
		for _, val := range konuId {
			hadisler = append(hadisler, verileriCek(i, val))
		}
	}
	hadislerJson, _ := json.MarshalIndent(&hadisler, "", " ")

	hadislerJson = bytes.Replace(hadislerJson, []byte("\\u003c"), []byte("<"), -1)
	hadislerJson = bytes.Replace(hadislerJson, []byte("\\u003e"), []byte(">"), -1)
	hadislerJson = bytes.Replace(hadislerJson, []byte("\\u0026"), []byte("&"), -1)

	ioutil.WriteFile("hadisler.json", hadislerJson, 0644)

	fmt.Println(len(hadisler))
}
