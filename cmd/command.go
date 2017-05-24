package cmd

import (
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type command struct {
	Run   func() error
	Setup func() error

	Cmd *kingpin.CmdClause
}

type runableCmd interface {
	Setup() error
	Run() error
	FullCommand() string
}

func (c *command) FullCommand() string {
	return c.Cmd.FullCommand()
}
