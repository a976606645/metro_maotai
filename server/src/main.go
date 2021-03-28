package main

import (
	"math/rand"
	"mdl/db"
	"mdl/wlog"
	"os/exec"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

var DB *gorm.DB
var isRun bool

func main() {
	wlog.DevelopMode()
	wlog.Info("Author:huwei365@hotmail.com")

	rand.Seed(time.Now().UnixNano())

	DB = db.NewSqlite()
	DB.AutoMigrate(&TableUser{})
	DB.AutoMigrate(&TableOrder{})
	DB.AutoMigrate(&TableStockRecord{})

	c := cron.New(cron.WithSeconds())
	c.AddFunc("59 59 18 * * *", Start)
	c.AddFunc("0 5 19 * * *", Stop)
	c.AddFunc("0 0 19 * * *", CollectStock)
	c.AddFunc("0 0 20 * * *", UpdateUser)
	go c.Run()

	g := gin.Default()
	g.Static("/m", "./www/")
	g.Static("/static", "./www/static")

	g.POST("user/sendSms", SendSms)
	g.POST("user/add", UserAdd)
	g.POST("user/list", UserList)
	g.POST("user/delete", UserDelete)
	g.POST("user/edit", UserEdit)
	g.POST("user/setStore", SetStore)

	g.POST("store/list", GetStoreList)
	g.POST("order/list", OrderList)

	if runtime.GOOS == "windows" {
		go func() {
			time.Sleep(1 * time.Second)
			err := exec.Command(`cmd`, `/c`, `start`, "http://localhost:27777/m").Start()
			if err != nil {
				wlog.Error("弹出浏览器失败:", err)
			}
		}()
	}

	g.Run(":27777")
	// TmpPurchase()
}
