package main

/*
(1)
$ curl -i http://localhost:9999/
HTTP/1.1 200 OK
Date: Mon, 12 Aug 2019 16:52:52 GMT
Content-Length: 18
Content-Type: text/html; charset=utf-8
<h1>Hello bottle</h1>

(2)
$ curl "http://localhost:9999/hello?name=jin"
hello jin, you're at /hello

(3)
$ curl "http://localhost:9999/login" -X POST -d 'username=jinli&password=1234'
{"password":"1234","username":"jinli"}

(4)
$ curl "http://localhost:9999/xxx"
404 NOT FOUND: /xxx
*/

import (
	"bottle"
	"log"
	"net/http"
	"time"
)

func onlyForV2() bottle.HandlerFunc {
	return func(c *bottle.Context) {
		// Start timer
		t := time.Now()
		// if a server error occurred
		c.Fail(500, "Internal Server Error")
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Request.RequestURI, time.Since(t))
	}
}

func main() {
	app := bottle.New()
	app.Use(bottle.Logger())
	app.LoadHTMLGlob("templates/*")

	v1 := app.Group("/v1")
	{
		v1.Static("/assets", "static")

		v1.GET("/", func(c *bottle.Context) {
			c.HTML(http.StatusOK, "css.html", nil)
		})

		v1.GET("/hello", func(c *bottle.Context) {
			// expect /hello?name=jin
			c.Text(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})


	}

	v2 := app.Group("/v2")
	v2.Use(onlyForV2()) // v2 group middleware
	{
		v2.GET("/hello/:name", func(c *bottle.Context) {
			c.Text(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})

		v2.POST("/login", func(c *bottle.Context) {
			c.JSON(http.StatusOK, bottle.D{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})
	}

	app.GET("/hello", func(c *bottle.Context) {
		// expect /hello?name=jin
		c.Text(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	app.GET("/hello/:name", func(c *bottle.Context) {
		// expect /hello/jin
		c.Text(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})

	app.POST("/login", func(c *bottle.Context) {
		c.JSON(http.StatusOK, bottle.D{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	app.Run(":9999")
}