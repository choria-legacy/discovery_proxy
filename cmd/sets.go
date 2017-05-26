package cmd

type setsCommand struct {
	command
}

func (s *setsCommand) Setup() error {
	s.Cmd = cli.app.Command("sets", "Node Set maintenance")

	return nil
}

func (s *setsCommand) Run() error {
	return nil
}
