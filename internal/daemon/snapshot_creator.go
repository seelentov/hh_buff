package daemon

import (
	"hh_buff/internal/models"
	"hh_buff/internal/repo"
	"hh_buff/pkg/hh"
	"log"
	"sync"
	"time"
)

type SnapshotCreatorDaemon struct {
	queryRepo    *repo.DBQueryRepo
	snapshotRepo *repo.DBSnapshotRepo
	hhClient     *hh.Client

	interval time.Duration
	period   time.Duration

	stop chan struct{}
	wg   sync.WaitGroup
}

func NewSnapshotCreator(
	queryRepo *repo.DBQueryRepo,
	snapshotRepo *repo.DBSnapshotRepo,
	hhClient *hh.Client,
	interval time.Duration,
	period time.Duration,
) *SnapshotCreatorDaemon {
	return &SnapshotCreatorDaemon{
		queryRepo:    queryRepo,
		snapshotRepo: snapshotRepo,
		hhClient:     hhClient,
		interval:     interval,
		period:       period,
		stop:         make(chan struct{}),
	}
}

func (d *SnapshotCreatorDaemon) Start() {
	d.wg.Add(1)
	go func() {
		defer d.wg.Done()
		d.run()
	}()
}

func (d *SnapshotCreatorDaemon) Stop() {
	close(d.stop)
	d.wg.Wait()
	log.Println("SnapshotCreatorDaemon stopped safely")
}

func (d *SnapshotCreatorDaemon) run() {
	ticker := time.NewTicker(d.interval)
	defer ticker.Stop()

	log.Printf("SnapshotCreatorDaemon started with interval %v and period %v", d.interval, d.period)
	d.processIteration()

	for {
		select {
		case <-ticker.C:
			d.processIteration()
		case <-d.stop:
			return
		}
	}
}

func (d *SnapshotCreatorDaemon) processIteration() {
	queries, err := d.findAllQueries()
	if err != nil {
		log.Printf("Error fetching queries: %v", err)
		return
	}

	for _, query := range queries {
		count, err := d.getVacanciesCount(query)
		if err != nil {
			log.Printf("Error getting count for query %s: %v", query.Name, err)
			return
		}

		if err := d.createSnapshot(query, count); err != nil {
			log.Printf("Error saving snapshot for query %s: %v", query.Name, err)
		}

		log.Printf("Snapshot created: %s :%v", query.Name, count)
		time.Sleep(d.period)
	}

	log.Printf("Iteration finished: processed %d queries", len(queries))
}

func (d *SnapshotCreatorDaemon) findAllQueries() ([]*models.DBQuery, error) {
	return d.queryRepo.GetAll()
}

func (d *SnapshotCreatorDaemon) getVacanciesCount(query *models.DBQuery) (int, error) {
	res, err := d.hhClient.GetVacancies(query.Query)
	if err != nil {
		return 0, err
	}
	return res.Found, nil
}

func (d *SnapshotCreatorDaemon) createSnapshot(query *models.DBQuery, vacanciesCount int) error {
	return d.snapshotRepo.Save(&models.DBSnapshot{
		QueryID: query.ID,
		Count:   vacanciesCount,
	})
}
