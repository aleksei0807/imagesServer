package main

import (
	"log"
	"path"
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
	}

	addr := viper.GetString("address")

	var frontendOrigins []string
	if viper.IsSet("frontendOrigins") {
		frontendOrigins = viper.GetStringSlice("frontendOrigins")
	}

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
		servepath := path.Clean(mapRoute.GetString("servepath"))
		savepath := path.Clean(mapRoute.GetString("savepath"))
		fullpath := mapRoute.GetString("fullpath")
		multiple := false
		if mapRoute.IsSet("multiple") {
			multiple = mapRoute.GetBool("multiple")
		}
		rename := true
		if mapRoute.IsSet("rename") {
			rename = mapRoute.GetBool("rename")
		}
		fileserve := path.Clean(mapRoute.GetString("fileserve")) + "/:route"
		r.POST(servepath, saveHandler(savepath, fullpath, multiple, rename, frontendOrigins))
		r.GET(fileserve, fasthttp.FSHandler(savepath, 2))
	}

	s := &fasthttp.Server{
		MaxRequestBodySize: 2<<29, // Harry Potter and magic number
		Handler:            r.Handler,
	}

	s.ListenAndServe(addr)
}
