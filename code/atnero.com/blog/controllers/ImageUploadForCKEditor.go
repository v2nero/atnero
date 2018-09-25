package controllers

import (
	blogSess "atnero.com/blog/models/session"
	"fmt"
	"github.com/astaxie/beego"
	"os"
	"time"
)

type ImageUploadForCKEditorController struct {
	beego.Controller
}

func (this *ImageUploadForCKEditorController) GetImageMaxSize() int64 {
	var maxSize int64 = 0
	for {
		if !blogSess.UserHasRight(&this.Controller, "upload_image_100K") {
			break
		}
		maxSize = 100 * 1024

		if !blogSess.UserHasRight(&this.Controller, "upload_image_400K") {
			break
		}
		maxSize = 400 * 1024

		if !blogSess.UserHasRight(&this.Controller, "upload_image_1M") {
			break
		}
		maxSize = 1024 * 1024

		if !blogSess.UserHasRight(&this.Controller, "upload_image_10M") {
			break
		}
		maxSize = 4 * 1024 * 1024

		break
	}
	return maxSize
}

func (this *ImageUploadForCKEditorController) GenValidFileName() (string, string, error) {
	var fileName string
	var filePath string
	_, userId, err := blogSess.GetUserBaseInfo(&this.Controller)
	if err != nil {
		return fileName, filePath, err
	}
	strUploadImageDir := beego.AppConfig.String("ckeditor::image_upload_dir")
	for {
		tNow := time.Now()
		fileName = fmt.Sprintf("%d_%x_%x",
			userId, tNow.Unix(), tNow.Nanosecond())
		filePath = fmt.Sprintf("%s/%s", strUploadImageDir, fileName)
		_, err := os.Stat(filePath)
		if err != nil && os.IsNotExist(err) {
			break
		}
		continue
	}

	return fileName, filePath, nil
}

func (this *ImageUploadForCKEditorController) Post() {
	result := struct {
		Uploaded int    `json:"uploaded"`
		FileName string `json:"filename"`
		Url      string `json:"url"`
		Error    string `json:"error"`
	}{
		Uploaded: 0,
		FileName: "",
		Url:      "",
		Error:    "",
	}

	for {
		if !blogSess.Logined(&this.Controller) {
			result.Error = "没有权限"
			break
		}

		maxSize := this.GetImageMaxSize()

		f, h, err := this.GetFile("upload")
		if err != nil {
			result.Error = "获取上传文件信息失败"
			break
		}
		defer f.Close()

		if h.Size > maxSize {
			result.Error = "上传文件过大"
			break
		}

		fileName, filePath, err := this.GenValidFileName()
		if err != nil {
			result.Error = "上传失败"
			break
		}
		strUploadImageViewDir := beego.AppConfig.String("ckeditor::image_upload_view_dir")
		err = this.SaveToFile("upload", filePath)
		if err != nil {
			result.Error = "上传失败"
			break
		}

		result.Uploaded = 1
		result.FileName = h.Filename
		result.Url = fmt.Sprintf("%s/%s", strUploadImageViewDir, fileName)
		break
	}
	this.Data["json"] = &result
	this.ServeJSON()
}
