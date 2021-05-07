package servis

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/djimenez/iconv-go"
	"github.com/muhammedaliakkaya/islamilimler/model"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func VerileriCek(id int, babNo int, sayfa int, kitapIsmi string) model.Hadis {
	url := fmt.Sprintf("http://islamilimleri.com/Ktphn/Kitablar/05/%03d/Arapca/%02d/%03d.htm", sayfa, id, babNo)
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
	iconv.Convert(data, yeni, "windows-1256", "utf-8")
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
	Baslık := BaslıkBul("CD33")
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
	hadis := model.Hadis{
		Kitap:  kitapIsmi,
		Konu:   Baslık,
		Numara: babNo,
		Metin:  metin,
	}
	return hadis
}
