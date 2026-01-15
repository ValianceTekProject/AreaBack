package download

import "github.com/gin-gonic/gin"

func DownloadApk(ctx *gin.Context) {
	ctx.File("./client.apk")
}
