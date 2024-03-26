package voucher

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"promotion/config"

	"github.com/shopspring/decimal"
)

type VoucherClient struct {
	client           *http.Client
	healthEndPoint   string
	creationEndPoint string
}

func NewVoucherClient(cfg *config.Config) *VoucherClient {
	healthEndPoint := fmt.Sprintf("http://%s:%s/%s", cfg.Voucher.Host, cfg.Voucher.Port, cfg.Voucher.Healthz)
	log.Println("healthEndPoint", healthEndPoint)
	creationEndPoint := fmt.Sprintf("http://%s:%s/%s", cfg.Voucher.Host, cfg.Voucher.Port, cfg.Voucher.EndPoint)
	log.Println("creationEndPoint", creationEndPoint)
	return &VoucherClient{
		client:           &http.Client{},
		healthEndPoint:   healthEndPoint,
		creationEndPoint: creationEndPoint,
	}
}

func (v *VoucherClient) Ping() (bool, error) {
	req, err := http.NewRequest(http.MethodGet, v.healthEndPoint, nil)
	if err != nil {
		return false, err
	}

	resp, err := v.client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	log.Println("Client addr:", v.healthEndPoint, ", ping status:", resp.StatusCode)
	return resp.StatusCode == http.StatusOK, nil
}

func (v *VoucherClient) CreateVoucher(
	ctx context.Context,
	name string,
	campaignExtID string,
	userExtID string,
	percent decimal.Decimal,
) (bool, error) {
	payload := CreateVoucher{
		CampaignExtID: campaignExtID,
		UserExtID:     userExtID,
		Name:          name,
		Value:         percent,
		FixedAmount:   false,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return true, fmt.Errorf("error marshalling JSON: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, v.creationEndPoint, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return true, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := v.client.Do(req)
	if err != nil {
		return true, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusConflict, http.StatusBadRequest, http.StatusBadGateway:
		return true, nil
	case http.StatusCreated, http.StatusOK:
		return false, nil
	default:
		return true, nil
	}
}
