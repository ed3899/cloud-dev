package main

import (
	"log"

	"github.com/ed3899/kumo/utils"
)

func init() {
	_, err := utils.GetBinaries()
	if err != nil {
		log.Fatalf("Error occurred while getting binaries: %v", err)
	}
}

func main() {
	

}
