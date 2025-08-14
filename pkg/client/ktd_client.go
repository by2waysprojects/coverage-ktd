package client

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/by2waysprojects/coverage-ktd/model"
	"github.com/by2waysprojects/coverage-ktd/pkg/reporting"
)

type SecurityClient struct {
	wsClient      *WebSocketClient
	reportData    model.ReportData
	responseMutex sync.Mutex
}

func NewSecurityClient(url string, reportData model.ReportData) *SecurityClient {
	websocketClient, err := NewWebSocketClient(url)
	if err != nil {
		panic(err)
	}

	securityClient := &SecurityClient{
		wsClient:   websocketClient,
		reportData: reportData,
	}

	websocketClient.SetMessageHandler(securityClient.handleIncomingMessage)
	return securityClient
}

func (c *SecurityClient) handleIncomingMessage(message []byte) error {
	var response map[string]interface{}
	if err := json.Unmarshal(message, &response); err != nil {
		return err
	}

	fmt.Println("Received message:", response)

	c.responseMutex.Lock()
	defer c.responseMutex.Unlock()

	c.reportData.AttacksDetected = append(c.reportData.AttacksDetected, response)
	c.reportData.Detected++
	return nil
}

func (c *SecurityClient) GenerateReport() {
	c.reportData.EndTime = time.Now()
	c.reportData.DetectionRate = (float64(c.reportData.Detected) / float64(c.reportData.TotalTests)) * 100

	reporting.GenerateHTMLReport(c.reportData)
	reporting.GenerateJSONReport(c.reportData)

	fmt.Println("\n=== Security Test Summary ===")
	fmt.Printf("Tests executed: %d\n", c.reportData.TotalTests)
	fmt.Printf("Threats detected: %d (%.1f%%)\n", c.reportData.Detected, c.reportData.DetectionRate)
	fmt.Printf("Duration: %s\n", c.reportData.EndTime.Sub(c.reportData.StartTime))
	fmt.Println("Reports generated in 'reports' directory")
}
