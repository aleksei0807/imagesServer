package main

import (
	"log"
	"github.com/valyala/fasthttp"
)

func saveHandler(savepath string, fullpath string, multiple bool, notrename bool) func(*fasthttp.RequestCtx) {
	log.Printf("savepath: %s", savepath)
	log.Printf("fullpath: %s", fullpath)
	log.Printf("multiple: %t", multiple)
	log.Printf("notrename: %t", notrename)
	return func(ctx *fasthttp.RequestCtx) {
		
	}
}
