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
	"time"
)

func Baslikgettir(kitapNo int) string {

	url := fmt.Sprintf("http://islamilimleri.com/Ktphn/Kitablar/05/%03d/Turkce/01/000.htm", kitapNo)
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
	var baslık string
	doc.Find("[valign=top] h1").Each(func(i int, selection *goquery.Selection) {
		baslık = selection.Text()
	})
	baslık = strings.ReplaceAll(baslık, "\n", "")
	baslık = strings.ReplaceAll(baslık, "\t", "")
	baslık = strings.ReplaceAll(baslık, "  ", "")
	baslık = strings.TrimSpace(baslık)
	fmt.Println(baslık)

	return baslık
}

func KonulariGettir(kitapNo int, konuId int) []int {
	url := fmt.Sprintf("http://islamilimleri.com/Ktphn/Kitablar/05/%03d/Turkce/%02d/000.htm", kitapNo, konuId)
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
	var baslık string
	doc.Find("[valign=top] h1").Each(func(i int, selection *goquery.Selection) {
		baslık = selection.Text()
	})
	baslık = strings.ReplaceAll(baslık, "\n", "")
	baslık = strings.ReplaceAll(baslık, "\t", "")
	fmt.Println(baslık)
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

func verileriCek(id int, babNo int, sayfa int, kitapIsmi string) Hadis {
	url := fmt.Sprintf("http://islamilimleri.com/Ktphn/Kitablar/05/%03d/Turkce/%02d/%03d.htm", sayfa, id, babNo)
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
			MetinTrim = strings.ReplaceAll(MetinTrim, "  ", " ")
			MetinTrim = strings.ReplaceAll(MetinTrim, "(=", "(")
			MetinTrim = strings.ReplaceAll(MetinTrim, ">=", ">")
			MetinTrim = strings.ReplaceAll(MetinTrim, "( =", "(")
			//MetinTrim = strings.ReplaceAll(MetinTrim, ",", "")
			MetinTrim = strings.ReplaceAll(MetinTrim, "„", "")
			MetinTrim = strings.ReplaceAll(MetinTrim, "|b|", "<b>")
			MetinTrim = strings.ReplaceAll(MetinTrim, "|/b|", "</b>")
			MetinTrim = strings.ReplaceAll(MetinTrim, "<b></b>", " ")
			metin = append(metin, MetinTrim)
		})

	}
	metinGetir("h3")
	metinGetir("p")
	hadis := Hadis{
		Kitap:  kitapIsmi,
		Konu:   Baslık,
		Numara: babNo,
		Metin:  metin,
	}
	return hadis
}

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
	var hadisler []Hadis
	var baslik string
	var konuId []int
	for i := 0; i < len(sayfalar); i++ {
		babSayisi := sayfalar[i][1]
		kitapNo := sayfalar[i][0]
		baslik = Baslikgettir(kitapNo)
		time.Sleep(1 * time.Second)
		for j := 1; j <= babSayisi; j++ {
			konuId = KonulariGettir(kitapNo, j)
			for _, val := range konuId {
				hadisler = append(hadisler, verileriCek(j, val, kitapNo, baslik))
			}
		}
	}

	hadislerJson, _ := json.MarshalIndent(&hadisler, "", " ")

	hadislerJson = bytes.Replace(hadislerJson, []byte("\\u003c"), []byte("<"), -1)
	hadislerJson = bytes.Replace(hadislerJson, []byte("\\u003e"), []byte(">"), -1)
	hadislerJson = bytes.Replace(hadislerJson, []byte("\\u0026"), []byte("&"), -1)

	ioutil.WriteFile("hadisler.json", hadislerJson, 0644)

	fmt.Println(len(hadisler))
}
