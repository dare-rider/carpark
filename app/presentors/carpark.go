package presentors

import "github.com/dare-rider/carpark/app/models/carpark"

type NearestCarparkResponse struct {
	Address        string  `json:"address"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	TotalSlots     int     `json:"total_slots"`
	AvailableSlots int     `json:"available_slots"`
}

func (rp *NearestCarparkResponse) SerializeFromModel(mod *carpark.Model) *NearestCarparkResponse {
	rp.Address = mod.Address
	rp.Latitude = mod.Latitude
	rp.Longitude = mod.Longitude
	totalSlots := 0
	availableSlots := 0
	for _, cpInfo := range mod.CarparkInfos {
		totalSlots += cpInfo.TotalLots
		availableSlots += cpInfo.LotsAvailable
	}
	rp.TotalSlots = totalSlots
	rp.AvailableSlots = availableSlots
	return rp
}
