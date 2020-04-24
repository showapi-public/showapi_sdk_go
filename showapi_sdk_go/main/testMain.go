package main

import (
	. "../src/com/show/api"
	"fmt"
)

func main() {
	res := ShowapiRequest("http://route.showapi.com/66-22", 1, "xxxxxxxxxxxxxx")
	res.AddTextPara("code", "")
	res.AddBase64Para("img_base64", "")
	fmt.Println(res.Post())
}
