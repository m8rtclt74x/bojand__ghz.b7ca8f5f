package database

import (
	"sync/atomic"

	"github.com/bojand/ghz/web/model"
)

// ListAllDetailsForReport lists all details for report
func (d *Database) ListAllDetailsForReport(rid uint) ([]*model.Detail, error) {
	r := &model.Report{}
	r.ID = rid

	s := make([]*model.Detail, 0)

	err := d.DB.Model(r).Related(&s).Error

	return s, err
}

// CreateDetailsBatch creates a batch of details
// Returns the number successfully created, and the number failed
func (d *Database) CreateDetailsBatch(rid uint, s []*model.Detail) (uint, uint) {
	nReq := len(s)

	NC := 10

	var nErr uint32

	sem := make(chan bool, NC)

	var nCreated, errCount uint

	for _, item := range s {
		sem <- true

		go func(detail *model.Detail) {
			defer func() { <-sem }()

			detail.ReportID = rid
			err := d.createDetail(detail)

			if err != nil {
				atomic.AddUint32(&nErr, 1)
			}
		}(item)
	}

	for i := 0; i < cap(sem)-1; i++ {
		sem <- true
	}

	errCount = uint(atomic.LoadUint32(&nErr))
	if errCount > 0 {
		nCreated = uint(nReq) - errCount - 1
	} else {
		nCreated = uint(nReq) - errCount
	}

	return nCreated, errCount
}

func (d *Database) createDetail(detail *model.Detail) error {
	return d.DB.Create(detail).Error
}
