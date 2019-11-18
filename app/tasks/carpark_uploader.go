package tasks

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/dare-rider/carpark/app/models/carpark"
	"github.com/dare-rider/carpark/constant"
	"github.com/dare-rider/carpark/utils"
	"github.com/dare-rider/carpark/utils/svy21"
	"log"
	"os"
	"strconv"
	"sync"
)

type carparkUploader struct {
	carparkUc       carpark.Usecase
	csvFileBasePath string
}

type CarparkUploader interface {
	Upload() error
}

func NewCarparkUploader(carparkUc carpark.Usecase, csvFileBasePath string) CarparkUploader {
	return &carparkUploader{
		carparkUc:       carparkUc,
		csvFileBasePath: csvFileBasePath,
	}
}

const (
	carparkFileName = "carpark.csv"
	// configurable batch size to increase/decrease parallelism
	carparkUploadBatchSize = 100
)

func (upldr *carparkUploader) Upload() error {

	filePath := utils.JoinURL(upldr.csvFileBasePath, carparkFileName)
	csvFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	reader := csv.NewReader(csvFile)
	defer csvFile.Close()
	allRecords, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// initiating sync wait group
	wg := new(sync.WaitGroup)
	for start := 0; start < len(allRecords); start += carparkUploadBatchSize {
		end := start + carparkUploadBatchSize
		if end > len(allRecords) {
			end = len(allRecords)
		}
		wg.Add(1)
		go func(recordsSubset [][]string) {
			for _, record := range recordsSubset {
				err := upldr.addCarpark(record)
				if err != nil {
					log.Println(err)
				}
			}
			// marking each waitgroup done
			wg.Done()
		}(allRecords[start:end])
	}
	wg.Wait()
	return nil
}

func (upldr *carparkUploader) addCarpark(line []string) error {
	if len(line) != 12 {
		return errors.New(fmt.Sprintf("%s, %v", constant.UploadRowInvalid, line))
	}
	cpMod := carpark.Model{
		CarParkNo:           line[0],
		Address:             line[1],
		CarParkType:         line[4],
		TypeOfParkingSystem: line[5],
		ShortTermParking:    line[6],
		FreeParking:         line[7],
		NightParking:        utils.StringToBool(line[8]),
		CarParkBasement:     utils.StringToBool(line[11]),
	}
	cpMod.XCoord, _ = strconv.ParseFloat(line[2], 64)
	cpMod.YCoord, _ = strconv.ParseFloat(line[3], 64)
	cpMod.Latitude, cpMod.Longitude = svy21.ToLatLon(cpMod.XCoord, cpMod.YCoord)
	cpMod.CarParkDecks, _ = strconv.Atoi(line[9])
	cpMod.GantryHeight, _ = strconv.ParseFloat(line[10], 64)

	err := upldr.carparkUc.InsertOrUpdate(&cpMod)
	return err
}
