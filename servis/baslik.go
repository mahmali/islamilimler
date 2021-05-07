package servis

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/djimenez/iconv-go"
	"log"
	"net/http"
	"strings"
)

func Baslikgettir(kitapNo int) string {

	url := fmt.Sprintf("http://islamilimleri.com/Ktphn/Kitablar/05/%03d/Arapca/01/000.htm", kitapNo)
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
	baslık = strings.ReplaceAll(baslık, "  ", "")
	baslık = strings.TrimSpace(baslık)
	fmt.Println(baslık)

	return baslık
}
