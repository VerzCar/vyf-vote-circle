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

		circleResponse := &model.CircleResponse{
			ID:          circle.ID,
			Name:        circle.Name,
			Description: circle.Description,
			ImageSrc:    circle.ImageSrc,
			Private:     circle.Private,
			Active:      circle.Active,
			Stage:       circle.Stage,
			CreatedFrom: circle.CreatedFrom,
			ValidFrom:   circle.ValidFrom,
			ValidUntil:  circle.ValidUntil,
			CreatedAt:   circle.CreatedAt,
			UpdatedAt:   circle.UpdatedAt,
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

		circlesResponse := make([]*model.CircleResponse, 0)

		for _, circle := range circles {
			circleResponse := &model.CircleResponse{
				ID:          circle.ID,
				Name:        circle.Name,
				Description: circle.Description,
				ImageSrc:    circle.ImageSrc,
				Private:     circle.Private,
				Active:      circle.Active,
				Stage:       circle.Stage,
				CreatedFrom: circle.CreatedFrom,
				ValidFrom:   circle.ValidFrom,
				ValidUntil:  circle.ValidUntil,
				CreatedAt:   circle.CreatedAt,
				UpdatedAt:   circle.UpdatedAt,
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

func (s *Server) CirclesByName() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		errResponse := model.Response{
			Status: model.ResponseError,
			Msg:    "cannot find circles with name",
			Data:   nil,
		}

		circleUriReq := &model.CircleByUriRequest{}

		err := ctx.ShouldBindUri(circleUriReq)

		if err != nil {
			s.log.Error(err)
			ctx.JSON(http.StatusBadRequest, errResponse)
			return
		}

		if err := s.validate.Struct(circleUriReq); err != nil {
			s.log.Warn(err)
			ctx.JSON(http.StatusBadRequest, errResponse)
			return
		}

		circles, err := s.circleService.CirclesFiltered(ctx.Request.Context(), &circleUriReq.Name)

		if err != nil {
			s.log.Errorf("service error: %v", err)
			ctx.JSON(http.StatusInternalServerError, errResponse)
			return
		}

		paginatedCirclesResponse := make([]*model.CirclePaginatedResponse, 0)

		for _, circle := range circles {
			c := &model.CirclePaginatedResponse{
				ID:          circle.ID,
				Name:        circle.Name,
				Description: circle.Description,
				ImageSrc:    circle.ImageSrc,
				Active:      circle.Active,
				Stage:       circle.Stage,
			}
			paginatedCirclesResponse = append(paginatedCirclesResponse, c)
		}

		response := model.Response{
			Status: model.ResponseSuccess,
			Msg:    "",
			Data:   paginatedCirclesResponse,
		}

		ctx.JSON(http.StatusOK, response)
	}
}

func (s *Server) CirclesOpenCommitments() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		errResponse := model.Response{
			Status: model.ResponseError,
			Msg:    "cannot find open commitments of circles related to user",
			Data:   nil,
		}

		circles, err := s.circleService.CirclesOpenCommitments(ctx.Request.Context())

		if err != nil {
			s.log.Errorf("service error: %v", err)
			ctx.JSON(http.StatusInternalServerError, errResponse)
			return
		}

		paginatedCirclesResponse := make([]*model.CirclePaginatedResponse, 0)

		for _, circle := range circles {
			c := &model.CirclePaginatedResponse{
				ID:              circle.ID,
				Name:            circle.Name,
				Description:     circle.Description,
				ImageSrc:        circle.ImageSrc,
				VotersCount:     &circle.VotersCount,
				CandidatesCount: &circle.CandidatesCount,
				Active:          circle.Active,
				Stage:           circle.Stage,
			}
			paginatedCirclesResponse = append(paginatedCirclesResponse, c)
		}

		response := model.Response{
			Status: model.ResponseSuccess,
			Msg:    "",
			Data:   paginatedCirclesResponse,
		}

		ctx.JSON(http.StatusOK, response)
	}
}

