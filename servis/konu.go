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
			deger := LatinceKaraktereCevir(selection.Text())
			m = append(m, deger)
		}
	})
	return m
}

func LatinceKaraktereCevir(Degisecek string) int {
	parcalanmis := strings.Split(Degisecek, "")
	var LatinceHali []int
	for _, deger := range parcalanmis {
		switch deger {
		case "٠":
			deger = "0"
		case "١":
			deger = "1"
		case "٢":
			deger = "2"
		case "۳":
			deger = "٣"
		case "٤":
			deger = "4"
		case "٥":
			deger = "5"
		case "٦":
			deger = "6"
		case "٧":
			deger = "7"
		case "٨":
			deger = "8"
		case "٩":
			deger = "9"
		default:
			fmt.Println("hataki giris")
		}

		temp, err := strconv.Atoi(deger)
		if err != nil {
			fmt.Println(err)
		}
		LatinceHali = append(LatinceHali, temp)
	}
	dondecek := LatinceHali[2] + LatinceHali[1] + LatinceHali[0]
	return dondecek
}
