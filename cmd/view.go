package cmd

import (
	"fmt"

	"github.com/choria-io/pdbproxy/choria"
)

type setsViewCommand struct {
	command

	set      *string
	discover *bool
}

func (s *setsViewCommand) Setup() error {
	s.Cmd = cli.sets.Cmd.Command("view", "View a set")
	s.discover = s.Cmd.Flag("discover", "Also show nodes matched by the set").Default("false").Bool()
	s.set = s.Cmd.Arg("set", "Set to view").Required().String()

	return nil
}

func (s *setsViewCommand) Run() error {
	c, err := choria.New(choria.UserConfig())
	if err != nil {
		return err
	}

	sets, err := choria.NewSets(c)
	if err != nil {
		return err
	}

	if err := sets.PrintSet(s.set, *s.discover); err != nil {
		return err
	}

	if !*s.discover {
		fmt.Println("Use --discover to view matching nodes\n")
	}

	return nil
}
