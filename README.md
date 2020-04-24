# showapi_sdk_go
showapi_sdk_go

| 方法名 | 说明| 参数|
|  -- | -- | -- |
| AddTextPara | 向请求中添加一个字符串类型的请求参数| 第一个参数:参数名,第二个参数: 字符串|
| AddFilePara | 向请求中添加一个File类型的请求参数| 第一个参数:参数名,第二个参数: File对象、byte字节、文件名String（三种类型皆可）|
| AddBase64Para | 向请求中添加一个字符串类型的请求参数| 第一个参数:参数名,第二个参数: 字符串文件名|
| Post | 以post的方法,发送参数到url地址,返回json字符串| |

## 第三方包下载说明
    需要到第三方包，mahonia（字符编码处理）,uuid（生成uuid）
    go env -w GO111MODULE=on
    go get -u "github.com/axgle/mahonia"
    go get -u "github.com/satori/go.uuid"

## 普通post请求示例
    import (
        . "../src/com/show/api"
        "fmt"
    )
    
    func main() {
        appid := "xxx"//要替换成自己的
        secret:="xxxxxxx"//要替换成自己的
        res := ShowapiRequest("http://route.showapi.com/64-19", appid, secret)
        res.AddTextPara("com", "zhongtong")
        res.AddTextPara("nu", "75312165465979")
        res.AddTextPara("senderPhone", "")
        res.AddTextPara("receiverPhone", "")
        fmt.Println(res.Post())
    }       


## 文件上传post请求示例

    import (
        . "../src/com/show/api"
        "fmt"
    )
    
    func main() {
        appid := "xxx"//要替换成自己的
        secret:="xxxxxxx"//要替换成自己的
        res := ShowapiRequest("http://route.showapi.com/1129-2", appid, secret)
        res.AddFilePara("imgFile", "替换为你的文件")//第一种：传入本地文件名
        //res.AddFilePara("imgFile", []byte{})//第二种：传入byte字节
        //res.AddFilePara("imgFile", os.File{})//第三种：传入本地文件
        fmt.Println(res.Post())
    }

## base64上传post请求示例
    import (
        . "../src/com/show/api"
        "fmt"
    )
    
    func main() {
        appid := "xxx"//要替换成自己的
        secret:="xxxxxxx"//要替换成自己的
        res := ShowapiRequest("http://route.showapi.com/1129-4", appid, secret)
        res.AddBase64Para("imgData", "替换为你的文件")//传入本地文件名
        fmt.Println(res.Post())
    }