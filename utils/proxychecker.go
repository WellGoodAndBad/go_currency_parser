package utils

import (
	"github.com/parnurzeal/gorequest"
)

func CheckProxy(proxies []string, userAgent string) string  {
	liveProxyChannel := make(chan string)
	for _, proxy := range proxies {
		proxyUrl := "http://" + proxy
		go requestForCheck(proxyUrl, userAgent, liveProxyChannel)
		}
	proxy := <-liveProxyChannel
	return proxy
	}


func requestForCheck(urlProxy string, userAgent string, retChan chan string) {
	request := gorequest.New().Proxy(urlProxy)
	resp, _, errs := request.Get("https://ya.ru").
		Set("User-Agent", userAgent).
		End()
	if errs != nil {
		//fmt.Println(errs)
		return
	} else {
		if resp.StatusCode == 200 {
			retChan <-urlProxy
		}
	}
}
