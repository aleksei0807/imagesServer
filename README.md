# imagesServer

Server for [react-images-uploader](https://github.com/aleksei0807/react-images-uploader)

## Usage

Download [release](https://github.com/aleksei0807/imagesServer/releases) and change config.yaml

## Config

- address - server address `example: :9090`
- frontendOrigins - origins your frontends for Access-Control-Allow-Origin `default: *`
- routes - object with your routes
  - servepath - serve path
  - savepath - path where you want to save images
  - fullpath - full path to folder with images
  - fileserve - serve images path
  - multiple - allows to upload a bunch of images `default: false`
  - rename - if false, then do not rename image `default: true`

### Config example

```yaml
address: :9090
frontendOrigins:
    - http://localhost:8181
routes:
    multipleFiles:
        servepath: /multiple
        savepath: ./static/multipleFiles
        fullpath: http://localhost:9090/static/multipleFiles
        fileserve: /static/multipleFiles
        multiple: true
        rename: true
    files:
        servepath: /notmultiple
        savepath: ./static/files
        fullpath: http://localhost:9090/static/files
        fileserve: /static/files
        multiple: false
        rename: true
```
