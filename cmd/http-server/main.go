package main

import (
	"http-proxy/internal"
	"log"
)

func main() {
	log.Print(httpProxySplash)

	conf, err := internal.LoadConfig("/etc/http-proxy/config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	log.Print(conf)

	hpd, err := internal.NewHttpProxyD(conf)
	if err != nil {
		log.Fatal(err)
	}

	if err := hpd.Run(); err != nil {
		log.Fatal(err)
	}
}

const httpProxySplash = `

 /$$   /$$ /$$$$$$$$ /$$$$$$$$ /$$$$$$$        /$$$$$$$                                        
| $$  | $$|__  $$__/|__  $$__/| $$__  $$      | $$__  $$                                       
| $$  | $$   | $$      | $$   | $$  \ $$      | $$  \ $$ /$$$$$$   /$$$$$$  /$$   /$$ /$$   /$$
| $$$$$$$$   | $$      | $$   | $$$$$$$/      | $$$$$$$//$$__  $$ /$$__  $$|  $$ /$$/| $$  | $$
| $$__  $$   | $$      | $$   | $$____/       | $$____/| $$  \__/| $$  \ $$ \  $$$$/ | $$  | $$
| $$  | $$   | $$      | $$   | $$            | $$     | $$      | $$  | $$  >$$  $$ | $$  | $$
| $$  | $$   | $$      | $$   | $$            | $$     | $$      |  $$$$$$/ /$$/\  $$|  $$$$$$$
|__/  |__/   |__/      |__/   |__/            |__/     |__/       \______/ |__/  \__/ \____  $$
                                                                                      /$$  | $$
                                                                                     |  $$$$$$/
                                                                                      \______/ 

`
