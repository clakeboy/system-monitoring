package command

import (
	"fmt"
	"os"
)

func StartDaemon() {
	fmt.Println(os.Getppid())
	os.Exit(0)
}
