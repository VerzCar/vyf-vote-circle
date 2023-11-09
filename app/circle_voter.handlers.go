package app

import (
	"github.com/VerzCar/vyf-vote-circle/api/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) CircleVoters() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		errResponse := model.Response{
			Status: model.ResponseError,
			Msg:    "cannot find circle voters",
			Data:   nil,
		}

		circleReq := &model.CircleUriRequest{}

		err := ctx.ShouldBindUri(circleReq)

		if err != nil {
			s.log.Error(err)
			ctx.JSON(http.StatusBadRequest, errResponse)
			return
		}

		circleVotersReq := &model.CircleVotersRequest{}

		err = ctx.ShouldBindJSON(circleVotersReq)

		if err != nil {
			s.log.Error(err)
			ctx.JSON(http.StatusBadRequest, errResponse)
			return
		}

		if err := s.validate.Struct(circleVotersReq); err != nil {
			s.log.Warn(err)
			ctx.JSON(http.StatusBadRequest, errResponse)
			return
		}

		voters, userVoter, err := s.circleVoterService.CircleVotersFiltered(
			ctx.Request.Context(),
			circleReq.CircleID,
			&circleVotersReq.Filter,
		)

		if err != nil {
			s.log.Errorf("service error: %v", err)
			ctx.JSON(http.StatusInternalServerError, errResponse)
			return
		}

		var votersRes []*model.CircleVoterResponse

		for _, voter := range voters {
			voterResponse := &model.CircleVoterResponse{
				ID:         voter.ID,
				Voter:      voter.Voter,
				Commitment: voter.Commitment,
				VotedFor:   voter.VotedFor,
				VotedFrom:  voter.VotedFrom,
				CreatedAt:  voter.CreatedAt,
				UpdatedAt:  voter.UpdatedAt,
			}
			votersRes = append(votersRes, voterResponse)
		}

		userVoterRes := &model.CircleVoterResponse{
			ID:         userVoter.ID,
			Voter:      userVoter.Voter,
			Commitment: userVoter.Commitment,
			VotedFor:   userVoter.VotedFor,
			VotedFrom:  userVoter.VotedFrom,
			CreatedAt:  userVoter.CreatedAt,
			UpdatedAt:  userVoter.UpdatedAt,
		}

		circleVotersRes := &model.CircleVotersResponse{
			Voters:    votersRes,
			UserVoter: userVoterRes,
		}

		response := model.Response{
			Status: model.ResponseSuccess,
			Msg:    "",
			Data:   circleVotersRes,
		}

		ctx.JSON(http.StatusOK, response)
	}
}
