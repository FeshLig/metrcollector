package agent

import (
	"fmt"
	"log"
	"net/http"

	"github.com/FeshLig/metrcollector/internal/handler"
)

type MetricsSender interface {
	Send(metricType, name, value string) error
}

type HTTPSender struct {
	BaseURL string
	Client  *http.Client
}

func NewSender(url string, client *http.Client) *HTTPSender {
	return &HTTPSender{
		BaseURL: url,
		Client:  client,
	}
}

func (h *HTTPSender) Send(metricType, name, value string) error {
	url := fmt.Sprintf("%s/update/%s/%s/%s", h.BaseURL, metricType, name, value)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/plain")

	resp, err := h.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("неверный статус: %d", resp.StatusCode)
	}

	return nil
}

func (h *HTTPSender) SendMetrics(storage handler.SnapshotMetrics) {

	var strValue string

	sender := h

	for name, value := range storage.SnapshotGauges() {
		strValue = fmt.Sprintf("%f", value)
		if err := sender.Send("gauge", name, strValue); err != nil {
			log.Printf("Ошибка отправки gauge %s: %v", name, err)
		}
	}

	for name, value := range storage.SnapshotCounters() {
		strValue = fmt.Sprintf("%d", value)
		if err := sender.Send("counter", name, strValue); err != nil {
			log.Printf("Ошибка отправки counter %s: %v", name, err)
		}
	}
}
