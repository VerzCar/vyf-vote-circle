package app

import (
	"github.com/VerzCar/vyf-vote-circle/api/model"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

func (s *Server) Rankings() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		errResponse := model.Response{
			Status: model.ResponseError,
			Msg:    "cannot find rankings",
			Data:   false,
		}

		rankingsReq := &model.RankingsUriRequest{}

		err := ctx.ShouldBindUri(rankingsReq)

		if err != nil {
			s.log.Error(err)
			ctx.JSON(http.StatusBadRequest, errResponse)
			return
		}

		if err := s.validate.Struct(rankingsReq); err != nil {
			s.log.Warn(err)
			ctx.JSON(http.StatusBadRequest, errResponse)
			return
		}

		rankings, err := s.rankingService.Rankings(ctx.Request.Context(), rankingsReq.CircleID)

		if err != nil {
			s.log.Errorf("service error: %v", err)
			ctx.JSON(http.StatusInternalServerError, errResponse)
			return
		}

		response := model.Response{
			Status: model.ResponseSuccess,
			Msg:    "",
			Data:   rankings,
		}

		ctx.JSON(http.StatusOK, response)
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// TODO check msg events - not working currently

func (s *Server) RankingsSubscription() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		for {
			rankings, _ := s.rankingSubscriptionService.Rankings(ctx.Request.Context(), 4)
			conn.WriteJSON(rankings)
			time.Sleep(time.Second)
		}
	}
}
