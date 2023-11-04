package app

import (
	"fmt"
	"github.com/VerzCar/vyf-vote-circle/api/model"
	"github.com/VerzCar/vyf-vote-circle/utils"
	"github.com/gin-gonic/gin"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
)

func (s *Server) Circle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		errResponse := model.Response{
			Status: model.ResponseError,
			Msg:    "cannot find circle",
			Data:   nil,
		}

		circleReq := &model.CircleUriRequest{}

		err := ctx.ShouldBindUri(circleReq)

		if err != nil {
			s.log.Error(err)
			ctx.JSON(http.StatusBadRequest, errResponse)
			return
		}

		circle, err := s.circleService.Circle(ctx.Request.Context(), circleReq.CircleID)

		if err != nil {
			s.log.Errorf("service error: %v", err)
			ctx.JSON(http.StatusInternalServerError, errResponse)
			return
		}

		var voters []*model.CircleVoterResponse

		for _, voter := range circle.Voters {
			voterResponse := &model.CircleVoterResponse{
				ID:         voter.ID,
				Voter:      voter.Voter,
				Commitment: voter.Commitment,
				VotedFor:   voter.VotedFor,
				VotedFrom:  voter.VotedFrom,
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

func (s *Server) EligibleToBeInCircle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		errResponse := model.Response{
			Status: model.ResponseError,
			Msg:    "user is not eligible to be in circle",
			Data:   nil,
		}

		circleReq := &model.CircleUriRequest{}

		err := ctx.ShouldBindUri(circleReq)

		if err != nil {
			s.log.Error(err)
			ctx.JSON(http.StatusBadRequest, errResponse)
			return
		}

		eligible, err := s.circleService.EligibleToBeInCircle(ctx.Request.Context(), circleReq.CircleID)

		if err != nil {
			s.log.Errorf("service error: %v", err)
			ctx.JSON(http.StatusInternalServerError, errResponse)
			return
		}

		if !eligible {
			ctx.JSON(http.StatusForbidden, errResponse)
			return
		}

		response := model.Response{
			Status: model.ResponseSuccess,
			Msg:    "",
			Data:   eligible,
		}

		ctx.JSON(http.StatusOK, response)
	}
}

