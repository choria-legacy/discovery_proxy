package cmd

import (
	"fmt"
	"os"

	"github.com/choria-io/discovery_proxy/choria"
	"github.com/choria-io/discovery_proxy/client"
)

type setsCreateCommand struct {
	command

	set   *string
	force *bool
	query *string
}

func (s *setsCreateCommand) Setup() error {
	s.Cmd = cli.sets.Cmd.Command("create", "Create or Update a Set")
	s.query = s.Cmd.Flag("query", "The query to use for the set").Default("").String()
	s.force = s.Cmd.Flag("update", "Updates an existing set").Short('f').Default("false").Bool()
	s.set = s.Cmd.Arg("set", "Set to create or udpate").Required().String()

	return nil
}

func (s *setsCreateCommand) Run() error {
	c, err := choria.New(choria.UserConfig())
	if err != nil {
		return err
	}

	sets, err := client.NewSets(c)
	if err != nil {
		return err
	}

	if sets.HaveSet(s.set) && !*s.force {
		fmt.Printf("A set called '%s' already exist, use --update to update it or pick a new name", *s.set)
		os.Exit(1)
	}

	query := ""

	for {
		if query == "" {
			query = askQuery()
			if query == "" {
				break
			}
		}

		if err := s.resolveAndPrintPQL(query, sets); err != nil {
			fmt.Printf("Could not perform query: %s\n\n", err.Error())
			query = ""
			continue
		}

		fmt.Println()

		if askYN("Do you want to store this query") {
			if *s.force {
				if err := sets.Update(*s.set, query); err != nil {
					return err
				}
			} else {
				if err := sets.Create(*s.set, query); err != nil {
					return err
				}
			}

			break
		} else {
			query = ""
			continue
		}
	}

	return nil
}

func (s *setsCreateCommand) resolveAndPrintPQL(pql string, sets *client.Sets) error {
	nodes, err := sets.ResolvePQL(pql)
	if err != nil {
		return err
	}

	fmt.Print("Matched Nodes:\n\n")

	sets.PrintNodes(nodes)

	return nil
}
