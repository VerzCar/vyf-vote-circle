package app

import (
	"github.com/VerzCar/vyf-vote-circle/api/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) CircleCandidates() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		errResponse := model.Response{
			Status: model.ResponseError,
			Msg:    "cannot find circle candidates",
			Data:   nil,
		}

		circleReq := &model.CircleUriRequest{}

		err := ctx.ShouldBindUri(circleReq)

		if err != nil {
			s.log.Error(err)
			ctx.JSON(http.StatusBadRequest, errResponse)
			return
		}

		circleCandidatesReq := &model.CircleCandidatesRequest{}

		err = ctx.ShouldBind(circleCandidatesReq)

		if err != nil {
			s.log.Error(err)
			ctx.JSON(http.StatusBadRequest, errResponse)
			return
		}

		if err := s.validate.Struct(circleCandidatesReq); err != nil {
			s.log.Warn(err)
			ctx.JSON(http.StatusBadRequest, errResponse)
			return
		}

		filterBy := &model.CircleCandidatesFilterBy{
			Commitment:   circleCandidatesReq.Commitment,
			HasBeenVoted: circleCandidatesReq.HasBeenVoted,
		}

		candidates, userCandidate, err := s.circleCandidateService.CircleCandidatesFiltered(
			ctx.Request.Context(),
			circleReq.CircleID,
			filterBy,
		)

		if err != nil {
			s.log.Errorf("service error: %v", err)
			ctx.JSON(http.StatusInternalServerError, errResponse)
			return
		}

		votersRes := make([]*model.CircleCandidateResponse, 0)

		for _, candidate := range candidates {
			voterResponse := &model.CircleCandidateResponse{
				ID:         candidate.ID,
				Candidate:  candidate.Candidate,
				Commitment: candidate.Commitment,
				CreatedAt:  candidate.CreatedAt,
				UpdatedAt:  candidate.UpdatedAt,
			}
			votersRes = append(votersRes, voterResponse)
		}

		var userCandidateRes *model.CircleCandidateResponse

		if userCandidate != nil {
			userCandidateRes = &model.CircleCandidateResponse{
				ID:         userCandidate.ID,
				Candidate:  userCandidate.Candidate,
				Commitment: userCandidate.Commitment,
				CreatedAt:  userCandidate.CreatedAt,
				UpdatedAt:  userCandidate.UpdatedAt,
			}
		}

		circleCandidatesRes := &model.CircleCandidatesResponse{
			Candidates:    votersRes,
			UserCandidate: userCandidateRes,
		}

		response := model.Response{
			Status: model.ResponseSuccess,
			Msg:    "",
			Data:   circleCandidatesRes,
		}

		ctx.JSON(http.StatusOK, response)
	}
}

func (s *Server) CircleCandidateCommitment() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		errResponse := model.Response{
			Status: model.ResponseError,
			Msg:    "cannot set commitment for candidate",
			Data:   nil,
		}

		circleReq := &model.CircleUriRequest{}

		err := ctx.ShouldBindUri(circleReq)

		if err != nil {
			s.log.Error(err)
			ctx.JSON(http.StatusBadRequest, errResponse)
			return
		}

		circleCandidateCommitmentReq := &model.CircleCandidateCommitmentRequest{}

		err = ctx.ShouldBindJSON(circleCandidateCommitmentReq)

		if err != nil {
			s.log.Error(err)
			ctx.JSON(http.StatusBadRequest, errResponse)
			return
		}

		if err := s.validate.Struct(circleCandidateCommitmentReq); err != nil {
			s.log.Warn(err)
			ctx.JSON(http.StatusBadRequest, errResponse)
			return
		}

		commitment, err := s.circleCandidateService.CircleCandidateCommitment(
			ctx.Request.Context(),
			circleReq.CircleID,
			circleCandidateCommitmentReq.Commitment,
		)

		if err != nil {
			s.log.Errorf("service error: %v", err)
			ctx.JSON(http.StatusInternalServerError, errResponse)
			return
		}

		response := model.Response{
			Status: model.ResponseSuccess,
			Msg:    "",
			Data:   commitment,
		}

		ctx.JSON(http.StatusOK, response)
	}
}

func (s *Server) CircleCandidateJoinCircle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		errResponse := model.Response{
			Status: model.ResponseError,
			Msg:    "cannot join as candidate in circle",
			Data:   nil,
		}

		circleReq := &model.CircleUriRequest{}

		err := ctx.ShouldBindUri(circleReq)

		if err != nil {
			s.log.Error(err)
			ctx.JSON(http.StatusBadRequest, errResponse)
			return
		}

		candidate, err := s.circleCandidateService.CircleCandidateJoinCircle(
			ctx.Request.Context(),
			circleReq.CircleID,
		)

		if err != nil {
			s.log.Errorf("service error: %v", err)
			ctx.JSON(http.StatusInternalServerError, errResponse)
			return
		}

		candidateRes := &model.CircleCandidateResponse{
			ID:         candidate.ID,
			Candidate:  candidate.Candidate,
			Commitment: candidate.Commitment,
			CreatedAt:  candidate.CreatedAt,
			UpdatedAt:  candidate.UpdatedAt,
		}

		response := model.Response{
			Status: model.ResponseSuccess,
			Msg:    "",
			Data:   candidateRes,
		}

		ctx.JSON(http.StatusOK, response)
	}
}

func (s *Server) CircleCandidateLeaveCircle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		errResponse := model.Response{
			Status: model.ResponseError,
			Msg:    "cannot leave as candidate from circle",
			Data:   nil,
		}

		circleReq := &model.CircleUriRequest{}

		err := ctx.ShouldBindUri(circleReq)

		if err != nil {
			s.log.Error(err)
			ctx.JSON(http.StatusBadRequest, errResponse)
			return
		}

		err = s.circleCandidateService.CircleCandidateLeaveCircle(
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
