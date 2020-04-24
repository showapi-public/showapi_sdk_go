package api

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	_ "encoding/json"
	"errors"
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/satori/go.uuid"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type NormalReq struct {
	url         string
	textMap     url.Values
	uploadMap   map[string]interface{}
	timeout     time.Duration
	headMap     url.Values
	charset     string
	headString  string
	bodyString  string
	res_headMap map[string]string
}

//用于请求官网
func ShowapiRequest(reqUrl string, appid int, sign string) *NormalReq {
	values := make(url.Values)
	values.Set("showapi_appid", strconv.Itoa(appid))
	values.Set("showapi_sign", sign)
	return &NormalReq{reqUrl, values, nil, 30000 * time.Millisecond, values, "utf-8", "", "", make(map[string]string)}
}

//PostAsBytes或GetAsBytes请求实例
func PostOrGetAsByteDemo() {
	r := NormalRequest("http://route.showapi.com/184-4")
	r.AddTextPara("showapi_appid", "")
	r.AddTextPara("showapi_sign", "")
	r.AddTextPara("typeId", "34")
	r.AddFilePara("image", "D:\\temp_img\\fc\\3.jpg")
	res, _ := r.PostAsBytes()
	//res,_ := r.GetAsBytes()
	fmt.Println(string(res))
}

//Post或Get请求实例
func TestPostOrGetDemo() {
	r := NormalRequest("http://route.showapi.com/184-5")
	r.AddHeadPara("Content-Type", "application/x-www-form-urlencoded")
	r.AddTextPara("showapi_appid", "")
	r.AddTextPara("showapi_sign", "")
	r.AddTextPara("img_base64", "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")
	r.AddTextPara("typeId", "34")
	res, _ := r.Post()
	//res,_ := r.Get()
	fmt.Println(res)
}

func NormalRequest(reqUrl string) *NormalReq {
	values := make(url.Values)
	values.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.92 Safari/537.36")
	return &NormalReq{reqUrl, make(url.Values), nil, 30000 * time.Millisecond, values, "utf-8", "", "", make(map[string]string)}
}

func (request *NormalReq) AddTextPara(key, value string) {
	request.textMap.Set(key, value)
}

func (request *NormalReq) AddFilePara(key string, fileName interface{}) {

	if request.uploadMap == nil {
		request.uploadMap = make(map[string]interface{})
	}
	request.uploadMap[key] = fileName
}

func (request *NormalReq) AddHeadPara(key, value string) {
	request.headMap.Set(key, value)
}
func (request *NormalReq) AddBase64Para(key, fileName string) {
	request.textMap.Set(key, Base64(fileName))
}

func (request *NormalReq) SetTimeOut(timeout time.Duration) {
	request.timeout = timeout * time.Millisecond
}
func (request *NormalReq) SetCharset(charset string) {
	request.charset = charset
}
func (request *NormalReq) SetBodyString(bodyString string) {
	request.bodyString = bodyString
}
func (request *NormalReq) SetHeadString(headString string) {
	request.headString = headString
}

func (request *NormalReq) Url() string {
	reqUrl := request.url
	return reqUrl
}
func (request *NormalReq) HeadMap() map[string]interface{} {
	headMap := make(map[string]interface{})
	for k, v := range request.headMap {
		headMap[k] = v[0]
	}
	return headMap
}
func (request *NormalReq) TextMap() map[string]interface{} {
	textMap := make(map[string]interface{})
	for k, v := range request.textMap {
		textMap[k] = v[0]
	}
	if request.uploadMap != nil {
		for k, v := range request.uploadMap {
			textMap[k] = v
		}
	}
	return textMap
}
func (request *NormalReq) Res_headMap() map[string]string {
	res_headMap := make(map[string]string)
	for k, v := range request.res_headMap {
		res_headMap[k] = v
	}
	return res_headMap
}

func (request *NormalReq) Get() (string, error) {
	tr := &http.Transport{
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           myDialContext,
		MaxIdleConns:          20,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	client := &http.Client{
		Timeout:   request.timeout,
		Transport: tr,
	}
	reqUrl := ""
	if strings.Contains(request.url, "?") {
		reqUrl = strings.TrimSpace(request.url) + "&" + request.textMap.Encode()
	} else {
		reqUrl = strings.TrimSpace(request.url) + "?" + request.textMap.Encode()
	}
	req, err := http.NewRequest("GET", reqUrl, nil)
	if request.headString != "" {
		strlist := strings.Split(request.headString, "\r\n")
		if len(strlist) <= 1 {
			strlist = strings.Split(request.headString, "\n")
		}
		for i := 0; i < len(strlist); i++ {
			setheadlist := strings.Split(strlist[i], ":")
			req.Header.Add(strings.TrimSpace(setheadlist[0]), strings.TrimSpace(setheadlist[1]))
		}
	}
	for k, v := range request.headMap {
		req.Header.Add(k, v[0])
	}
	if err != nil {
		return "request err", err
	}
	resp, err := client.Do(req)
	for k, v := range resp.Header {
		request.res_headMap[k] = v[0]
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "response err", err
	}
	return string(body), nil
}
func (request *NormalReq) GetAsBytes() ([]byte, error) {
	body, err := request.Get()
	if err != nil {
		return nil, err
	}
	return []byte(body), nil
}

func (request *NormalReq) PostAsBytes() ([]byte, error) {
	if request.bodyString != "" {
		request.AddHeadPara("Content-Type", "application/octet-stream;charset="+request.charset)
	} else if len(request.headMap.Get("Content-Type")) <= 0 {
		request.AddHeadPara("Content-Type", "application/x-www-form-urlencoded;charset="+request.charset)
	}
	body, err := request.Post()
	if err != nil {
		return nil, err
	}
	bodyCoder := mahonia.NewDecoder(request.charset)
	bodyResult := bodyCoder.ConvertString(string(body))
	tagCoder := mahonia.NewDecoder("UTF-8")
	_, data, _ := tagCoder.Translate([]byte(bodyResult), true)
	return data, nil
}

//文件转base64
func Base64(fileName string) string {
	fileBase64, _ := ioutil.ReadFile(fileName)
	return base64.StdEncoding.EncodeToString(fileBase64)
}

var myDial = &net.Dialer{
	Timeout:   30 * time.Second,
	KeepAlive: 30 * time.Second,
	DualStack: true,
}
var myDialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
	network = "tcp4" //仅使用ipv4
	//network = "tcp6" //仅使用ipv6
	return myDial.DialContext(ctx, network, addr)
}

