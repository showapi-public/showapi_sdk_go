package main

import (
	"showSdk/normalRequest"
	"fmt"
)

func main() {
	res := normalRequest.ShowapiRequest("http://route.showapi.com/66-22", 1, "xxxxxxxxxxxxxx")
	res.AddTextPara("code", "6912345678901")
	fmt.Println(res.Get())
}
