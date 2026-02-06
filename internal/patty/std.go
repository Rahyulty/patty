package patty

import "os"

func Stdout() *os.File { return os.Stdout }
func Stderr() *os.File { return os.Stderr }