func (request *NormalReq) Post() (string, error) {
	tr := &http.Transport{
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           myDialContext,
		MaxIdleConns:          20,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	client := &http.Client{
		Timeout:   request.timeout,
		Transport: tr,
	}
	var resp *http.Response
	if request != nil && request.uploadMap != nil && len(request.uploadMap) > 0 {
		var fileName, filePath string
		for k, v := range request.uploadMap {
			fileName = k
			typeOfUpload := reflect.TypeOf(v)
			if typeOfUpload.Kind() == reflect.String {
				filePath = v.(string)
			} else if typeOfUpload.Kind() == reflect.Ptr {
				file := v.(*os.File)
				if file == nil {
					return "", errors.New("请确认上传的文件是否正确")
				}
				filePath = file.Name()
			} else if typeOfUpload.Kind() == reflect.Slice {
				fileBytes := v.([]byte)
				fileType := http.DetectContentType(fileBytes) //文件后缀
				u, _ := uuid.NewV4()
				fileName := strings.Replace(u.String(), "-", "", -1)
				fileEndings, err := mime.ExtensionsByType(fileType)
				if err != nil {
					return "", err
				}
				newPath := filepath.Join("/NFS/hby/upload", fileName+fileEndings[0]) //生成完整的文件名
				newFile, err := os.Create(newPath)
				if err != nil {
					return "", err
				}
				if _, err := newFile.Write(fileBytes); err != nil {
					return "", err
				}
				defer os.Remove(newFile.Name())
				filePath = newFile.Name()
				newFile.Close()
			} else {
				return "", errors.New("上传格式为string或File或[]Byte")
			}
		}
		//只支持单文件上传
		file, err := os.Open(filePath)
		if err != nil {
			return "", err
		}
		defer file.Close()
		postbody := &bytes.Buffer{}
		writer := multipart.NewWriter(postbody)
		part, err := writer.CreateFormFile(fileName, filepath.Base(filePath))
		if err != nil {
			return "", err
		}
		_, err = io.Copy(part, file)
		for k, v := range request.textMap {
			writer.WriteField(k, v[0])
		}
		err = writer.Close()
		if err != nil {
			return "request err", err
		}
		req, err := http.NewRequest("POST", request.Url(), postbody)
		req.Header.Add("Content-Type", writer.FormDataContentType())
		for k, v := range request.headMap {
			if k != "Content-Type" {
				req.Header.Add(k, v[0])
			}
		}
		if request.headString != "" {
			strlist := strings.Split(request.headString, "\r\n")
			if len(strlist) <= 1 {
				strlist = strings.Split(request.headString, "\n")
			}
			for i := 0; i < len(strlist); i++ {
				setheadlist := strings.Split(strlist[i], ":")
				req.Header.Add(strings.TrimSpace(setheadlist[0]), strings.TrimSpace(setheadlist[1]))
			}
		}
		resp, err = client.Do(req)
		if err != nil {
			return "request err", err
		}

	} else {
		bodystr := strings.TrimSpace(request.textMap.Encode())
		if request.bodyString != "" {
			bodystr = bodystr + strings.TrimSpace(request.bodyString)
		}
		req, err := http.NewRequest("POST", strings.TrimSpace(request.url), strings.NewReader(bodystr))
		if request.headString != "" {
			strlist := strings.Split(request.headString, "\r\n")
			if len(strlist) <= 1 {
				strlist = strings.Split(request.headString, "\n")
			}
			for i := 0; i < len(strlist); i++ {
				setheadlist := strings.Split(strlist[i], ":")
				req.Header.Add(strings.TrimSpace(setheadlist[0]), strings.TrimSpace(setheadlist[1]))
			}
		}
		for k, v := range request.headMap {
			if k == "Content-Type" && !strings.Contains(v[0], "charset") {
				req.Header.Add(k, v[0]+";charset="+request.charset)
			} else {
				req.Header.Add(k, v[0])
			}
		}
		resp, err = client.Do(req)
		if err != nil {
			return "request err", err
		}
	}
	for k, v := range resp.Header {
		request.res_headMap[k] = v[0]
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "response err", err
	}
	return string(body), nil
}
