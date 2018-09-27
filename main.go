package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/globalsign/mgo"
	"github.com/vitalyo61/genkeys/config"
	"github.com/vitalyo61/genkeys/db"
	"github.com/vitalyo61/genkeys/generator"
	"github.com/vitalyo61/genkeys/web"
)

func main() {
	conf := config.Make()
	err := conf.Open("genkeys.conf")
	if err != nil {
		panic(err)
	}

	db, err := db.Make(conf.DataBase.URL)
	if err != nil {
		panic(err)
	}

	chData := make(chan *generator.Data)

	server := web.Make(conf.Server, db, chData)

	cLast, err := db.CodeLast()
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}

	gen, err := generator.Make(cLast.Number)
	if err != nil {
		panic(err)
	}

	gen.Start(chData)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		for {
			switch s := <-sigs; {
			case s == syscall.SIGTERM || s == syscall.SIGINT:
				server.Stop()

				ch := make(chan *generator.Result)
				chData <- &generator.Data{
					Cmd:      generator.CmdStop,
					ChanCode: ch,
				}
				<-ch

				os.Exit(0)
			}
		}
	}()

	server.Start()
}
