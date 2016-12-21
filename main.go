package main

import (
	"log"
	"strings"

	"github.com/buaazp/fasthttprouter"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
)

func main() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
		return
	}

	addr := viper.GetString("address")
	r := fasthttprouter.New()
	log.Printf("serve on %s", addr)

	routes := viper.Sub("routes")
	keys := viper.GetStringMap("routes")
	for key := range keys {
		if strings.Contains(key, ".") {
			continue
		}
		mapRoute := routes.Sub(key)
		if mapRoute == nil {
			log.Fatal("Config is invalid!")
		}
		servepath := mapRoute.GetString("servepath")
		savepath := mapRoute.GetString("savepath")
		fullpath := mapRoute.GetString("fullpath")
		multiple := false
		if mapRoute.IsSet("multiple") {
			multiple = mapRoute.GetBool("multiple")
		}
		rename := true
		if mapRoute.IsSet("rename") {
			rename = mapRoute.GetBool("rename")
		}
		fileserve := mapRoute.GetString("fileserve") + "/:route"
		r.POST(servepath, saveHandler(savepath, fullpath, multiple, rename))
		r.GET(fileserve, fasthttp.FSHandler(savepath, 2))
	}

	s := &fasthttp.Server{
		MaxRequestBodySize: 1024 * 1024 * 1024, // Harry Potter and magic number
		Handler:            r.Handler,
	}

	s.ListenAndServe(addr)
}
