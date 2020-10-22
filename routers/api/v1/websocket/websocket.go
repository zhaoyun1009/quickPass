package websocket

import (
	"QuickPass/pkg/app"
	"QuickPass/pkg/e"
	"QuickPass/pkg/mq"
	"QuickPass/pkg/util"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type ReqConnectionWebSocket struct {
	Token string `form:"token" json:"token" binding:"required"`
}

func TokenWebSocketHandler(c *gin.Context) {
	var (
		form ReqConnectionWebSocket
	)
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.ErrorResp(c, e.INVALID_PARAMS, err.Error())
		return
	}
	claims, err := util.ParseToken(form.Token)
	if err != nil {
		code := e.ERROR_AUTH_CHECK_TOKEN_FAIL
		switch err.(*jwt.ValidationError).Errors {
		case jwt.ValidationErrorExpired:
			code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
		}
		app.UnauthorizedResp(c, code, "")
		return
	}

	// websocket协议更新
	conn, _, _, err := ws.UpgradeHTTP(c.Request, c.Writer)
	if err != nil {
		app.ErrorResp(c, e.ERROR, err.Error())
		return
	}

	statusInfosChan := make(chan []byte, 5)
	info := &mq.RegisterInfo{
		Agency:   claims.Agency,
		Username: claims.Username,
		Role:     claims.Role,
		DataChan: statusInfosChan,
	}
	mq.SocketChanReceiverInstance.RegisterChan(info)
	// 释放资源
	defer func() {
		// 移除通道
		mq.SocketChanReceiverInstance.RemoveChan(info)
		// 关闭通道
		close(statusInfosChan)
		// 关闭连接
		_ = conn.Close()
	}()
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()
		for {
			// 通道阻塞
			if info, ok := <-statusInfosChan; ok {
				err = wsutil.WriteServerText(conn, info)
				// 关闭通道时err不为空，结束携程
				if err != nil {
					break
				}
			} else {
				break
			}
		}
	}()
	for {
		// 检测客户端连接是否关闭，客户端连接关闭时err不为空
		_, _, err := wsutil.ReadClientData(conn)
		if err != nil {
			break
		}
	}
}
