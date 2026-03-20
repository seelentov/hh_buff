package http

import (
	"hh_buff/internal/repo"
	"sync"

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
	qs, err := rc.snapshotRepo.GetAll()
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, qs)
}

func (rc *RestController) Current(ctx *gin.Context) {
	qs, err := rc.queryRepo.GetAll()
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
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
