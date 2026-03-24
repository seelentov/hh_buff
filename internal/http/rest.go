package http

import (
	"errors"
	"hh_buff/internal/models"
	"hh_buff/internal/repo"
	"hh_buff/pkg/hh"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type RestController struct {
	queryRepo    *repo.DBQueryRepo
	snapshotRepo *repo.DBSnapshotRepo
	hhClient     *hh.Client
}

func NewRestController(queryRepo *repo.DBQueryRepo, snapshotRepo *repo.DBSnapshotRepo, hhClient *hh.Client) *RestController {
	return &RestController{queryRepo, snapshotRepo, hhClient}
}

// Data godoc
// @Summary      Получить исторические данные
// @Description  Возвращает срезы данных по ID запросов и диапазону дат
// @Tags         data
// @Produce      json
// @Param        start_date  query     string  false  "Дата начала (YYYY-MM-DD)"
// @Param        end_date    query     string  false  "Дата конца (YYYY-MM-DD)"
// @Param        queries     query     []int   true   "Список ID запросов" collectionFormat(multi)
// @Success      200         {array}   models.DBSnapshot
// @Failure      400         {object}  ErrorResponse
// @Failure      500         {object}  ErrorResponse
// @Router       /data [get]
func (rc *RestController) Data(ctx *gin.Context) {
	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")

	var startDate *time.Time
	var endDate *time.Time

	if startDateStr != "" {
		sd, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid start_date format. Use YYYY-MM-DD"})
			return
		}
		startDate = &sd
	}

	if endDateStr != "" {
		ed, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid end_date format. Use YYYY-MM-DD"})
			return
		}

		endDate = &ed
	}

	queriesRaw := ctx.QueryArray("queries")
	var queryIDs []uint

	for _, s := range queriesRaw {
		id, err := strconv.ParseUint(s, 10, 32)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "queries must be a list of integers"})
			return
		}
		queryIDs = append(queryIDs, uint(id))
	}

	qs, err := rc.snapshotRepo.GetByQueryIDsAndDate(queryIDs, startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(200, qs)
}

// Queries godoc
// @Summary      Получить список запросов
// @Description  Возвращает список запросов
// @Tags         queries
// @Produce      json
// @Success      200         {array}   models.DBQuery
// @Failure      400         {object}  ErrorResponse
// @Failure      500         {object}  ErrorResponse
// @Router       /queries [get]
func (rc *RestController) Queries(ctx *gin.Context) {
	qs, err := rc.queryRepo.GetAll()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(200, qs)
}

type UploadQueryReq struct {
	Name  string                 `binding:"required"`
	Query hh.GetVacanciesRequest `binding:"required"`
}

// UploadQuery godoc
// @Summary      Создать новый запрос
// @Tags         queries
// @Accept       json
// @Produce      json
// @Param        request  body      UploadQueryReq  true  "Данные запроса"
// @Success      200      {object}  MessageResponse
// @Failure      400         {object}  ErrorResponse
// @Failure      500         {object}  ErrorResponse
// @Router       /queries [post]
func (rc *RestController) UploadQuery(ctx *gin.Context) {
	var q UploadQueryReq
	if err := ctx.ShouldBindJSON(&q); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	qu := models.DBQuery{
		Name:  q.Name,
		Query: q.Query,
	}
	err := rc.queryRepo.Save(&qu)
	if err != nil {
		if errors.Is(err, repo.ErrAlreadyExists) {
			ctx.JSON(http.StatusConflict, ErrorResponse{Error: err.Error()})
			return
		}

		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	res, err := rc.hhClient.GetVacancies(q.Query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
	}

	err = rc.snapshotRepo.Save(&models.DBSnapshot{
		QueryID: qu.ID,
		Count:   res.Found,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
	}

	ctx.JSON(http.StatusOK, MessageResponse{Message: "created"})
}
