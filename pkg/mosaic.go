package mosaic

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
)

type RGB struct {
	R, G, B uint8
}

type SourceImage struct {
	Path     string
	img      image.Image
	averages []RGB
	tiles    []image.Rectangle
}

type LibraryImage struct {
	Path    string
	img     image.Image
	average RGB
}

type Match struct {
	SourceIndex  int
	LibraryIndex int
}

func (si *SourceImage) load() error {
	sourceFile, err := os.Open(si.Path)
	if err != nil {
		log.Fatalf("While loading source image, os.Open() failed with %s\n", err)
	}
	defer sourceFile.Close()
	img, _, err := image.Decode(sourceFile)
	if err != nil {
		log.Fatalf("While loading source image, image.Decode() failed with %s\n", err)
	}
	si.img = img
	return err
}

func splitSourceToTiles(si *SourceImage, gridSize []int) {
	if gridSize == nil {
		gridSize = []int{32, 32}
	}
	si.tiles = []image.Rectangle{}
	width, height := si.img.Bounds().Size().X, si.img.Bounds().Size().Y
	fmt.Println(width, height)
	m, n := gridSize[0], gridSize[1]
	w, h := int(width/m), int(height/n)

	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			si.tiles = append(si.tiles, image.Rectangle{
				image.Point{i * w, j * h},
				image.Point{(i + 1) * w, (j + 1) * h},
			})
			fmt.Printf("[%v:%v] (%v,%v)-(%v,%v)\n", i, j, i*w, j*h, (i+1)*w, (j+1)*h)
		}
	}
}

func calculateSourceAverages(si *SourceImage) {
	si.averages = make([]RGB, len(si.tiles))
	fmt.Println("Number of tiles: ", len(si.tiles))
	for tileIndex, tile := range si.tiles {
		tileImage := si.img.(*image.RGBA).SubImage(tile)
		// fmt.Printf("Tile %v: %v\n", tileIndex, tile)
		var (
			cumR uint64
			cumG uint64
			cumB uint64
		)
		totalPixels := uint64(0)
		for x := tile.Min.X; x < tile.Max.X; x++ {
			for y := tile.Min.Y; y < tile.Max.Y; y++ {
				totalPixels++
				r, g, b, _ := tileImage.At(x, y).RGBA()
				cumR += uint64(r)
				cumG += uint64(g)
				cumB += uint64(b)
				// fmt.Printf("[%v] (%v, %v) R:%v G:%v B:%v\n", tileIndex, x, y, uint8(cumR), uint8(cumG), uint8(cumB))
			}
		}
		si.averages[tileIndex] = RGB{
			uint8((cumR / totalPixels) >> 8),
			uint8((cumG / totalPixels) >> 8),
			uint8((cumB / totalPixels) >> 8),
		}
		fmt.Println(si.averages[tileIndex])
	}
}

func readLibraryImages(libraryPath string) ([]LibraryImage, error) {
	files := []string{}
	acceptedFileTypes := []string{"*.png", "*.jpg", "*.jpeg"}
	var err error
	for _, accepted := range acceptedFileTypes {
		found, err := filepath.Glob(libraryPath + "/" + accepted)
		if err != nil {
			log.Fatalf("While loading library images, filepath.Glob() failed with %s\n", err)
		}
		files = append(files, found...)
	}
	fmt.Println("Found library images: ", len(files))

	libraryImages := []LibraryImage{}
	for _, filePath := range files {
		libraryFile, err := os.Open(filePath)
		if err != nil {
			log.Fatalf("While loading a library image, os.Open() failed with %s\n", err)
		}
		defer libraryFile.Close()
		img, _, err := image.Decode(libraryFile)
		if err != nil {
			log.Fatalf("While loading a library image, image.Decode() failed with %s\n", err)
		}

		// Calculate average color of the library image
		width, height := img.Bounds().Size().X, img.Bounds().Size().Y
		var (
			cumR uint64
			cumG uint64
			cumB uint64
		)
		totalPixels := uint64(0)
		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				totalPixels++
				r, g, b, _ := img.At(x, y).RGBA()
				cumR += uint64(r)
				cumG += uint64(g)
				cumB += uint64(b)
			}
		}
		libraryImages = append(libraryImages,
			LibraryImage{
				Path: filePath,
				img:  img,
				average: RGB{
					uint8((cumR / totalPixels) >> 8),
					uint8((cumG / totalPixels) >> 8),
					uint8((cumB / totalPixels) >> 8),
				},
			},
		)
	}

	return libraryImages, err
}

