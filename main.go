package main

import (
	"encoding/json"
	"fmt"
	"github.com/higashigo/go_currency_parser/utils"
	"os"
)


type Config struct {
	Database struct {
		Host     string `json:"host"`
		Password string `json:"password"`
		User 	 string `json:"user"`
		DBname   string `json:"dbname`
		Port     string `json:"port`
	} `json:"database"`
	ProxyUrl  string `json:"proxyUrl"`
	CurrUrl   string `json:"currUrl"`
	UserAgent string `json:"userAgent"`
}

func LoadConfiguration(file string) Config {
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}

func main() {
	confin := LoadConfiguration("config.json")
	// string for connectiong to database
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=disable",
							confin.Database.User,
							confin.Database.DBname,
							confin.Database.Password,
							confin.Database.Host,
							confin.Database.Port)

	userAgent := confin.UserAgent
	proxyPageUrl := confin.ProxyUrl
	currencyUrl := confin.CurrUrl

	htmlChannel := make(chan string)
	// create paeser
	parser := utils.Parser{UserAgent: userAgent}

	// search live proxy
	go parser.GetRequest(proxyPageUrl, "", htmlChannel)
	proxyHtml := <-htmlChannel
	liveProxy := parser.ParseProxyPage(proxyHtml)

	// parse data last 30 days
	datesForParse := utils.DatesParse()
	for _, date := range datesForParse{
		currUrl := fmt.Sprintf(currencyUrl, date)
		go parser.GetRequest(currUrl,liveProxy, htmlChannel)
	}

	for i:=0; i<len(datesForParse); i++{
		htmlCurPage := <-htmlChannel
		parser.ParseCurrPage(&htmlCurPage, &connStr)
	}
}
