package cmd

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/chzyer/readline"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type application struct {
	app        *kingpin.Application
	server     serverCommand
	sets       setsCommand
	setsView   setsViewCommand
	setsRm     setsRmCommand
	setsCreate setsCreateCommand
	setsList   setsListCommand
	commands   []runableCmd
	command    string
}

const version = "0.0.1"

var cli = application{}
var debug = false

func ParseCLI() error {
	cli.app = kingpin.New("pdbproxy", "Choria Discovery Proxy for PuppetDB")
	cli.app.Version(version)
	cli.app.Author("R.I. Pienaar <rip@devco.net>")
	cli.app.Flag("debug", "Enable debug logging").Short('d').BoolVar(&debug)

	cli.server.Setup()
	cli.sets.Setup()
	cli.setsList.Setup()
	cli.setsView.Setup()
	cli.setsRm.Setup()
	cli.setsCreate.Setup()

	cli.command = kingpin.MustParse(cli.app.Parse(os.Args[1:]))

	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.JSONFormatter{})

	return nil
}

func Run() error {
	var err error

	if debug {
		log.SetLevel(log.DebugLevel)
	}

	// these are all runableCmd figure out how to stick them in cli.commands
	// and do this by iteration
	switch cli.command {
	case cli.server.FullCommand():
		err = cli.server.Run()
	case cli.setsCreate.FullCommand():
		err = cli.setsCreate.Run()
	case cli.setsList.FullCommand():
		err = cli.setsList.Run()
	case cli.setsRm.FullCommand():
		err = cli.setsRm.Run()
	case cli.setsView.FullCommand():
		err = cli.setsView.Run()
	}

	return err
}

func askYN(prompt string) bool {
	rl, _ := readline.NewEx(&readline.Config{
		Prompt: prompt + " (y/n)> ",
	})

	for {
		ans, _ := rl.Readline()
		if ans != "" {
			if ans == "y" || ans == "Y" {
				return true
			}

			return false
		}
	}
}

func askQuery() string {
	fmt.Println("Please enter a PQL query, you can scroll back for history and use normal shell editing short cuts:")

	query := ""
	u, _ := user.Current()

	rl, err := readline.NewEx(&readline.Config{
		Prompt:      "pql> ",
		HistoryFile: filepath.Join(u.HomeDir, ".choria_history"),
	})
	if err != nil {
		return query
	}
	defer rl.Close()

	for {
		query, err := rl.Readline()
		if err != nil {
			return query
		}

		if query != "" {
			return query
		}
	}
}