func (s *Server) Circles() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		errResponse := model.Response{
			Status: model.ResponseError,
			Msg:    "cannot find circles",
			Data:   nil,
		}

		circles, err := s.circleService.Circles(ctx.Request.Context())

		if err != nil {
			s.log.Errorf("service error: %v", err)
			ctx.JSON(http.StatusInternalServerError, errResponse)
			return
		}

		if circles == nil {
			response := model.Response{
				Status: model.ResponseSuccess,
				Msg:    "Has no circles",
				Data:   []*model.Circle{},
			}
			ctx.JSON(http.StatusNoContent, response)
			return
		}

		var circlesResponse []*model.CircleResponse

		for _, circle := range circles {
			var voters []*model.CircleVoterResponse

			for _, voter := range circle.Voters {
				voterResponse := &model.CircleVoterResponse{
					ID:         voter.ID,
					Voter:      voter.Voter,
					Commitment: voter.Commitment,
					VotedFor:   voter.VotedFor,
					VotedFrom:  voter.VotedFrom,
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

			circlesResponse = append(circlesResponse, circleResponse)
		}

		response := model.Response{
			Status: model.ResponseSuccess,
			Msg:    "",
			Data:   circlesResponse,
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

		if err := s.validate.Struct(circleCreateReq); err != nil {
			s.log.Warn(err)
			ctx.JSON(http.StatusBadRequest, errResponse)
			return
		}

		circle, err := s.circleService.CreateCircle(ctx.Request.Context(), circleCreateReq)

		if err != nil {
			s.log.Errorf("service error: %v", err)
			ctx.JSON(http.StatusInternalServerError, errResponse)
			return
		}

		var voters []*model.CircleVoterResponse

		for _, voter := range circle.Voters {
			voterResponse := &model.CircleVoterResponse{
				ID:         voter.ID,
				Voter:      voter.Voter,
				Commitment: voter.Commitment,
				VotedFor:   voter.VotedFor,
				VotedFrom:  voter.VotedFrom,
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

func (s *Server) UpdateCircle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		errResponse := model.Response{
			Status: model.ResponseError,
			Msg:    "circle cannot be updated",
			Data:   nil,
		}

		circleUpdateReq := &model.CircleUpdateRequest{}

		err := ctx.ShouldBindJSON(circleUpdateReq)

		if err != nil {
			s.log.Error(err)
			ctx.JSON(http.StatusBadRequest, errResponse)
			return
		}

		if err := s.validate.Struct(circleUpdateReq); err != nil {
			s.log.Warn(err)
			ctx.JSON(http.StatusBadRequest, errResponse)
			return
		}

		circle, err := s.circleService.UpdateCircle(ctx.Request.Context(), circleUpdateReq)

		if err != nil {
			s.log.Errorf("service error: %v", err)
			ctx.JSON(http.StatusInternalServerError, errResponse)
			return
		}

		var voters []*model.CircleVoterResponse

		for _, voter := range circle.Voters {
			voterResponse := &model.CircleVoterResponse{
				ID:         voter.ID,
				Voter:      voter.Voter,
				Commitment: voter.Commitment,
				VotedFor:   voter.VotedFor,
				VotedFrom:  voter.VotedFrom,
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

func (s *Server) AddToGlobalCircle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		errResponse := model.Response{
			Status: model.ResponseError,
			Msg:    "user cannot be added to global circle",
			Data:   nil,
		}

		err := s.circleService.AddToGlobalCircle(ctx.Request.Context())

		if err != nil {
			s.log.Errorf("service error: %v", err)
			ctx.JSON(http.StatusInternalServerError, errResponse)
			return
		}

		response := model.Response{
			Status: model.ResponseSuccess,
			Msg:    "Added to global circle",
			Data:   nil,
		}

		ctx.JSON(http.StatusOK, response)
	}
}

func (s *Server) UploadCircleImage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		errResponse := model.Response{
			Status: model.ResponseError,
			Msg:    "cannot upload file",
			Data:   nil,
		}

		circleReq := &model.CircleUriRequest{}

		err := ctx.ShouldBindUri(circleReq)

		if err != nil {
			s.log.Error(err)
			ctx.JSON(http.StatusBadRequest, errResponse)
			return
		}

		multiPartFile, err := ctx.FormFile("circleImageFile")

		if err != nil {
			s.log.Errorf("service error: %v", err)
			ctx.JSON(http.StatusInternalServerError, errResponse)
			return
		}

		contentFile, err := multiPartFile.Open()
		if err != nil {
			s.log.Errorf("service error: %v", err)
			ctx.JSON(http.StatusInternalServerError, errResponse)
			return
		}

		defer contentFile.Close()

		bytes, err := io.ReadAll(contentFile)
		if err != nil {
			s.log.Errorf("service error: %v", err)
			ctx.JSON(http.StatusInternalServerError, errResponse)
			return
		}

		mimeType := http.DetectContentType(bytes)

		if !utils.IsImageMimeType(mimeType) {
			s.log.Errorf("file type is wrong type: %s", mimeType)
			errResponse.Msg = "file type is not an image"
			ctx.JSON(http.StatusNotAcceptable, errResponse)
			return
		}

		_, _ = contentFile.Seek(0, 0)

		circleImage, _, err := image.Decode(contentFile)

		st := utils.Resize(circleImage, image.Point{X: 200, Y: 200})

		filePath := fmt.Sprintf("circle/image/%d/%s", circleReq.CircleID, "main")

		_, err = s.extStorageService.Upload(
			ctx.Request.Context(),
			filePath,
			st,
		)

		if err != nil {
			s.log.Errorf("service error: %v", err)
			ctx.JSON(http.StatusInternalServerError, errResponse)
			return
		}

		imageEndpoint := fmt.Sprintf("%s/%s", s.extStorageService.ObjectEndpoint(), filePath)

		updateCircleReq := &model.CircleUpdateRequest{
			ImageSrc: &imageEndpoint,
		}

		circle, err := s.circleService.UpdateCircle(ctx.Request.Context(), updateCircleReq)

		if err != nil {
			s.log.Errorf("service error: %v", err)
			ctx.JSON(http.StatusInternalServerError, errResponse)
			return
		}

		response := model.Response{
			Status: model.ResponseSuccess,
			Msg:    "",
			Data:   circle.ImageSrc,
		}

		ctx.JSON(http.StatusOK, response)
	}
}
