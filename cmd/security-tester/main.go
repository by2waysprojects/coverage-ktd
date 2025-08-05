package main

import (
	"os"
	"time"

	"github.com/by2waysprojects/coverage-ktd/model"
	"github.com/by2waysprojects/coverage-ktd/pkg/attacks"
	"github.com/by2waysprojects/coverage-ktd/pkg/client"
)

const (
	attackConfigDir = "attacks"
)

func main() {
	securityAppUri := os.Getenv("TARGET_SECURITY_APP")
	vulnerableAppUri := os.Getenv("TARGET_VULNERABLE_APP")

	if securityAppUri == "" || vulnerableAppUri == "" {
		panic("TARGET_SECURITY_APP and TARGET_VULNERABLE_APP environment variables must be set")
	}

	attackExecutor := attacks.NewAttackExecutor(vulnerableAppUri)
	attackExecutor.LoadAttacks(attackConfigDir)
	attacks := attackExecutor.GetAttacks()

	reportData := model.ReportData{
		StartTime:  time.Now(),
		Tests:      attacks,
		TotalTests: len(attacks),
	}
	securityClient := client.NewSecurityClient(securityAppUri, reportData)
	attackExecutor.RunAll()

	timer := time.NewTimer(10 * time.Minute)
	<-timer.C

	securityClient.GenerateReport()
}
