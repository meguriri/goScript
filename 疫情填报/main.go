package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/spf13/viper"
)

var (
	UserName string
	PassWord string
	Name     string
	Stunum   string
	Tel      string
	College  string
	Province string
	City     string
	Id       string
)

func init() {
	viper.SetConfigName("info")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln(err)
	}
	UserName = viper.GetString("user.username")
	PassWord = viper.GetString("user.password")
	Name = viper.GetString("form.name")
	Stunum = viper.GetString("form.stunum")
	Tel = viper.GetString("form.tel")
	College = viper.GetString("form.college")
	Province = viper.GetString("form.province")
	City = viper.GetString("form.city")
	Id = viper.GetString("form.id")
}

func reader(res *http.Response) interface{} {
	s := make([]byte, 0)
	body := bytes.NewBuffer(s)
	body.ReadFrom(res.Body)
	message := make(map[string]interface{})
	err := json.Unmarshal(body.Bytes(), &message)
	if err != nil {
		log.Fatalln(err)
	}
	return message["m"]
}

func getCookie() []*http.Cookie {
	res, err := http.PostForm(
		"https://ucapp.sau.edu.cn/wap/login/invalid",
		url.Values{
			"username": {UserName},
			"password": {PassWord},
		},
	)
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()

	log.Println("获取Cookies: ", reader(res))

	return res.Cookies()
}

func UploadInfo(cookies []*http.Cookie) {

	data := time.Now().Format("2006-01-02")
	temp := []string{"36.3", "36.3", "36.3"}
	c := new(http.Client)

	req, err := http.NewRequest(
		"POST",
		"https://app.sau.edu.cn/form/wap/default/save?formid=10",
		strings.NewReader(url.Values{
			"xingming":                       {Name},
			"xuehao":                         {Stunum},
			"shoujihao":                      {Tel},
			"danweiyuanxi":                   {College},
			"dangqiansuozaishengfen":         {Province},
			"dangqiansuozaichengshi":         {City},
			"shifouyuhubeiwuhanrenyuanmiqie": {"否"},
			"shifoujiankangqingkuang":        {"是"},
			"shifoujiechuguohubeihuoqitayou": {"否"},
			"fanhuididian":                   {""},
			"shifouweigelirenyuan":           {"否"},
			"shentishifouyoubushizhengzhuan": {"否"},
			"shifouyoufare":                  {"否"},
			"qitaxinxi":                      {""},
			"tiwen":                          {temp[0]},
			"tiwen1":                         {temp[1]},
			"tiwen2":                         {temp[2]},
			"riqi":                           {data},
			"id":                             {Id},
		}.Encode()),
	)
	if err != nil {
		log.Fatalln(err)
	}

	for _, v := range cookies {
		req.AddCookie(v)
	}

	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Accept-Encoding", "")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Length", "528")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Host", "app.sau.edu.cn")
	req.Header.Set("Origin", "https://app.sau.edu.cn")
	req.Header.Set("Referer", "https://app.sau.edu.cn/form/wap/default/index?formid=10&nn=7026.582142720368")
	req.Header.Set("sec-ch-ua", "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"100\", \"Google Chrome\";v=\"100\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"Windows\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	res, err := c.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()

	log.Println("上传填报信息: ", reader(res))
}

func main() {
	cookies := getCookie()
	UploadInfo(cookies)
}
