package main

import "net/http"

var accept = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"
var acceptEncoding = "gzip, deflate"
var acceptLanguage = "zh-CN,zh;q=0.9"
var userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36"

func setCommonHeaders(req *http.Request) {
	req.Header.Set("Accept", accept)
	req.Header.Set("Accept-Encoding", acceptEncoding)
	req.Header.Set("Accept-Language", acceptLanguage)
	req.Header.Set("User-Agent", userAgent)
}
