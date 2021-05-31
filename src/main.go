package main

import (
	"bottle"
	"fmt"
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
	app := bottle.Default()
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

	app.GET("/panic", func(c *bottle.Context) {
		// expect /hello?name=jin
		text := []string{"lorem ipsum"}
		b := text[100]
		fmt.Println(b)

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