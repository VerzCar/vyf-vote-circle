package app

import (
	"github.com/VerzCar/vyf-vote-circle/api/model"
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

		v, ok := ctx.Get("clientChan")
		if !ok {
			s.log.Error("Did not find client Channel in context")
			ctx.JSON(http.StatusBadRequest, errResponse)
			return
		}

		clientChan, ok := v.(ClientChan)

		if !ok {
			s.log.Error("Did not could cast client Channel")
			ctx.JSON(http.StatusBadRequest, errResponse)
			return
		}

		go func() {
			for {
				rankings, err := s.rankingSubscriptionService.Rankings(ctx.Request.Context(), rankingsReq.CircleID)

				if err != nil {
					s.log.Errorf("service error: %v", err)
					ctx.JSON(http.StatusInternalServerError, errResponse)
					return
				}

				s.serverEventService.MessageSub(rankings)
			}
		}()

		ctx.Stream(
			func(w io.Writer) bool {
				if msg, ok := <-clientChan; ok {
					ctx.SSEvent("message", msg)
					return true
				}
				return false
			},
		)

	}
}
