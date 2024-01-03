package config

import (
	"flag"
	"os"
)

var Options struct {
	HostAddr     string
	DatabaseAddr string
	AccrualAddr  string
}

func InitConfig() {
	initFlags()
	initEnv()
}

func initFlags() {
	flag.StringVar(&Options.HostAddr, "a", ":8081", "Адрес и порт запуска сервиса")
	flag.StringVar(&Options.DatabaseAddr, "d", "", "Адрес подключения к базе данных")
	flag.StringVar(&Options.AccrualAddr, "r", ":8080", "Адрес системы расчёта начислений")
	flag.Parse()
}

func initEnv() {
	if envHostAddr := os.Getenv("RUN_ADDRESS"); envHostAddr != "" {
		Options.HostAddr = envHostAddr
	}
	if envDatabaseAddr := os.Getenv("DATABASE_URI"); envDatabaseAddr != "" {
		Options.DatabaseAddr = envDatabaseAddr
	}
	if envAccrualAddr := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAccrualAddr != "" {
		Options.AccrualAddr = envAccrualAddr
	}
}
