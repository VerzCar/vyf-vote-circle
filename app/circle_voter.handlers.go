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

		err = ctx.ShouldBind(circleVotersReq)

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

		filterBy := &model.CircleVotersFilterBy{
			HasBeenVoted: circleVotersReq.HasBeenVoted,
		}

		voters, userVoter, err := s.circleVoterService.CircleVotersFiltered(
			ctx.Request.Context(),
			circleReq.CircleID,
			filterBy,
		)

		if err != nil {
			s.log.Errorf("service error: %v", err)
			ctx.JSON(http.StatusInternalServerError, errResponse)
			return
		}

		votersRes := make([]*model.CircleVoterResponse, 0)

		for _, voter := range voters {
			voterResponse := &model.CircleVoterResponse{
				ID:         voter.ID,
				Voter:      voter.Voter,
				Commitment: voter.Commitment,
				VotedFor:   voter.VotedFor,
				CreatedAt:  voter.CreatedAt,
				UpdatedAt:  voter.UpdatedAt,
			}
			votersRes = append(votersRes, voterResponse)
		}

		var userVoterRes *model.CircleVoterResponse

		if userVoter != nil {
			userVoterRes = &model.CircleVoterResponse{
				ID:         userVoter.ID,
				Voter:      userVoter.Voter,
				Commitment: userVoter.Commitment,
				VotedFor:   userVoter.VotedFor,
				CreatedAt:  userVoter.CreatedAt,
				UpdatedAt:  userVoter.UpdatedAt,
			}
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

func (s *Server) CircleVoterJoinCircle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		errResponse := model.Response{
			Status: model.ResponseError,
			Msg:    "cannot join as voter in circle",
			Data:   nil,
		}

		circleReq := &model.CircleUriRequest{}

		err := ctx.ShouldBindUri(circleReq)

		if err != nil {
			s.log.Error(err)
			ctx.JSON(http.StatusBadRequest, errResponse)
			return
		}

		voter, err := s.circleVoterService.CircleVoterJoinCircle(
			ctx.Request.Context(),
			circleReq.CircleID,
		)

		if err != nil {
			s.log.Errorf("service error: %v", err)
			ctx.JSON(http.StatusInternalServerError, errResponse)
			return
		}

		voterRes := &model.CircleVoterResponse{
			ID:         voter.ID,
			Voter:      voter.Voter,
			Commitment: voter.Commitment,
			VotedFor:   voter.VotedFor,
			CreatedAt:  voter.CreatedAt,
			UpdatedAt:  voter.UpdatedAt,
		}

		response := model.Response{
			Status: model.ResponseSuccess,
			Msg:    "",
			Data:   voterRes,
		}

		ctx.JSON(http.StatusOK, response)
	}
}

func (s *Server) CircleVoterLeaveCircle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		errResponse := model.Response{
			Status: model.ResponseError,
			Msg:    "cannot leave as voter from circle",
			Data:   nil,
		}

		circleReq := &model.CircleUriRequest{}

		err := ctx.ShouldBindUri(circleReq)

		if err != nil {
			s.log.Error(err)
			ctx.JSON(http.StatusBadRequest, errResponse)
			return
		}

		err = s.circleVoterService.CircleVoterLeaveCircle(
			ctx.Request.Context(),
			circleReq.CircleID,
		)

		if err != nil {
			s.log.Errorf("service error: %v", err)
			ctx.JSON(http.StatusInternalServerError, errResponse)
			return
		}

		response := model.Response{
			Status: model.ResponseSuccess,
			Msg:    "",
			Data:   "",
		}

		ctx.JSON(http.StatusOK, response)
	}
}

func (s *Server) CircleVoterAddToCircle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		errResponse := model.Response{
			Status: model.ResponseError,
			Msg:    "cannot add voter to circle",
			Data:   nil,
		}

		circleReq := &model.CircleUriRequest{}

		err := ctx.ShouldBindUri(circleReq)

		if err != nil {
			s.log.Error(err)
			ctx.JSON(http.StatusBadRequest, errResponse)
			return
		}

		circleVoterReq := &model.CircleVoterRequest{}

		err = ctx.ShouldBindJSON(circleVoterReq)

		if err != nil {
			s.log.Error(err)
			ctx.JSON(http.StatusBadRequest, errResponse)
			return
		}

		if err := s.validate.Struct(circleVoterReq); err != nil {
			s.log.Warn(err)
			ctx.JSON(http.StatusBadRequest, errResponse)
			return
		}

		voter, err := s.circleVoterService.CircleVoterAddToCircle(
			ctx.Request.Context(),
			circleReq.CircleID,
			circleVoterReq,
		)

		if err != nil {
			s.log.Errorf("service error: %v", err)
			ctx.JSON(http.StatusInternalServerError, errResponse)
			return
		}

		voterRes := &model.CircleVoterResponse{
			ID:         voter.ID,
			Voter:      voter.Voter,
			Commitment: voter.Commitment,
			VotedFor:   voter.VotedFor,
			CreatedAt:  voter.CreatedAt,
			UpdatedAt:  voter.UpdatedAt,
		}

		response := model.Response{
			Status: model.ResponseSuccess,
			Msg:    "",
			Data:   voterRes,
		}

		ctx.JSON(http.StatusOK, response)
	}
}

func (s *Server) CircleVoterRemoveFromCircle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		errResponse := model.Response{
			Status: model.ResponseError,
			Msg:    "cannot add voter to circle",
			Data:   nil,
		}

		circleReq := &model.CircleUriRequest{}

		err := ctx.ShouldBindUri(circleReq)

		if err != nil {
			s.log.Error(err)
			ctx.JSON(http.StatusBadRequest, errResponse)
			return
		}

		circleVoterReq := &model.CircleVoterRequest{}

		err = ctx.ShouldBindJSON(circleVoterReq)

		if err != nil {
			s.log.Error(err)
			ctx.JSON(http.StatusBadRequest, errResponse)
			return
		}

		if err := s.validate.Struct(circleVoterReq); err != nil {
			s.log.Warn(err)
			ctx.JSON(http.StatusBadRequest, errResponse)
			return
		}

		voter, err := s.circleVoterService.CircleVoterRemoveFromCircle(
			ctx.Request.Context(),
			circleReq.CircleID,
			circleVoterReq,
		)

		if err != nil {
			s.log.Errorf("service error: %v", err)
			ctx.JSON(http.StatusInternalServerError, errResponse)
			return
		}

		voterRes := &model.CircleVoterResponse{
			ID:         voter.ID,
			Voter:      voter.Voter,
			Commitment: voter.Commitment,
			VotedFor:   voter.VotedFor,
			CreatedAt:  voter.CreatedAt,
			UpdatedAt:  voter.UpdatedAt,
		}

		response := model.Response{
			Status: model.ResponseSuccess,
			Msg:    "",
			Data:   voterRes,
		}

		ctx.JSON(http.StatusOK, response)
	}
}
