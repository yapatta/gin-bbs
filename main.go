package main

import (
	//"log"
	//"net/http"
	"strconv"
	//"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type Tweet struct {
	gorm.Model
	Name    string
	Comment string
}

func dbInit() {
	db := dbOpen()
	db.AutoMigrate(&Tweet{})
	defer db.Close()
}

func dbOpen() *gorm.DB {
	db, err := gorm.Open("mysql", "gorm:password@/bbs?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("データベースが開けません！")
	}
	return db
}

func dbInsert(name string, comment string) {
	db := dbOpen()
	db.Create(&Tweet{Name: name, Comment: comment})
	defer db.Close()
}

func dbUpdate(id int, name string, comment string) {
	db := dbOpen()
	var tweet Tweet
	db.First(&tweet, id)
	tweet.Name = name
	tweet.Comment = comment
	db.Save(&tweet)
	db.Close()
}

func dbDelete(id int) {
	db := dbOpen()
	var tweet Tweet
	db.First(&tweet, id)
	db.Delete(&tweet)
	db.Close()
}

func dbGetAll() []Tweet {
	db := dbOpen()
	var tweets []Tweet
	db.Order("created_at desc").Find(&tweets)
	db.Close()
	return tweets
}

func dbGetOne(id int) Tweet {
	db := dbOpen()
	var tweet Tweet
	db.First(&tweet, id)
	db.Close()
	return tweet
}

// main...

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templetes/*.html")

	dbInit()

	router.GET("/", func(ctx *gin.Context) {
		tweets := dbGetAll()
		ctx.HTML(200, "index.html", gin.H{"tweets": tweets})
	})

	router.POST("/new", func(ctx *gin.Context) {
		name := ctx.PostForm("name")
		comment := ctx.PostForm("comment")
		dbInsert(name, comment)
		ctx.Redirect(302, "/")
	})

	//Detail
	router.GET("/detail/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		tweet := dbGetOne(id)
		ctx.HTML(200, "detail.html", gin.H{"tweet": tweet})
	})

	//Update
	router.POST("/update/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		name := ctx.PostForm("name")
		comment := ctx.PostForm("comment")
		dbUpdate(id, name, comment)
		ctx.Redirect(302, "/")
	})
	//削除確認
	router.GET("/delete_check/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		tweet := dbGetOne(id)
		ctx.HTML(200, "delete.html", gin.H{"tweet": tweet})
	})

	//Delete
	router.POST("/delete/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		dbDelete(id)
		ctx.Redirect(302, "/")
	})

	router.Run()
}
