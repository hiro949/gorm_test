package models

import (
	"fmt"
	"gormTest/config"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Model struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type Item struct {
	Model
	JanCode      string     `gorm:"size:255" json:"jan_code,omitempty"`
	ItemName     string     `gorm:"size:255" json:"item_name,omitempty"`
	Price        int        `json:"price,omitempty"`
	CategoryId   int        `json:"category_id,omitempty"`
	SeriesId     int        `json:"series_id,omitempty"`
	Stock        int        `json:"stock,omitempty"`
	Discontinued bool       `json:"discontinued"`
	ReleaseDate  *time.Time `json:"release_date,omitempty"`
}

var Db *gorm.DB

func GetAllItems(items *[]Item) {
	Db.Find(&items)
}

func GetSingleItem(item *Item, key string) {
	Db.First(&item, key)
}

func InsertItem(item *Item) {
	Db.NewRecord(item)
	Db.Create(&item)
}

func DeleteItem(key string) {
	Db.Where("id = ?", key).Delete(&Item{})
}

func UpdateItem(item *Item, key string) {
	Db.Model(&item).Where("id = ?", key).Updates(
		map[string]interface{}{
			"jan_code":     item.JanCode,
			"item_name":    item.ItemName,
			"price":        item.Price,
			"category_id":  item.CategoryId,
			"series_id":    item.SeriesId,
			"stock":        item.Stock,
			"discontinued": item.Discontinued,
			"release_date": item.ReleaseDate,
		})
}

// データベースの初期化
func init() {
	var err error
	dbConnectInfo := fmt.Sprintf(
		`%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local`,
		config.Config.DbUserName,
		config.Config.DbUserPassword,
		config.Config.DbHost,
		config.Config.DbPort,
		config.Config.DbName,
	)

	// configから読み込んだ情報を元に、データベースに接続します
	Db, err = gorm.Open(config.Config.DbDriverName, dbConnectInfo)
	if err != nil {
		log.Fatalln(err)
	} else {
		fmt.Println("Successfully connect database..")
	}

	// 接続したデータベースにitemsテーブルを作成します
	Db.Set("gorm:table_options", "ENGINE = InnoDB").AutoMigrate(&Item{})
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Successfully created table..")
	}
}
