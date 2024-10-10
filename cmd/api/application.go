package main

import (
	"gitlab.com/sunflower4026/interview-flash-mobile/common/constants"
	"gitlab.com/sunflower4026/interview-flash-mobile/common/httpserver"
	"gitlab.com/sunflower4026/interview-flash-mobile/common/httpservice"
	"gitlab.com/sunflower4026/interview-flash-mobile/helper"
	"gitlab.com/sunflower4026/interview-flash-mobile/toolkit/db/migrations"
	"gitlab.com/sunflower4026/interview-flash-mobile/toolkit/db/postgres"
	"gitlab.com/sunflower4026/interview-flash-mobile/toolkit/log"
	"gitlab.com/sunflower4026/interview-flash-mobile/toolkit/runtimekit"
)

func main() {
	var err error

	helper.SetDefaultTimezone()

	appContext, cancel := runtimekit.NewRuntimeContext()
	defer func() {
		cancel()

		if err != nil {
			log.FromCtx(appContext).Error(err, "found error")
		}
	}()

	appConfig, err := helper.EnvConfigVariable(".env")
	if err != nil {
		return
	}

	// Setup Postgres
	postgresConf, err := postgres.NewFromConfig(appConfig, constants.DefaultPostgresPath)
	if err != nil {
		return
	}

	sqlDB, err := postgresConf.DB()
	if err != nil {
		return
	}

	err = migrations.Migrate(sqlDB)
	if err != nil {
		return
	}

	svc := httpservice.NewService(appConfig, postgresConf)

	httpserver.RunHTTPService(appContext, appConfig, svc)
}
