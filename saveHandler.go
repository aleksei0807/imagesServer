package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"strings"

	"github.com/valyala/fasthttp"
	"github.com/kirillDanshin/myutils"
)

/* this is not craziness, this is optimizations */
const (
	ACAO = "Access-Control-Allow-Origin"
	ACAM = "Access-Control-Allow-Methods"
	methods = "GET,PUT,POST,DELETE"
	ACAH = "Access-Control-Allow-Headers"
	CType = "Content-Type, *"
	ACAC = "Access-Control-Allow-Credentials"
	trueStr = "true"
	imageFiles = "imageFiles"
	lBr = "["
	rBr = "]"
	badReq = "Bad Request"
	serverErr = "Internal Server Error"
	dot = "."
	slash = "/"
	quote = "\""
	coma = ","
	jsonType = "application/json"
)

func saveHandler(
	savepath string,
	fullpath string,
	multiple bool,
	rename bool,
	frontendOrigins []string,
) func(*fasthttp.RequestCtx) {

	origin := "*"
	if len(frontendOrigins) > 0 {
		origin = ""
		for k, v := range frontendOrigins {
			origin = myutils.Concat(origin, v)
			if k < len(frontendOrigins) - 1 {
				origin = myutils.Concat(origin, ", ")
			}
		}
	}

	return func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set(ACAO, origin)
		ctx.Response.Header.Set(ACAM, methods)
		ctx.Response.Header.Set(ACAH, CType)
		ctx.Response.Header.Set(ACAC, trueStr)

		multipart, err := ctx.MultipartForm()
		if err != nil {
			ctx.Error(badReq, fasthttp.StatusBadRequest)
			return
		}
		files := multipart.File[imageFiles]
		result := lBr
		for k, fileHeader := range files {
			originalFilename := fileHeader.Filename
			fileName := originalFilename
			file, err := fileHeader.Open()
			if err != nil {
				ctx.Error(badReq, fasthttp.StatusBadRequest)
				return
			}

			var fileData []byte
			fdb := new(bytes.Buffer)
			io.Copy(fdb, file)
			fileData = fdb.Bytes()

			defer file.Close()

			if rename != false {
				hash := md5.New()

				_, err := hash.Write(fileData)
				if err != nil {
					ctx.Error(badReq, fasthttp.StatusBadRequest)
					return
				}

				hashInBytes := hash.Sum(nil)[:16]

				fileHash := hex.EncodeToString(hashInBytes)

				idx := strings.LastIndex(fileName, dot)
				if idx >= 0 && len(fileName) > idx {
					fileName = fileHash + fileName[idx:]
				} else {
					fileName = fileHash
				}
			}

			filePath := myutils.Concat(savepath, slash, fileName)
			servePath := myutils.Concat(fullpath, slash, fileName)

			if _, err := os.Stat(filePath); err == nil {
				result = myutils.Concat(result, quote, servePath, quote)
				if k < len(files)-1 {
					result = myutils.Concat(result, coma)
				}
				file.Close()
				continue
			}

			newFile, err := os.Create(filePath)
			if err != nil {
				ctx.Error(serverErr, fasthttp.StatusInternalServerError)
				return
			}

			_, err = newFile.Write(fileData)
			if err != nil {
				ctx.Error(serverErr, fasthttp.StatusInternalServerError)
				return
			}

			result = myutils.Concat(result, quote, servePath, quote)
			if k < len(files)-1 {
				result = myutils.Concat(result, coma)
			}

			file.Close()
		}

		result = myutils.Concat(result, rBr)
		ctx.SetContentType(jsonType)
		ctx.WriteString(result)
	}
}
