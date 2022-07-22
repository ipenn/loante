package v1

import (
	"encoding/base64"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io/ioutil"
	"loante/global"
	"loante/service/resp"
	"loante/tools"
	"os"
	"strings"
)

type upload struct { }

func NewUpload() *upload {
	return &upload{}
}

type uploadImg struct {
	Name string `json:"name"`
}

type uploadBase64Img struct {
	Base64 string `json:"base64"`
}


func (s *upload)UploadBase64(c *fiber.Ctx) error {
	input := new(uploadBase64Img)
	if err := c.BodyParser(input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	if len(input.Base64) < 23{
		return resp.Err(c, 1, "数据长度不足")
	}
	path := fmt.Sprintf("static/upload/%s/", tools.GetFormatTime()[0:7])
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return resp.Err(c, 1, "上传失败-1")
	}
	fileName := tools.NewUUID()
	if input.Base64[11]=='j'{
		input.Base64 = input.Base64[23:]
	}else if input.Base64[11]=='p'{
		input.Base64 = input.Base64[22:]
	}else if input.Base64[11]=='g'{
		input.Base64 = input.Base64[22:]
	}
	filePath := fmt.Sprintf("%s%s.%s", path, fileName, "png")
	unbased, err := base64.StdEncoding.DecodeString(input.Base64)
	if err != nil{
		return resp.Err(c, 1, err.Error())
	}
	err = ioutil.WriteFile(filePath, unbased, 0666)
	if err != nil{
		return resp.Err(c, 1, err.Error())
	}
	return resp.OK(c, fmt.Sprintf("%s/%s", global.C.Http.Host, filePath))
}

func (s *upload)Upload(c *fiber.Ctx) error {
	input := new(uploadImg)
	if err := c.BodyParser(input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	path := fmt.Sprintf("static/upload/%s/", tools.GetFormatTime()[0:7])
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return resp.Err(c, 1, "上传失败-1")
	}
	file, _ := c.FormFile("file")
	fileArr := strings.Split(file.Filename, ".")
	if len(fileArr) < 2{
		return resp.Err(c, 1, "上传失败-2")
	}
	if file.Size > 10485760{
		return resp.Err(c, 1, "文件过大")
	}
	ext := fileArr[len(fileArr)-1]
	fileName :=  tools.NewUUID()
	filePath := fmt.Sprintf("%s%s.%s", path, fileName, ext)
	err = c.SaveFile(file, filePath)
	if err != nil {
		return resp.Err(c, 1, "上传失败")
	}
	return resp.OK(c, fmt.Sprintf("%s/%s", global.C.Http.Host, filePath))
}
