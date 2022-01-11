package main

import (
	"errors"
	"log"
	mosaic "mosaic/pkg"
	"os"
	"strconv"
	"strings"
)

func main() {
	args := os.Args[1:]
	if len(args) != 4 {
		errors.New("Usage: <source_image> <library_images> <gridsize(cols,rows)> <output_directory")
	}
	gridSize := strings.Split(args[2], ",")
	gR, gC := gridSize[0], gridSize[1]
	rows, err := strconv.Atoi(gR)
	if err != nil {
		log.Fatal(err)
	}
	cols, err := strconv.Atoi(gC)
	if err != nil {
		log.Fatal(err)
	}
	mosaic.CreateMosaic(
		args[0],
		args[1],
		[]int{rows, cols},
		args[3],
	)
}
