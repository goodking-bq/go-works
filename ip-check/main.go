package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/valyala/fasthttp"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	myIpDB "ip-check/ipdb"

	"ip-check/ipdb-go"
)

var (
	port    string
	host    string
	help    bool
	fRegion string
	fCity   string
)

func getCityInfo(ip string) *ipdb.CityInfo {
	body, _ := myIpDB.Asset("ipipfree.ipdb")
	db, err := ipdb.NewCityBytes(body)
	if err != nil {
		log.Fatal(err)
	}
	info, _ := db.FindInfo(ip, "CN")

	return info
}


func getRemoteIp(ctx *fasthttp.RequestCtx) string{
	remoteIp := string(ctx.QueryArgs().Peek("ip"))
	if remoteIp == nil {
		remoteIp = ctx.RemoteIP().String()
	}
	return remoteIp
}


// Get get url
func Get(ctx *fasthttp.RequestCtx) {
	remoteIp:=getRemoteIp(ctx)
	info, _ := json.Marshal(getCityInfo(remoteIp))
	i,_:=fmt.Fprintln(ctx,info)
	fmt.Println(i)
}



func checkRequest(w http.ResponseWriter, r *http.Request) {
	remoteIp:=getRemoteIp(r)
	info := getCityInfo(remoteIp)
	if info.RegionName != "" && strings.Contains(fRegion, info.RegionName) {
		io.WriteString(w, "0")
		return
	}
	if info.CityName != "" && strings.Contains(fCity, info.CityName) {
		io.WriteString(w, "0")
		return
	}
	io.WriteString(w, "1")
}

func helpRequest(w http.ResponseWriter, r *http.Request) {
	html := `<html>

<head>
    <title>帮助</title>
</head>

<body>
    <ul>
        <li> 获取IP的区域信息
            <ul>
                <li>URL:
                    <code>/get</code> <strong>获取本机</strong>
                </li>
                <li>URL: <code>/get?ip=1.1.1.1</code> <strong>获取特定ip</strong></li>
                返回JSON字符串：<br /><code>{"country_name":"中国","region_name":"四川","city_name":"成都","owner_domain":"","isp_domain":"","latitude":"","longitude":"","timezone":"","utc_offset":"","china_admin_code":"","idd_code":"","country_code":"","continent_code":"","idc":"","base_station":"","country_code3":"","european_union":"","currency_code":"","currency_name":"","anycast":"","line":"","district_info":{"country_name":"","region_name":"","city_name":"","district_name":"","china_admin_code":"","covering_radius":"","latitude":"","longitude":""},"route":"","asn":"","asn_info":null,"area_code":""}</code>
            </ul>
        </li>
        <br />
        <li> 检查ip是否满足条件
            <ul>
                <li>URL:
                    <code>/check</code> <strong>检查本机</strong>
                </li>
                <li>URL: <code>/check?ip=1.1.1.1</code> <strong>检查特定ip</strong></li>
                返回数字：<br /><code>0(未通过) 或 1(通过) </code>
            </ul>
        </li>
    </ul>
</body>

</html>`
	io.WriteString(w, html)
}

func main() {
	flag.BoolVar(&help, "h", false, "帮助")
	flag.StringVar(&host, "host", "", "侦听ip")
	flag.StringVar(&port, "p", "8000", "侦听端口")
	flag.StringVar(&fRegion, "r", "", "检查时禁用省份，逗号隔开")
	flag.StringVar(&fCity, "c", "", "检查时禁用城市，逗号隔开")
	flag.Parse()
	if help {
		fmt.Fprintf(os.Stderr, `Usage: ipcheck [-host [host]] [-p [port]]

Options:
`)
		flag.PrintDefaults()
		return
	}
	http.HandleFunc("/get", doRequest)      //   设置访问路由
	http.HandleFunc("/help", helpRequest)   //   帮助
	http.HandleFunc("/check", checkRequest) //   帮助
	fmt.Printf("start server at %s:%s", host, port)
	err := http.ListenAndServe(host+":"+port, nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
