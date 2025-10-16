package services

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/adtoba/grinbid-backend/src/models"
)

type PaystackService struct {
	SecretKey string
}

func NewPaystackService(secretKey string) *PaystackService {
	return &PaystackService{SecretKey: secretKey}
}

func (ps *PaystackService) InitializeTransaction(amountNaira string, email string, transaction models.Transaction) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"amount":   amountNaira,
		"email":    email,
		"metadata": transaction,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	url := "https://api.paystack.co/transaction/initialize"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+ps.SecretKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response["data"].(map[string]interface{}), nil
}
