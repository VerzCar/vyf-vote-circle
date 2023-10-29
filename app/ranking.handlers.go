package app

import (
	"github.com/VerzCar/vyf-vote-circle/api/model"
	"github.com/VerzCar/vyf-vote-circle/app/router/server_event"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
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

		var rankingsResponse []*model.RankingResponse

		for _, ranking := range rankings {
			rankingResponse := &model.RankingResponse{
				ID:         ranking.ID,
				IdentityID: ranking.IdentityID,
				Number:     ranking.Number,
				Votes:      ranking.Votes,
				Placement:  ranking.Placement,
				CreatedAt:  ranking.CreatedAt,
				UpdatedAt:  ranking.UpdatedAt,
			}
			rankingsResponse = append(rankingsResponse, rankingResponse)
		}

		response := model.Response{
			Status: model.ResponseSuccess,
			Msg:    "",
			Data:   rankingsResponse,
		}

		ctx.JSON(http.StatusOK, response)
	}
}

func (s *Server) RankingsSubscription() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		errResponse := model.Response{
			Status: model.ResponseError,
			Msg:    "cannot stream rankings",
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

		serverEventChan, err := server_event.ContextToServerEventChan[[]*model.Ranking](ctx)

		if err != nil {
			s.log.Error(err)
			ctx.JSON(http.StatusBadRequest, errResponse)
			return
		}

		rankingsObs, err := s.rankingSubscriptionService.RankingsChan(ctx, rankingsReq.CircleID)

		if err != nil {
			s.log.Error(err)
			ctx.JSON(http.StatusBadRequest, errResponse)
			return
		}

		go func() {
			for rankings := range rankingsObs {
				s.rankingsServerEventService.Publish(rankings)
			}
		}()

		ctx.Stream(
			func(w io.Writer) bool {
				if msg, ok := <-serverEventChan; ok {
					ctx.SSEvent("rankings", msg)
					return true
				}
				return false
			},
		)
	}
}