func (s *Server) CirclesOfInterest() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		errResponse := model.Response{
			Status: model.ResponseError,
			Msg:    "cannot find circles related to user",
			Data:   nil,
		}

		circles, err := s.circleService.CirclesOfInterest(ctx.Request.Context())

		if err != nil {
			s.log.Errorf("service error: %v", err)
			ctx.JSON(http.StatusInternalServerError, errResponse)
			return
		}

		paginatedCirclesResponse := make([]*model.CirclePaginatedResponse, 0)

		for _, circle := range circles {
			c := &model.CirclePaginatedResponse{
				ID:              circle.ID,
				Name:            circle.Name,
				Description:     circle.Description,
				ImageSrc:        circle.ImageSrc,
				VotersCount:     &circle.VotersCount,
				CandidatesCount: &circle.CandidatesCount,
				Active:          circle.Active,
				Stage:           circle.Stage,
			}
			paginatedCirclesResponse = append(paginatedCirclesResponse, c)
		}

		response := model.Response{
			Status: model.ResponseSuccess,
			Msg:    "",
			Data:   paginatedCirclesResponse,
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

		circleResponse := &model.CircleResponse{
			ID:          circle.ID,
			Name:        circle.Name,
			Description: circle.Description,
			ImageSrc:    circle.ImageSrc,
			Private:     circle.Private,
			Active:      circle.Active,
			Stage:       circle.Stage,
			CreatedFrom: circle.CreatedFrom,
			ValidFrom:   circle.ValidFrom,
			ValidUntil:  circle.ValidUntil,
			CreatedAt:   circle.CreatedAt,
			UpdatedAt:   circle.UpdatedAt,
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

		circleReq := &model.CircleUriRequest{}

		err := ctx.ShouldBindUri(circleReq)

		if err != nil {
			s.log.Error(err)
			ctx.JSON(http.StatusBadRequest, errResponse)
			return
		}

		circleUpdateReq := &model.CircleUpdateRequest{}

		err = ctx.ShouldBindJSON(circleUpdateReq)

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

		circle, err := s.circleService.UpdateCircle(ctx.Request.Context(), circleReq.CircleID, circleUpdateReq)

		if err != nil {
			s.log.Errorf("service error: %v", err)
			ctx.JSON(http.StatusInternalServerError, errResponse)
			return
		}

		circleResponse := &model.CircleResponse{
			ID:          circle.ID,
			Name:        circle.Name,
			Description: circle.Description,
			ImageSrc:    circle.ImageSrc,
			Private:     circle.Private,
			Active:      circle.Active,
			Stage:       circle.Stage,
			CreatedFrom: circle.CreatedFrom,
			ValidFrom:   circle.ValidFrom,
			ValidUntil:  circle.ValidUntil,
			CreatedAt:   circle.CreatedAt,
			UpdatedAt:   circle.UpdatedAt,
		}

		response := model.Response{
			Status: model.ResponseSuccess,
			Msg:    "",
			Data:   circleResponse,
		}

		ctx.JSON(http.StatusOK, response)
	}
}

func (s *Server) DeleteCircle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		errResponse := model.Response{
			Status: model.ResponseError,
			Msg:    "circle cannot be deleted",
			Data:   nil,
		}

		circleReq := &model.CircleUriRequest{}

		err := ctx.ShouldBindUri(circleReq)

		if err != nil {
			s.log.Error(err)
			ctx.JSON(http.StatusBadRequest, errResponse)
			return
		}

		err = s.circleService.DeleteCircle(ctx.Request.Context(), circleReq.CircleID)

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

		imageSrc, err := s.circleUploadService.UploadImage(ctx.Request.Context(), multiPartFile, circleReq.CircleID)

		if err != nil {
			s.log.Errorf("service error: %v", err)
			ctx.JSON(http.StatusInternalServerError, errResponse)
			return
		}

		response := model.Response{
			Status: model.ResponseSuccess,
			Msg:    "",
			Data:   imageSrc,
		}

		ctx.JSON(http.StatusOK, response)
	}
}

func (s *Server) DeleteCircleImage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		errResponse := model.Response{
			Status: model.ResponseError,
			Msg:    "cannot delete file",
			Data:   nil,
		}

		circleReq := &model.CircleUriRequest{}

		err := ctx.ShouldBindUri(circleReq)

		if err != nil {
			s.log.Error(err)
			ctx.JSON(http.StatusBadRequest, errResponse)
			return
		}

		imageSrc, err := s.circleUploadService.DeleteImage(ctx.Request.Context(), circleReq.CircleID)

		if err != nil {
			s.log.Errorf("service error: %v", err)
			ctx.JSON(http.StatusInternalServerError, errResponse)
			return
		}

		response := model.Response{
			Status: model.ResponseSuccess,
			Msg:    "",
			Data:   imageSrc,
		}

		ctx.JSON(http.StatusOK, response)
	}
}
