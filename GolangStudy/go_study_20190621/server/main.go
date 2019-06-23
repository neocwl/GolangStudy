package main

import (
	"GolangStudy/GolangStudy/go_study_20190621/modle"
	"GolangStudy/GolangStudy/go_study_20190621/modle/bookSet"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"net/http"
)

const (
	KEY_BOOK_DETAIL_IN_REDIS = "book_detail"
	KEY_BOOK_IN_REDIS        = "book"
)

//定义一个全局的pool
var pool *redis.Pool
var conn redis.Conn

func init() {

	pool = &redis.Pool{
		MaxIdle:     0,   //最大空闲连接数
		MaxActive:   0,   //表示和数据库的最大连接数，0表示没有限制
		IdleTimeout: 100, //最大空闲时间单位：秒
		Dial: func() (conn redis.Conn, e error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}
}

func main() {
	conn = pool.Get()
	defer conn.Close()

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.GET("/spider/bookset/:key", getBooks)
	r.Run(":8880")
}

func getBooks(c *gin.Context) {
	key := c.Param("key")
	start := c.DefaultQuery("start", "0")
	end := c.Query("end")

	//获取结果
	result, err := redis.Strings(conn.Do("lrange", key, start, end))

	if err != nil {
		msg := modle.Message{ErrCode: modle.MESSAGE_CODE_QUERY_FAILED,
			Error: err.Error(),
			Data:  "",
		}
		msg.Send(c)
	} else {
		//反序列化到数组中
		if key == KEY_BOOK_IN_REDIS {
			books := make([]bookSet.Book, 0)
			for i, _ := range result {
				bookStr := result[i]
				book := bookSet.Book{}
				json.Unmarshal([]byte(bookStr), &book)
				books = append(books, book)
			}
			//设置到消息类中
			msg := modle.Message{ErrCode: modle.MESSAGE_CODE_QUERY_SUCCESS,
				Error: "",
				Data:  books}
			msg.Send(c)
		} else if key == KEY_BOOK_DETAIL_IN_REDIS {
			bookDetails := make([]bookSet.BookDetail, 0)
			for i, _ := range result {
				bookDetailStr := result[i]
				bookDetail := bookSet.BookDetail{}
				json.Unmarshal([]byte(bookDetailStr), &bookDetail)
				bookDetails = append(bookDetails, bookDetail)
			}
			//设置到消息类中
			msg := modle.Message{ErrCode: modle.MESSAGE_CODE_QUERY_SUCCESS,
				Error: "",
				Data:  bookDetails}
			msg.Send(c)
		}

	}
}