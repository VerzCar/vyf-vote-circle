package app

import (
	"github.com/VerzCar/vyf-vote-circle/api/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) CreateVote() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		errResponse := model.Response{
			Status: model.ResponseError,
			Msg:    "vote cannot be created",
			Data:   false,
		}

		circleReq := &model.CircleUriRequest{}

		err := ctx.ShouldBindUri(circleReq)

		if err != nil {
			s.log.Error(err)
			ctx.JSON(http.StatusBadRequest, errResponse)
			return
		}

		voteCreateReq := &model.VoteCreateRequest{}

		err = ctx.ShouldBindJSON(voteCreateReq)

		if err != nil {
			s.log.Error(err)
			ctx.JSON(http.StatusBadRequest, errResponse)
			return
		}

		if err := s.validate.Struct(voteCreateReq); err != nil {
			s.log.Warn(err)
			ctx.JSON(http.StatusBadRequest, errResponse)
			return
		}

		result, err := s.voteService.CreateVote(ctx.Request.Context(), circleReq.CircleID, voteCreateReq)

		if err != nil {
			s.log.Errorf("service error: %v", err)
			ctx.JSON(http.StatusInternalServerError, errResponse)
			return
		}

		response := model.Response{
			Status: model.ResponseSuccess,
			Msg:    "",
			Data:   result,
		}

		ctx.JSON(http.StatusOK, response)
	}
}
