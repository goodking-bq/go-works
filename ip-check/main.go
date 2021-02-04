package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"log"
	"os"
	"strings"

	myIpDB "ip-check/ipdb"
)

var (
	port    string
	host    string
	help    bool
	fRegion string
	fCity   string
	db      *ipdb.City
)

func getCityInfo(ip string) *ipdb.CityInfo {
	info, _ := db.FindInfo(ip, "CN")
	return info
}

func getRemoteIp(ctx *fasthttp.RequestCtx) string {
	remoteIp := string(ctx.QueryArgs().Peek("ip"))
	if remoteIp == "" {
		remoteIp = string(ctx.Request.Header.Peek("X-Forwarded-For"))
		if index := strings.IndexByte(remoteIp, ','); index >= 0 {
			remoteIp = remoteIp[0:index]
			//获取最开始的一个 即 1.1.1.1
		}
		remoteIp = strings.TrimSpace(remoteIp)
		if len(remoteIp) > 0 {
			return remoteIp
		}
		remoteIp = strings.TrimSpace(string(ctx.Request.Header.Peek("X-Real-Ip")))
		if len(remoteIp) > 0 {
			return remoteIp
		}
		return ctx.RemoteIP().String()
	}
	return remoteIp

}

// Get get url
func Get(ctx *fasthttp.RequestCtx) {
	remoteIp := getRemoteIp(ctx)
	fmt.Println(remoteIp)
	info, _ := json.Marshal(getCityInfo(remoteIp))
	_, _ = fmt.Fprintf(ctx, string(info))
}

func Check(ctx *fasthttp.RequestCtx) {
	remoteIp := getRemoteIp(ctx)
	info := getCityInfo(remoteIp)
	if info.RegionName != "" && strings.Contains(fRegion, info.RegionName) {
		_, _ = fmt.Fprintf(ctx, "0")
		return
	}
	if info.CityName != "" && strings.Contains(fCity, info.CityName) {
		_, _ = fmt.Fprintf(ctx, "0")
		return
	}
	_, _ = fmt.Fprintf(ctx, "1")
}

func Help(ctx *fasthttp.RequestCtx) {
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
	ctx.Response.Header.Set("Content-Type", "text/html; charset=utf-8")
	_, _ = fmt.Fprintf(ctx, html)
}

func Test(ctx *fasthttp.RequestCtx) {
	_, _ = fmt.Fprintf(ctx, "this is the second part of body\n")
}
func main() {
	flag.BoolVar(&help, "h", false, "帮助")
	flag.StringVar(&host, "host", "", "侦听ip")
	flag.StringVar(&port, "p", "8000", "侦听端口")
	flag.StringVar(&fRegion, "r", "", "检查时禁用省份，逗号隔开")
	flag.StringVar(&fCity, "c", "", "检查时禁用城市，逗号隔开")
	flag.Parse()
	if help {
		_, _ = fmt.Fprintf(os.Stderr, `Usage: ipcheck [-host [host]] [-p [port]]

Options:
`)
		flag.PrintDefaults()
		return
	}
	body, _ := myIpDB.Asset("ipipfree.ipdb")
	db, _ = ipdb.NewCityBytes(body)

	router := fasthttprouter.New()
	router.GET("/get", Get)
	router.GET("/check", Check)
	router.GET("/Help", Help)
	router.GET("/test", Test)

	fmt.Printf("start server at %s:%s\n", host, port)
	log.Fatal(fasthttp.ListenAndServe(host+":"+port, router.Handler))

}
