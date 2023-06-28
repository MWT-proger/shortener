package main

import (
	"flag"

	"github.com/MWT-proger/shortener/configs"
)

// parseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func parseFlags(conf *configs.Config) {

	flag.StringVar(&conf.HostServer, "a", conf.HostServer, "Адрес и порт для запуска сервера")
	flag.StringVar(&conf.BaseURLShortener, "b", conf.BaseURLShortener, "Базовый URl  который будет использоваться для короткой ссылки")
	flag.StringVar(&conf.LogLevel, "l", "info", "Уровень логирования")
	flag.Parse()
}
