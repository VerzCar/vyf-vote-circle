package app

import (
	"github.com/VerzCar/vyf-vote-circle/api/model"
	"github.com/gin-gonic/gin"
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

		response := model.Response{
			Status: model.ResponseSuccess,
			Msg:    "",
			Data:   rankings,
		}

		ctx.JSON(http.StatusOK, response)
	}
}
