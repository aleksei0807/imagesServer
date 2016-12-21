package main

import (
	"log"
	"github.com/valyala/fasthttp"
)

func serveHandler(ctx *fasthttp.RequestCtx) {
	route := ctx.UserValue("route").(string)
	log.Printf("Serve route: %s", route)
	ctx.WriteString(route)
}
