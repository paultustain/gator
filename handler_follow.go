package main

import "errors"

func handlerFollow(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("not enough arguements")
	}

	return nil
}
