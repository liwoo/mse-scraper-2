package utils

import (
	"fmt"
	"os"
)

func EnsureDirsExist(dirs []string) {
	for _, dir := range dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0755)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}
}
