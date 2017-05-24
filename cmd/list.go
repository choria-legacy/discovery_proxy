package cmd

import (
	"fmt"

	"github.com/choria-io/pdbproxy/choria"
)

type setsListCommand struct {
	command
}

func (s *setsListCommand) Setup() error {
	s.Cmd = cli.sets.Cmd.Command("list", "List known sets").Default().Alias("ls")

	return nil
}

func (s *setsListCommand) Run() error {
	c, err := choria.New(choria.UserConfig())
	if err != nil {
		return err
	}

	sets, err := choria.NewSets(c)
	if err != nil {
		return err
	}

	return sets.List(func(results []string) error {
		fmt.Printf("Found %d set(s)\n\n", len(results))

		choria.SliceGroups(results, 3, func(group []string) {
			for _, set := range group {
				fmt.Printf("   %-20s", set)
			}

			fmt.Println("")
		})

		fmt.Println("")

		return nil
	})
}
