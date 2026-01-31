package sweego

import "fmt"

type SweegoDomainListInformation struct {
	Id                   int64  `json:"id"`
	ClientId             int64  `json:"client_id"`
	Uuid                 string `json:"uuid"`
	CreationDate         string `json:"creation_dt"`
	LastVerificationDate string `json:"last_verification_dt"`
	TrackingOpenEnabled  bool   `json:"tracking_open_enabled"`
	TrackingClickEnabled bool   `json:"tracking_click_enabled"`
	IsVerified           bool   `json:"is_verified"`
	Domain               string `json:"domain"`
}

type SweegoDomainRecord struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Data     string `json:"data"`
	Verified bool   `json:"verified"`
}

type SweegoDomainDetails struct {
	Uuid                 string               `json:"uuid"`
	IsVerified           bool                 `json:"is_verified"`
	TrackingOpenEnabled  bool                 `json:"tracking_open_enabled"`
	TrackingClickEnabled bool                 `json:"tracking_click_enabled"`
	Domain               string               `json:"domain"`
	DomainRecord         SweegoDomainRecord   `json:"domain_record"`
	DkimRecord           SweegoDomainRecord   `json:"dkim_record"`
	DmarcRecord          SweegoDomainRecord   `json:"dmarc_record"`
	InboundRecordList    []SweegoDomainRecord `json:"inbound_record_list"`
	TrackingRecord       SweegoDomainRecord   `json:"tracking_record"`
}

func (api *SweegoApi) ListDomains() ([]SweegoDomainListInformation, error) {
	api.logger.Debug("ListDomains")

	var response []SweegoDomainListInformation
	err := api.executeGetRequest(fmt.Sprintf("clients/%s/domains", api.clientId), &response)
	return response, err
}

func (api *SweegoApi) GetDomain(uuid string) (SweegoDomainDetails, error) {
	api.logger.Debug(fmt.Sprintf("ListDomains(%#v)", uuid))

	var response SweegoDomainDetails
	err := api.executeGetRequest(fmt.Sprintf("clients/%s/domains/%s", api.clientId, uuid), &response)
	return response, err
}

func (api *SweegoApi) CreateDomain(domain string) (SweegoDomainDetails, error) {
	api.logger.Debug(fmt.Sprintf("CreateDomain(%#v)", domain))

	var response SweegoDomainDetails
	err := api.executeJsonRequest(
		"POST",
		fmt.Sprintf("clients/%s/domains", api.clientId),
		map[string]string{"domain": domain},
		&response,
	)

	return response, err
}

func (api *SweegoApi) DeleteDomain(uuid string) error {
	api.logger.Debug(fmt.Sprintf("DeleteDomain(%#v)", uuid))

	return api.executePlainRequest("DELETE", fmt.Sprintf("clients/%s/domains/%s", api.clientId, uuid), nil)
}
