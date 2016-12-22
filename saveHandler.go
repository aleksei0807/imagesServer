package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"strings"

	"github.com/valyala/fasthttp"
)

func saveHandler(savepath string, fullpath string, multiple bool, rename bool) func(*fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("Access-Control-Allow-Origin", strings.TrimRight(string(ctx.Request.Header.Referer()), "/"))
		ctx.Response.Header.Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE")
		ctx.Response.Header.Set("Access-Control-Allow-Headers", "Content-Type, *")
		ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")

		multipart, err := ctx.MultipartForm()
		if err != nil {
			ctx.Error("Bad Request", 400)
			return
		}
		files := multipart.File["imageFiles"]
		result := "["
		for k, fileHeader := range files {
			originalFilename := fileHeader.Filename
			fileName := originalFilename
			file, err := fileHeader.Open()
			if err != nil {
				ctx.Error("Bad Request", 400)
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
					ctx.Error("Bad Request", 400)
					return
				}

				hashInBytes := hash.Sum(nil)[:16]

				fileHash := hex.EncodeToString(hashInBytes)

				idx := strings.LastIndex(fileName, ".")
				if idx >= 0 && len(fileName) > idx {
					fileName = fileHash + fileName[idx:]
				} else {
					fileName = fileHash
				}
			}

			filePath := savepath + "/" + fileName
			servePath := fullpath + "/" + fileName

			if _, err := os.Stat(filePath); err == nil {
				result += "\"" + servePath + "\""
				if k < len(files)-1 {
					result += ","
				}
				file.Close()
				continue
			}

			newFile, err := os.Create(filePath)
			if err != nil {
				ctx.Error("Server Error", 500)
				return
			}

			_, err = newFile.Write(fileData)
			if err != nil {
				ctx.Error("Server Error", 500)
				return
			}

			result += "\"" + servePath + "\""
			if k < len(files)-1 {
				result += ","
			}

			file.Close()
		}

		result += "]"
		ctx.SetContentType("application/json")
		ctx.WriteString(result)
	}
}
