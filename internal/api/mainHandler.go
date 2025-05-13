package api

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/star-find-cloud/star-mall/model"
	"github.com/star-find-cloud/star-mall/pkg/database"
	"net/http"
)

func MainHandler(c *gin.Context) {
	session := sessions.Default(c)
	userinfo, ok := session.Get("userinfo").(model.Admin)
	if ok {
		// 权限查询
		c.HTML(http.StatusOK, "index.html", gin.H{
			"authList": []model.Auth{},
			"userinfo": userinfo.UserName,
			"isSuper":  userinfo.IsSuper,
		})
	} else {
		c.Redirect(http.StatusFound, "/login")
	}
}

func ChangeStatusHandler(c *gin.Context) {
	id := c.PostForm("id")
	if id == "" {
		c.JSON(http.StatusOK, gin.H{"success": false, "msg": "非法请求"})
		return
	}
	table := c.PostForm("table")
	field := c.PostForm("field")
	db := database.GetDB()
	_, err := db.ExecContext(c, "update ? set ? =ABS(? - 1) where id = ?", table, field, field, id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "msg": "更新失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "msg": "更新成功"})
}
