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
	"net/http"
)

func main() {
	r := bottle.New()
	r.GET("/", func(c *bottle.Context) {
		c.HTML(http.StatusOK, "<h1>Hello bottle</h1>")
	})
	r.GET("/hello", func(c *bottle.Context) {
		// expect /hello?name=jin
		c.Text(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.GET("/hello/:name", func(c *bottle.Context) {
		// expect /hello/jin
		c.Text(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})

	r.GET("/assets/*filepath", func(c *bottle.Context) {
		c.JSON(http.StatusOK, bottle.H{"filepath": c.Param("filepath")})
	})

	r.POST("/login", func(c *bottle.Context) {
		c.JSON(http.StatusOK, bottle.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	r.Run(":9999")
}