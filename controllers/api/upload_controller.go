package api

import (
	"io/ioutil"

	"github.com/kataras/iris"
	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"

	"github.com/mlogclub/mlog/services"
	"github.com/mlogclub/mlog/utils/oss"
)

const uploadMaxBytes int64 = 1024 * 1024 * 3 // 1M

type UploadController struct {
	Ctx iris.Context
}

func (this *UploadController) Post() *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(this.Ctx)
	if user == nil {
		return simple.Error(simple.ErrorNotLogin)
	}

	file, header, err := this.Ctx.FormFile("image")
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	defer file.Close()

	if header.Size > uploadMaxBytes {
		return simple.ErrorMsg("图片不能超过3M")
	}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}

	logrus.Info("上传文件：", header.Filename, " size:", header.Size)

	url, err := oss.UploadImage(fileBytes)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	return simple.NewEmptyRspBuilder().Put("url", url).JsonResult()
}
