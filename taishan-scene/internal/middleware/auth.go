package middleware

import (
	"github.com/gin-gonic/gin"
)

type ResponseData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func SessionAuthMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 在白名单中添加需要忽略中间件保护的路由
		whitelist := []string{""}
		for _, path := range whitelist {
			if ctx.Request.URL.Path == path {
				ctx.Next()
				return
			}
		}
		//client := &http.Client{}
		// 创建 GET 请求
		//request, err := http.NewRequest("GET", conf.Conf.Url.Account+"/user/info", nil)
		//if err != nil {
		//	log.Logger.Error(http.StatusInternalServerError, err)
		//	ctx.AbortWithStatusJSON(http.StatusOK, ResponseData{
		//		Code:    errno.ErrToken,
		//		Message: errno.CodeAlertMap[errno.ErrToken],
		//	})
		//	return
		//}
		//openUserId := ctx.Request.Header.Get("Open-User-Id")
		//if openUserId != "" {
		//	// 自定义用户
		//	request.Header.Set("Open-User-Id", openUserId)
		//} else {
		//	// 普通用户
		//	auth := ctx.Request.Header.Get("Authorization")
		//	request.Header.Set("Authorization", auth)
		//}
		//response, err := client.Do(request)
		//if err != nil {
		//	log.Logger.Error("middleware.auth.SessionAuthMiddleWare.clientDo(): ", http.StatusInternalServerError, err)
		//	ctx.AbortWithStatusJSON(http.StatusOK, ResponseData{
		//		Code:    errno.ErrToken,
		//		Message: errno.CodeAlertMap[errno.ErrToken],
		//	})
		//	return
		//}
		//defer response.Body.Close()
		//
		//// 检查响应是否成功
		//if response.StatusCode != http.StatusOK {
		//	log.Logger.Error("middleware.auth.SessionAuthMiddleWare.responseStatusCode(): ", response.StatusCode)
		//	ctx.AbortWithStatusJSON(http.StatusOK, ResponseData{
		//		Code:    errno.ErrToken,
		//		Message: errno.CodeAlertMap[errno.ErrToken],
		//	})
		//	return
		//}
		//// 解析响应体
		//var result rao.CommonResponse[UserInfo]
		//if err = json.NewDecoder(response.Body).Decode(&result); err != nil {
		//	// 处理响应体解析错误
		//	log.Logger.Error("middleware.auth.SessionAuthMiddleWare.Decode(): ", http.StatusInternalServerError, err)
		//	ctx.AbortWithStatusJSON(http.StatusOK, ResponseData{
		//		Code:    errno.ErrToken,
		//		Message: errno.CodeAlertMap[errno.ErrToken],
		//	})
		//	return
		//}
		//if result.Code != (0) {
		//	log.Logger.Error("middleware.auth.SessionAuthMiddleWare.resultCode not 0: ", http.StatusInternalServerError, err)
		//	ctx.AbortWithStatusJSON(http.StatusOK, ResponseData{
		//		Code:    errno.ErrToken,
		//		Message: errno.CodeAlertMap[errno.ErrToken],
		//	})
		//	return
		//}
		//
		//id := result.Data.ID
		//name := result.Data.Name
		//if !(id > 0) || name == "" {
		//	log.Logger.Error("middleware.auth.SessionAuthMiddleWare.error id or name: ", http.StatusInternalServerError, err)
		//	ctx.AbortWithStatusJSON(http.StatusOK, ResponseData{
		//		Code:    errno.ErrToken,
		//		Message: errno.CodeAlertMap[errno.ErrToken],
		//	})
		//	return
		//}
		id = int32(1)
		name := "简单随风"
		// 将用户ID和用户名写入上下文中
		ctx.Set("userID", id)
		ctx.Set("userName", name)

		ctx.Next()
	}
}

type UserInfo struct {
	AvatarURL   string      `json:"avatarUrl"`
	Channel     int         `json:"channel"`
	CreatedTime string      `json:"created_time"`
	Email       string      `json:"email"`
	FsUserID    string      `json:"fs_user_id"`
	ID          int32       `json:"id"`
	IsDeleted   int         `json:"isDeleted"`
	Name        string      `json:"name"`
	Nickname    string      `json:"nickname"`
	OpenID      string      `json:"open_id"`
	Phone       string      `json:"phone"`
	Roles       interface{} `json:"roles"`
	Status      int         `json:"status"`
}
