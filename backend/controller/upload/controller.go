package upload

import (
	"github.com/gin-gonic/gin"

	"github.com/yycy134679/school-secondhand-trading-system/backend/common/errors"
	"github.com/yycy134679/school-secondhand-trading-system/backend/common/resp"
	"github.com/yycy134679/school-secondhand-trading-system/backend/common/util"
)

// UploadController 处理文件上传
type UploadController struct{}

// NewUploadController 创建上传控制器
func NewUploadController() *UploadController {
	return &UploadController{}
}

// UploadImage 接收图片文件并返回可访问 URL
func (uc *UploadController) UploadImage(c *gin.Context) {
	if _, exists := c.Get("user_id"); !exists {
		resp.Error(c, errors.CodeUnauthenticated, "用户未登录")
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		resp.Error(c, errors.CodeInvalidParams, "请上传有效的图片文件")
		return
	}
	defer file.Close()

	url, err := util.SaveImage(file, header)
	if err != nil {
		resp.Error(c, errors.CodeInvalidParams, err.Error())
		return
	}

	resp.Success(c, gin.H{
		"url": url,
	})
}
