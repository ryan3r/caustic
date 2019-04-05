package main

import (
	"os"
	"io/ioutil"
	"log"
	"runner"
	"fmt"
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