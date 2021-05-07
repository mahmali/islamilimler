package servis

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/djimenez/iconv-go"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func KonulariGettir(kitapNo int, konuId int) []int {
	url := fmt.Sprintf("http://islamilimleri.com/Ktphn/Kitablar/05/%03d/Arapca/%02d/000.htm", kitapNo, konuId)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	//charset := detectContentCharset(resp.Body)
	utfBody, err := iconv.NewReader(resp.Body, "windows-1256", "utf-8")
	if err != nil {
		fmt.Println(err)
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
	doc.Find("select[name=CD63] option").Each(func(i int, selection *goquery.Selection) {
		fmt.Println(selection.Text())
		if selection.Text() != "باب" {
			deger, err := strconv.Atoi(selection.Text())
			if err != nil {
				fmt.Println(err)
			}
			m = append(m, deger)
		}
	})
	return m
}
