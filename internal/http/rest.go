package http

import (
	"hh_buff/internal/models"
	"hh_buff/internal/repo"
	"hh_buff/pkg/hh"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type RestController struct {
	queryRepo    *repo.DBQueryRepo
	snapshotRepo *repo.DBSnapshotRepo
}

func NewRestController(queryRepo *repo.DBQueryRepo, snapshotRepo *repo.DBSnapshotRepo) *RestController {
	return &RestController{queryRepo, snapshotRepo}
}

func (rc *RestController) Data(ctx *gin.Context) {
	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")

	var startDate *time.Time
	var endDate *time.Time

	if startDateStr != "" {
		sd, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format. Use YYYY-MM-DD"})
			return
		}
		startDate = &sd
	}

	if endDateStr != "" {
		ed, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format. Use YYYY-MM-DD"})
			return
		}

		endDate = &ed
	}

	queriesRaw := ctx.QueryArray("queries")
	var queryIDs []uint

	for _, s := range queriesRaw {
		id, err := strconv.ParseUint(s, 10, 32)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "queries must be a list of integers"})
			return
		}
		queryIDs = append(queryIDs, uint(id))
	}

	qs, err := rc.snapshotRepo.GetByQueryIDsAndDate(queryIDs, startDate, endDate)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, qs)
}

func (rc *RestController) Current(ctx *gin.Context) {
	qs, err := rc.queryRepo.GetAll()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data := make(map[string]int)
	var wg sync.WaitGroup
	for _, q := range qs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			count, err := rc.snapshotRepo.GetCurrentCount(q.ID)
			if err != nil {
				return
			}
			data[q.Name] = count
		}()
	}
	wg.Wait()

	ctx.JSON(200, data)
}

func (rc *RestController) Queries(ctx *gin.Context) {
	qs, err := rc.queryRepo.GetAll()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, qs)
}

type UploadQueryReq struct {
	Name  string                 `binding:"required"`
	Query hh.GetVacanciesRequest `binding:"required"`
}

func (rc *RestController) UploadQuery(ctx *gin.Context) {
	var q UploadQueryReq
	if err := ctx.ShouldBindJSON(&q); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	created, err := rc.queryRepo.Save(&models.DBQuery{
		Name:  q.Name,
		Query: q.Query,
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"created": created})
}
