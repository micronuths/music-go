package lyric

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"music-go/utils"
)

// 歌词
func Lyric(c *gin.Context) {
	id := c.Query("id")
	params, encSecKey, encErr := utils.Encrypt(fmt.Sprintf(`{"id":"%s","csrf_token":""}`, id))
	if encErr != nil {
		log.Println(encErr)
	}
	res, resErr := utils.DoPostRequest(c.Request.Cookies(), "https://music.163.com/weapi/song/lyric?os=osx&id="+id+"&lv=-1&kv=-1&tv=-1", params, encSecKey)
	if resErr != nil {
		log.Println(resErr)
	}
	c.String(200, res)
}
