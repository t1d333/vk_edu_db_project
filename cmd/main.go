package main

import (
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

func main() {
	router := routing.New()
	fasthttp.ListenAndServe(":5000", router.HandleRequest)
}
