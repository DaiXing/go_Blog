package main

import "github.com/DaiXing/go_Blog/blogx"

func main() {
	blogx.InitLogger()

	blogx.LoadParams()

	blogx.DbInit()

	blogx.WebInit()
}
