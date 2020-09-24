package utils

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/parnurzeal/gorequest"
	"log"
	"strings"
)


type DataCur struct {
	CurrencyCode string
	CurrencyName string
	UnitsPerUSD string
	USDPerUnit string
}

type Parser struct {
	UserAgent string
	Proxy string

}

func (p Parser)GetRequest(url string, proxy string, retChan chan string)  {
	request := gorequest.New()
	if proxy !=  ""{
		request.Proxy(proxy)
	}
	resp, body, errs := request.Get(url).
		Set("User-Agent", p.UserAgent).
		End()
	if errs != nil {
		log.Fatalln(errs)
	}

	statusCode := resp.StatusCode

	if statusCode != 200 {
		log.Fatal("Status code == %s", statusCode)
	}
	retChan <- body
}

func (p Parser) ParseProxyPage(htmlbody string) string {

	doc := GoqueryDoc(htmlbody)
	textArea := doc.Find("textarea.form-control").Text()
	textAreaStr := strings.Split(textArea, "\n")
	proxiesList := textAreaStr[3:len(textAreaStr)-1]
	liveProxy := CheckProxy(proxiesList, p.UserAgent)

	return liveProxy
}

func (p Parser) ParseCurrPage(htmlBody *string, connStr *string){

	var dataForDb []DataCur
	doc := GoqueryDoc(*htmlBody)

	dateParse, _ := doc.Find("input#ratesDate").Attr("value") //get date parse from page

	doc.Find("table#historicalRateTbl").Find("tbody").Find("tr").Each(func(_ int, s *goquery.Selection) {
		var dataTr DataCur
		s.Find("td").Each(func(i int, selTd *goquery.Selection){
			if i == 0 {
				dataTr.CurrencyCode = strings.TrimSpace(selTd.Text())
			}
			if i == 1 {
				dataTr.CurrencyName = strings.TrimSpace(selTd.Text())
			}
			if i == 2 {
				dataTr.UnitsPerUSD = strings.TrimSpace(selTd.Text())
			}
			if i == 3 {
				dataTr.USDPerUnit = strings.TrimSpace(selTd.Text())
			}
		})
		dataForDb = append(dataForDb, dataTr)
	})
	InsertData(&dataForDb, dateParse, *connStr)

}

func GoqueryDoc(htmlbody string) *goquery.Document {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlbody))
	if err != nil {
		log.Fatal(err)
	}
	return doc
}