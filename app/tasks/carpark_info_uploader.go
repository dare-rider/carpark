package tasks

import (
	"github.com/dare-rider/carpark/app/extservices/govsgcarpark"
	"github.com/dare-rider/carpark/app/models/carparkinfo"
	"log"
	"sync"
)

type carparkInfoUploader struct {
	govSgCarparkUc govsgcarpark.Usecase
	carparkInfoUc  carparkinfo.Usecase
}

type CarparkInfoUploader interface {
	Upload() error
}

const (
	carparkInfoUploadBatchSize = 100
)

func NewCarparkInfoUploader(govSgCarparkUc govsgcarpark.Usecase, carparkInfoUc carparkinfo.Usecase) CarparkInfoUploader {
	return &carparkInfoUploader{
		govSgCarparkUc: govSgCarparkUc,
		carparkInfoUc:  carparkInfoUc,
	}
}

func (upldr *carparkInfoUploader) Upload() error {
	cpInfos, err := upldr.govSgCarparkUc.CarparkInfos()
	if err != nil {
		return err
	}
	// initiating sync wait group
	wg := new(sync.WaitGroup)
	for start := 0; start < len(cpInfos); start += carparkInfoUploadBatchSize {
		end := start + carparkInfoUploadBatchSize
		if end > len(cpInfos) {
			end = len(cpInfos)
		}
		wg.Add(1)
		go func(recordsSubset []carparkinfo.Model) {
			for _, record := range recordsSubset {
				err := upldr.carparkInfoUc.InsertOrUpdateByCarParkNo(&record)
				if err != nil {
					log.Println(err)
				}
			}
			// marking each waitgroup done
			wg.Done()
		}(cpInfos[start:end])
	}
	wg.Wait()
	return nil
}
