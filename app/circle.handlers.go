package app

import (
	"github.com/VerzCar/vyf-vote-circle/api/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) Circle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		errResponse := model.Response{
			Status: model.ResponseError,
			Msg:    "cannot find circle",
			Data:   nil,
		}

		circleReq := &model.CircleRequest{}

		err := ctx.ShouldBindJSON(circleReq)

		if err != nil {
			s.log.Error(err)
			ctx.JSON(http.StatusBadRequest, errResponse)
			return
		}

		circle, err := s.circleService.Circle(ctx.Request.Context(), circleReq.ID)

		if err != nil {
			s.log.Errorf("service error: %v", err)
			ctx.JSON(http.StatusInternalServerError, errResponse)
			return
		}

		voters := []*model.CircleVoterResponse{}

		for _, voter := range circle.Voters {
			voterResponse := &model.CircleVoterResponse{
				ID:         voter.ID,
				Voter:      voter.Voter,
				Commitment: voter.Commitment,
				CreatedAt:  voter.CreatedAt,
				UpdatedAt:  voter.UpdatedAt,
			}
			voters = append(voters, voterResponse)
		}

		circleResponse := &model.CircleResponse{
			ID:          circle.ID,
			Name:        circle.Name,
			Description: circle.Description,
			ImageSrc:    circle.ImageSrc,
			Voters:      voters,
			Private:     circle.Private,
			Active:      circle.Active,
			CreatedFrom: circle.CreatedFrom,
			ValidUntil:  circle.ValidUntil,
		}

		response := model.Response{
			Status: model.ResponseSuccess,
			Msg:    "",
			Data:   circleResponse,
		}

		ctx.JSON(http.StatusOK, response)
	}
}

func (s *Server) CreateCircle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		errResponse := model.Response{
			Status: model.ResponseError,
			Msg:    "circle cannot be created",
			Data:   nil,
		}

		circleCreateReq := &model.CircleCreateRequest{}

		err := ctx.ShouldBindJSON(circleCreateReq)

		if err != nil {
			s.log.Error(err)
			ctx.JSON(http.StatusBadRequest, errResponse)
			return
		}

		circle, err := s.circleService.CreateCircle(ctx.Request.Context(), circleCreateReq)

		if err != nil {
			s.log.Errorf("service error: %v", err)
			ctx.JSON(http.StatusInternalServerError, errResponse)
			return
		}

		voters := []*model.CircleVoterResponse{}

		for _, voter := range circle.Voters {
			voterResponse := &model.CircleVoterResponse{
				ID:         voter.ID,
				Voter:      voter.Voter,
				Commitment: voter.Commitment,
				CreatedAt:  voter.CreatedAt,
				UpdatedAt:  voter.UpdatedAt,
			}
			voters = append(voters, voterResponse)
		}

		circleResponse := &model.CircleResponse{
			ID:          circle.ID,
			Name:        circle.Name,
			Description: circle.Description,
			ImageSrc:    circle.ImageSrc,
			Voters:      voters,
			Private:     circle.Private,
			Active:      circle.Active,
			CreatedFrom: circle.CreatedFrom,
			ValidUntil:  circle.ValidUntil,
		}

		response := model.Response{
			Status: model.ResponseSuccess,
			Msg:    "",
			Data:   circleResponse,
		}

		ctx.JSON(http.StatusOK, response)
	}
}
