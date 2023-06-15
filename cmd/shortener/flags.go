package main

import (
	"flag"

	"github.com/MWT-proger/shortener/configs"
)

// parseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func parseFlags(conf *configs.Config) {

	flag.StringVar(&conf.HostServer, "a", conf.HostServer, "address and port to run server")
	flag.StringVar(&conf.BaseURLShortener, "b", conf.BaseURLShortener, "base URL for a short link")
	flag.Parse()
}
