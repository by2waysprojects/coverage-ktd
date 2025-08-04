package model

import "time"

type ReportData struct {
	StartTime       time.Time
	EndTime         time.Time
	TotalTests      int
	Detected        int
	Tests           map[string]Attack
	AttacksDetected []interface{}
	DetectionRate   float64
}
