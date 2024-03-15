package app

import (
	"github.com/VerzCar/vyf-vote-circle/api/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) UserOption() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		errResponse := model.Response{
			Status: model.ResponseError,
			Msg:    "cannot find user option",
			Data:   nil,
		}

		optionResponse, err := s.userOptionService.UserOption(ctx.Request.Context())

		if err != nil {
			s.log.Errorf("service error: %v", err)
			ctx.JSON(http.StatusInternalServerError, errResponse)
			return
		}

		response := model.Response{
			Status: model.ResponseSuccess,
			Msg:    "",
			Data:   optionResponse,
		}

		ctx.JSON(http.StatusOK, response)
	}
}
