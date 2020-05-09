package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/xujiajun/nutsdb"
)

var (
	bucket = "dataList"
	key    = []byte("cd")
	keyDay = []byte("cd_day")
	today  = []byte("today")
)

//CjData data type
type CjData struct {
	HouseType int     `json:"house_type"` //1 商品房 2 二手房
	ZoneType  int     `json:"zone_type"`  // 1 中心城区 2 郊区新城 3 全市
	DataType  int     `json:"data_type"`  // 数据类型，1 总面积 2 住宅套数(套) 3 住宅面积(平方米) 4 非住宅面积
	Data      float64 `json:"data"`       // data
	DateTime  string  `json:"date_time"`
	Date      string  `json:"date"`
}

// Json return json string
func (d *CjData) Json() []byte {
	r, _ := json.Marshal(d)
	return r
}

// Save save data
func (d *CjData) Save() error {
	todayKey := fmt.Sprintf("%s_%d_%d_%d", today, d.HouseType, d.ZoneType, d.DataType)
	todayData := getToday(todayKey)
	if todayData == nil {
		_ = setToday(todayKey, d)
	} else {
		if todayData.Date == d.Date {
			_ = setToday(todayKey, d)
		} else {
			_ = deleteToday(todayKey)
			_ = db.Update(
				func(tx *nutsdb.Tx) error {
					return tx.RPush(bucket, keyDay, todayData.Json())
				})
		}
	}
	return db.Update(
		func(tx *nutsdb.Tx) error {
			return tx.RPush(bucket, key, d.Json())
		})
}

func getToday(tKey string) *CjData {
	var d *CjData
	if err := db.View(
		func(tx *nutsdb.Tx) error {
			if e, err := tx.Get(bucket, []byte(tKey)); err != nil {
				return err
			} else {
				d = &CjData{}
				_ = json.Unmarshal(e.Value, d)
			}
			return nil
		}); err != nil {
		log.Println(err)
	}
	return d
}

func setToday(tKey string, d *CjData) error {
	return db.Update(
		func(tx *nutsdb.Tx) error {
			if err := tx.Put(bucket, []byte(tKey), d.Json(), 0); err != nil {
				return err
			}
			return nil
		})
}

func deleteToday(tKey string) error {
	return db.Update(
		func(tx *nutsdb.Tx) error {
			if err := tx.Delete(bucket, []byte(tKey)); err != nil {
				return err
			}
			return nil
		})
}

func Query(zoneType, houseType, dataType int, day bool) []CjData {
	var k []byte
	var data []CjData

	if day == true {
		k = keyDay
	} else {
		k = key
	}
	if err := db.View(
		func(tx *nutsdb.Tx) error {
			if items, err := tx.LRange(bucket, k, 0, -1); err != nil {
				return err
			} else {
				for _, item := range items {
					d := &CjData{}
					_ = json.Unmarshal(item, d)
					if zoneType == d.ZoneType && dataType == d.DataType && houseType == d.HouseType {
						data = append(data, *d)
					}
				}
			}
			return nil
		}); err != nil {
		log.Println(err)
	}
	return data
}
