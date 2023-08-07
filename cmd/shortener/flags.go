package main

import (
	"flag"

	"github.com/MWT-proger/shortener/configs"
)

// parseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func parseFlags(conf *configs.Config) {

	flag.StringVar(&conf.HostServer, "a", conf.HostServer, "адрес и порт для запуска сервера")
	flag.StringVar(&conf.DatabaseDSN, "d", conf.DatabaseDSN, "строка с адресом подключения к БД")
	flag.StringVar(&conf.JSONFileDB, "f", conf.JSONFileDB, "полное имя файла, куда сохраняются данные в формате JSON")
	flag.StringVar(&conf.BaseURLShortener, "b", conf.BaseURLShortener, "базовый URl  который будет использоваться для короткой ссылки")
	flag.StringVar(&conf.LogLevel, "l", "info", "уровень логирования")
	flag.Parse()
}
