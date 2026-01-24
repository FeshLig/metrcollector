package agent

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/FeshLig/metrcollector/internal/handler"
	"github.com/FeshLig/metrcollector/internal/repository"
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
		return fmt.Errorf("%w", err)
	}
	req.Header.Set("Content-Type", "text/plain")

	resp, err := h.Client.Do(req)
	if err != nil {
		return fmt.Errorf("%w", err)
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
		strValue = strconv.FormatFloat(float64(value), 'f', -1, 64)
		if err := sender.Send("gauge", name, strValue); err != nil {
			log.Printf("Ошибка отправки gauge %s: %v", name, err)
		}
	}

	for name, value := range storage.SnapshotCounters() {
		strValue = strconv.FormatInt(int64(value), 10)
		if err := sender.Send("counter", name, strValue); err != nil {
			log.Printf("Ошибка отправки counter %s: %v", name, err)
		}
	}

}

func RunSender(storage *repository.MemStorage, options Options) {
	client := &http.Client{}

	collector := NewMetricCollector(storage)
	sender := NewSender("http://"+options.Address.String(), client)

	pollTicker := time.NewTicker(options.PollInterval.Duration)
	reportTicker := time.NewTicker(options.ReportInterval.Duration)
	defer pollTicker.Stop()
	defer reportTicker.Stop()

	for {

		select {

		case <-pollTicker.C:
			go collector.CollectMetrics()

		case <-reportTicker.C:
			go func() {
				sender.SendMetrics(storage)
				storage.SetCounter("PollCount", 0)
			}()
		}
	}

}
