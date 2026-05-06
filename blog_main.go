package main

import "github.com/DaiXing/go_Blog/blogx"

func main() {
	blogx.ConfigLogger()
	blogx.ConfigLoadParams()

	blogx.DbInit()

	blogx.WebInit()

}
