package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runner"
)

func main() {
	os.Chdir("/mnt/submissions")

	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fmt.Printf("%v: %v\n", file.Name(), runner.Test(file.Name(), "1 2 3"))
	}
}
