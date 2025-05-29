package exchangerate

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"io"
	"maxwellzp/golang-db-api-context/pkg/config"
	"maxwellzp/golang-db-api-context/pkg/models"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
}

func NewClient(cfg config.ApiConfig) *Client {
	return &Client{
		httpClient: &http.Client{
			// HTTP client timeout (transport layer)
			// For transport-layer safety net
			Timeout: cfg.Timeout,
		},
		baseURL: cfg.URL,
		apiKey:  cfg.ApiKey,
	}
}

type Response struct {
	Result             string             `json:"result"`
	TimeLastUpdateUnix int64              `json:"time_last_update_unix"`
	ConversionRates    map[string]float64 `json:"conversion_rates"`
}

func (c *Client) GetExchangeRates(ctx context.Context, baseCurrency string) ([]models.ExchangeRate, error) {
	// The stricter timeout wins. HTTP Client Timeout VS Context Timeout
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/%s/latest/%s", c.baseURL, c.apiKey, baseCurrency),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	if response.Result != "success" {
		return nil, fmt.Errorf("error result: %s", response.Result)
	}

	rates := make([]models.ExchangeRate, 0, len(response.ConversionRates))
	dateUpdated := time.Unix(response.TimeLastUpdateUnix, 0)
	for currencyCode, rate := range response.ConversionRates {
		rateD := decimal.NewFromFloat(rate)

		rates = append(rates, models.ExchangeRate{
			CurrencyCode:     currencyCode,
			BaseCurrencyCode: baseCurrency,
			Rate:             rateD,
			DateUpdated:      dateUpdated,
		})
	}

	return rates, nil
}
