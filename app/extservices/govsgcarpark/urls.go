package govsgcarpark

import "github.com/dare-rider/carpark/utils"

func (rp *repo) carparkAvailabilitylURL() string {
	return utils.JoinURL(rp.baseURL, "/v1/transport/carpark-availability")
}
