/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/antchfx/htmlquery"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cobra"
)

// exporterCmd represents the exporter command
var exporterCmd = &cobra.Command{
	Use:   "exporter",
	Short: "成都住房交易数据 for prometheus exporter",
	Long:  `成都住房交易数据`,
	Run: func(cmd *cobra.Command, args []string) {
		url, _ := cmd.Flags().GetString("url")
		e := NewExporter("house", url)
		// 注册一个采集器
		reg := prometheus.NewRegistry()
		reg.MustRegister(e)
		prometheus.Unregister(prometheus.NewGoCollector())
		gatherers := prometheus.Gatherers{
			prometheus.DefaultGatherer,
			reg,
		}

		h := promhttp.HandlerFor(gatherers,
			promhttp.HandlerOpts{
				ErrorHandling: promhttp.ContinueOnError,
			})
		http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		})
		http.HandleFunc("/health", func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte("ok"))
		})
		log.Println("Start server at :9000")
		if err := http.ListenAndServe(":9000", nil); err != nil {
			log.Printf("Error occur when start server %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(exporterCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// exporterCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	exporterCmd.Flags().String("url", URL, "Help message for toggle")
}

const URL = "https://zw.cdzj.chengdu.gov.cn/zwdt/SCXX/Default.aspx?action=ucEveryday"

var (
	dataTypes = map[int]string{
		1: "area",
		2: "number",
		3: "area",
		4: "area",
	}
	zoneTypes = map[int]string{
		1: "center",
		2: "outer",
		3: "all",
	}
	houseTypes = map[int]string{
		1: "new",
		2: "old",
	}
)

type Exporter struct {
	ns         string
	url        string
	AreaDesc   *prometheus.Desc
	NumberDesc *prometheus.Desc
}

func NewExporter(ns, url string) *Exporter {
	return &Exporter{
		ns:  ns,
		url: url,
		AreaDesc: prometheus.NewDesc(prometheus.BuildFQName(ns, "",
			"sold_area"), "销售面积", []string{"zone", "type"}, nil),
		NumberDesc: prometheus.NewDesc(prometheus.BuildFQName(ns, "",
			"sold_number"), "销售套数", []string{"zone", "type"}, nil),
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.AreaDesc
	ch <- e.NumberDesc
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	client := &http.Client{
		Timeout: 3 * time.Second,
	}
	req, err := http.NewRequest("GET", e.url, nil)
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows; U; Windows NT 6.1; en-US) AppleWebKit/532.5 (KHTML, like Gecko) Chrome/4.0.249.0 Safari/532.5")
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Println("status code error: " + resp.Status)
	}
	//body, err := ioutil.ReadAll(resp.Body)

	doc, err := htmlquery.Parse(resp.Body) //htmlquery.LoadDoc("/Volumes/work/go-program-private/src/zfcj/doc") // htmlquery.LoadURL(url)

	if err != nil {
		log.Println("crawl err: ", err)
	}
	for i, table := range htmlquery.Find(doc, "//table") {
		for j, tr := range htmlquery.Find(table, "/tbody/tr") {
			if j < 2 {
				continue
			}
			tds := htmlquery.Find(tr, "/td")
			for dataType, td := range tds {
				switch dataType {
				case 2:
					d, _ := strconv.ParseFloat(strings.Trim(htmlquery.InnerText(td), "\n "), 32)
					if m, err := prometheus.NewConstMetric(e.NumberDesc, prometheus.GaugeValue, d, zoneTypes[j-1],
						houseTypes[(i+1)/3]); err == nil {
						ch <- m
					}
				case 3:
					d, _ := strconv.ParseFloat(strings.Trim(htmlquery.InnerText(td), "\n "), 32)
					if m, err := prometheus.NewConstMetric(e.AreaDesc, prometheus.GaugeValue, d, zoneTypes[j-1],
						houseTypes[(i+1)/3]); err == nil {
						ch <- m
					}
				}
			}

		}
	}
}
