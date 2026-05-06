package main

import "gitee.com/dx/go_blog/blogx"

func main() {
	blogx.ConfigLogger()
	blogx.ConfigLoadParams()

	blogx.DbInit()

	blogx.WebInit()

}
