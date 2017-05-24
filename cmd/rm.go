package cmd

import (
	"fmt"
	"os"

	"github.com/choria-io/pdbproxy/choria"
)

type setsRmCommand struct {
	command

	set   *string
	force *bool
}

func (s *setsRmCommand) Setup() error {
	s.Cmd = cli.sets.Cmd.Command("rm", "Delete a Set")
	s.force = s.Cmd.Flag("force", "Force delete without prompting").Short('f').Default("false").Bool()
	s.set = s.Cmd.Arg("set", "Set to delete").Required().String()

	return nil
}

func (s *setsRmCommand) Run() error {
	c, err := choria.New(choria.UserConfig())
	if err != nil {
		return err
	}

	sets, err := choria.NewSets(c)
	if err != nil {
		return err
	}

	if !sets.HaveSet(s.set) {
		fmt.Printf("Could not find set %s", *s.set)
		os.Exit(1)
	}

	if !*s.force {
		if err := sets.PrintSet(s.set, false); err != nil {
			return err
		}

		if !askYN("Are you sure you wish to delete this node set?") {
			fmt.Printf("Node set %s was not deleted.", *s.set)
			os.Exit(0)
		}
	}

	return sets.Rm(s.set)
}
