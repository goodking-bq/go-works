/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

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
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/spf13/cobra"
)

// crawlUrl crawl
func crawlUrl(url string) {
	defer http.Get("http://localhost:8088/reload")
	client := &http.Client{
		Timeout: 3 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
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
	now := time.Now()
	for i, table := range htmlquery.Find(doc, "//table") {
		for j, tr := range htmlquery.Find(table, "/tbody/tr") {
			if j < 2 {
				continue
			}
			tds := htmlquery.Find(tr, "/td")
			for dataType, td := range tds {
				if dataType > 0 {
					d, _ := strconv.ParseFloat(strings.Trim(htmlquery.InnerText(td), "\n "), 32)
					data := &CjData{
						HouseType: (i + 1) / 3,
						ZoneType:  j - 1,
						DataType:  dataType,
						DateTime:  now.Format("2006-01-02 15:04:05"),
						Data:      d,
						Date:      now.Format("2006-01-02"),
					}
					if err := data.Save(); err != nil {
						log.Println(err)
					}
				}
			}

		}
	}
	log.Println("crawl success.")
}

// crawlCmd represents the crawl command
var crawlCmd = &cobra.Command{
	Use:   "crawl",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("crawl called")
		crawlUrl("https://zw.cdzj.chengdu.gov.cn/py/SCXX/Default.aspx?action=ucEveryday")
	},
}

func init() {
	rootCmd.AddCommand(crawlCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// crawlCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// crawlCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
