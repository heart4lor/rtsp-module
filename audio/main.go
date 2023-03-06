package main

import (
	"github.com/aler9/rtsp-simple-server/internal/core"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/exec"
)

type PublishRequest struct {
	NetworkMode string
	RtspServer  string // needed when networkMode is public
	Codec       string
	Bitrate     string
}

type SubscribeRequest struct {
	Username    string // username you want to listen
	NetworkMode string
	RtspServer  string // needed when networkMode is public
}

var user string
var port string

func publish(c *gin.Context) {
	var publishRequest PublishRequest
	if err := c.BindJSON(&publishRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"msg": "操作不合法，请检查"})
		return
	}
	if publishRequest.NetworkMode == "local" {
		publishRequest.RtspServer = "localhost"
	}
	if publishRequest.Codec == "" {
		publishRequest.Codec = "libmp3lame"
	}
	if publishRequest.Bitrate == "" {
		publishRequest.Bitrate = "32k"
	}
	rtspUrl := "rtsp://" + publishRequest.RtspServer + ":8554/audio/" + user
	// ffmpeg -f avfoundation -i ":0" -acodec libmp3lame -ab 32k -ac 1 -f rtsp rtsp://hk1.sunyongfei.cn:8554/audio/syf
	publishCommand := exec.Command(
		"ffmpeg", "-f", "avfoundation", "-i", ":0", "-acodec", publishRequest.Codec, "-ab", publishRequest.Bitrate, "-ac", "1", "-f", "rtsp", rtspUrl)
	go func() {
		err := publishCommand.Run()
		if err != nil {
			print("error while publishing stream, cmd:", publishCommand.Args)
		}
	}()
	c.IndentedJSON(http.StatusOK, gin.H{
		"success": true,
		"model":   publishCommand.Args,
	})
}

func subscribe(c *gin.Context) {
	var subscribeRequest SubscribeRequest
	if err := c.BindJSON(&subscribeRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"msg": "操作不合法，请检查"})
		return
	}
	rtspUrl := "rtsp://" + subscribeRequest.RtspServer + ":8554/audio/" + subscribeRequest.Username
	// ffplay rtsp://hk1.sunyongfei.cn:8554/audio/syf
	subscribeCommand := exec.Command("ffplay", rtspUrl)
	go func() {
		err := subscribeCommand.Run()
		if err != nil {
			print("error while publishing stream, cmd:", subscribeCommand.Args)
		}
	}()
	c.IndentedJSON(http.StatusOK, gin.H{
		"success": true,
		"model":   subscribeCommand.Args,
	})
}

func main() {
	// New allocates a core.
	user = os.Args[1]
	port = os.Args[2]
	go core.New(nil)
	router := gin.Default()
	router.POST("/publish", publish)
	router.POST("/subscribe", subscribe)
	err := router.Run("0.0.0.0:" + port)
	if err != nil {
		return
	}
}
