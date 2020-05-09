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
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/xujiajun/nutsdb"

	"zfcj/cmd/bindata"

	"github.com/spf13/cobra"
)

func apiHandler(w http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	var (
		zoneType  = 0
		houseType = 0
		dataType  = 0
		timeType  = true
	)
	if params.Get("zoneType") != "" {
		zoneType, _ = strconv.Atoi(params.Get("zoneType"))
	}
	if params.Get("houseType") != "" {
		houseType, _ = strconv.Atoi(params.Get("houseType"))
	}
	if params.Get("dataType") != "" {
		dataType, _ = strconv.Atoi(params.Get("dataType"))
	}
	if params.Get("timeType") != "" {
		t, _ := strconv.Atoi(params.Get("timeType"))
		if t == 1 {
			timeType = false
		}
	}
	w.Header().Set("Content-Type", "application/json")
	data := Query(zoneType, houseType, dataType, timeType)
	dataStr, _ := json.Marshal(data)
	io.WriteString(w, string(dataStr))
}

func reloadHandler(w http.ResponseWriter, req *http.Request) {
	_ = db.Close()
	opt := nutsdb.DefaultOptions
	opt.Dir = "/tmp/nutsdb" //这边数据库会自动创建这个目录文件
	_db, err := nutsdb.Open(opt)
	if err != nil {
		log.Fatal(err)
	}
	db = _db
	_, _ = w.Write([]byte("ok"))
}

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("serve called")
		http.Handle("/",
			http.FileServer(
				&assetfs.AssetFS{Asset: bindata.Asset, AssetDir: bindata.AssetDir, AssetInfo: bindata.AssetInfo, Prefix: "ui/dist/"}))
		http.HandleFunc("/api", apiHandler)
		http.HandleFunc("/reload", reloadHandler)
		http.ListenAndServe(":8088", nil)
	},
}

func startCrawl(cmd *cobra.Command, args []string) {
	tick := time.NewTicker(time.Hour)
	crawlUrl("https://zw.cdzj.chengdu.gov.cn/py/SCXX/Default.aspx?action=ucEveryday")
	go func(tick *time.Ticker) {
		for {
			select {
			case <-tick.C:
				crawlUrl("https://zw.cdzj.chengdu.gov.cn/py/SCXX/Default.aspx?action=ucEveryday")
			}
			time.Sleep(time.Millisecond)
		}
	}(tick)
}

func init() {
	rootCmd.AddCommand(serveCmd)
	//serveCmd.PreRun = startCrawl
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