func findBestMatches(sourceImage *SourceImage, libraryImages []LibraryImage) []Match {
	matches := []Match{}

	// For each grid in the source image ...
	for sIndex, sAverage := range sourceImage.averages {
		// Calculate the distance of the averages from the library ...
		var minimumIndex int
		var minimumDistance float64 = math.Inf(0)
		for libraryIndex, libraryImage := range libraryImages {
			// Then choose the closest one
			tAverage := libraryImage.average
			distance := (math.Pow((float64(tAverage.R)-float64(sAverage.R)), 2.0) +
				math.Pow((float64(tAverage.G)-float64(sAverage.G)), 2.0) +
				math.Pow((float64(tAverage.B)-float64(sAverage.B)), 2.0))
			if distance < minimumDistance {
				minimumDistance = distance
				minimumIndex = libraryIndex
			}
		}
		fmt.Printf("Source grid: [%02d] %3v \tClosest img: [%02d] %3v\n", sIndex, sAverage, minimumIndex, libraryImages[minimumIndex].average)
		matches = append(matches, Match{
			SourceIndex:  sIndex,
			LibraryIndex: minimumIndex,
		})
	}

	// TODO return error if any
	return matches
}

func generateMosaic(sourceImage *SourceImage, libraryImages []LibraryImage, matches []Match, outputDir string) (string, error) {
	outputDir = filepath.Dir(outputDir)
	outputExtension := filepath.Ext(sourceImage.Path)
	inputFileName := strings.Split(filepath.Base(sourceImage.Path), ".")[0]
	outputPath := filepath.Join(outputDir, inputFileName+"_out"+outputExtension)

	img := image.NewRGBA(sourceImage.img.Bounds())
	for _, match := range matches {
		dp := sourceImage.tiles[match.SourceIndex].Min
		libImage := libraryImages[match.LibraryIndex].img
		// Resizing to grid size
		// Resize the cropped image to width = 200px preserving the aspect ratio.
		w, h := sourceImage.tiles[match.SourceIndex].Max.X, sourceImage.tiles[match.SourceIndex].Max.Y
		resizedLibImage := imaging.Fill(libImage, w, h, imaging.Center, imaging.Lanczos)
		sr := libImage.Bounds()
		r := image.Rectangle{dp, dp.Add(sourceImage.tiles[match.SourceIndex].Size())}

		// fmt.Println(
		// 	">>",
		// 	r,
		// 	sr.Min,
		// )
		// To copy from a rectangle sr in the source image
		// to a rectangle starting at a point dp in the destination,
		// convert the source rectangle into the destination image’s co-ordinate space:
		draw.Draw(
			img,             // dst
			r,               // r
			resizedLibImage, // src
			sr.Min,          // sp
			draw.Src,        // op
		)
	}

	// Encode as PNG.
	f, err := os.Create(outputPath)
	if err != nil {
		log.Fatal(err)
	}

	switch outputExtension {
	case ".png":
		err = png.Encode(f, img)
		if err != nil {
			log.Fatal(err)
		}
	case ".jpg":
	case ".jpeg":
		opt := jpeg.Options{
			Quality: 90,
		}
		err = jpeg.Encode(f, img, &opt)
		if err != nil {
			log.Fatal(err)
		}
	}

	return outputPath, err
}

func CreateMosaic(sourcePath string, libraryPath string, gridSize []int, outputDir string) string {
	// Reading source image
	sourceImage := SourceImage{
		Path:     sourcePath,
		img:      nil,
		averages: []RGB{},
		tiles:    []image.Rectangle{},
	}
	err := sourceImage.load()
	splitSourceToTiles(&sourceImage, gridSize)
	calculateSourceAverages(&sourceImage)
	if err != nil {
		log.Fatalf("Got an error loading the source image %s\n", err)
	}
	// Reading library images ("thumbnails")
	libraryImages, err := readLibraryImages(libraryPath)
	if err != nil {
		log.Fatalf("Got an error loading the library images %s\n", err)
	}
	fmt.Println(libraryImages[len(libraryImages)-1])
	// Find best matches for each grid of the source image
	// and the library images
	matches := findBestMatches(
		&sourceImage,
		libraryImages,
	)
	// for _, v := range matches {
	// 	fmt.Println("Match", v)
	// }

	mosaic, err := generateMosaic(
		&sourceImage,
		libraryImages,
		matches,
		outputDir,
	)
	if err != nil {
		log.Fatalf("Got an error generating the mosaic %s\n", err)
	}

	fmt.Print("Output written to : ", mosaic)

	// TODO: Later on, return (path to new image AND an error)
	return mosaic
}
