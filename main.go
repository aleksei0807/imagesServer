package main

import (
	"log"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"github.com/buaazp/fasthttprouter"
)

func main()  {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
		return
	}

	addr := viper.GetString("addres")
	r := fasthttprouter.New()
	log.Printf("serve on %s", addr)

	routes := viper.GetStringMap("routes")
	for _, route := range routes {
		mapRoute := route.(map[string]interface{})
		servepath := mapRoute["servepath"].(string)
		savepath := mapRoute["savepath"].(string)
		fullpath := mapRoute["fullpath"].(string)
		multipleI := mapRoute["multiple"]
		var multiple bool
		if (multipleI == nil) {
			multiple = false
		} else {
			multiple = multipleI.(bool)
		}
		renameI := mapRoute["notrename"]
		var rename bool
		if (renameI == nil) {
			rename = true
		} else {
			rename = renameI.(bool)
		}
		fileserve := mapRoute["fileserve"].(string) + "/:route"
		r.POST(servepath, saveHandler(savepath, fullpath, multiple, rename))
		r.GET(fileserve, serveHandler)
	}

	fasthttp.ListenAndServe(addr, r.Handler)
}
