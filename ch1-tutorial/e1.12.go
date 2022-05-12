// Modify the Lissajous server to read parameter values from the URL.
// For example, you might arrange it so that a URL like
// http://localhost:8000/?cycles=20 sets the number of cycles to 20
// instead of the default 5. Use the strconv.Atoi function to convert
// the string parameter into an integer. You can see its documentation
// with go doc stringconv.Atoi.

package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

var palette = []color.Color{color.White, color.Black}

const (
	whiteIndex = 0 // first color in palette
	blackIndex = 1 // next color in palette
)

func main() {
	// The sequence of images is deterministic unless we seed
	// the pseudo-random number generator using the current time.
	// Thanks to Randall McPherson for pointing out the omission.
	rand.Seed(time.Now().UTC().UnixNano())

	if len(os.Args) > 1 && os.Args[1] == "web" {
		handler := func(w http.ResponseWriter, r *http.Request) {
			lissajous(w, r)
		}
		http.HandleFunc("/", handler)
		log.Fatal(http.ListenAndServe("localhost:8000", nil))
		return
	}
	lissajous(os.Stdout, nil)
}

func lissajous(out io.Writer, req *http.Request) {
	const (
		defaultCycles = 5
		res           = 0.001 // angular resolution
		size          = 100   // image canvas covers [-size..+size]
		nframes       = 64    // number of animation frames
		delay         = 8     // delay between frames in 10ms units
	)

	cycles := defaultCycles

	if req != nil {
		for query, queryArgs := range req.URL.Query() {
			if query == "cycles" && len(queryArgs) > 0 {
				var err error
				cycles, err = strconv.Atoi(queryArgs[0])

				if err != nil || cycles < 1 {
					cycles = defaultCycles
				}
			}
		}
	}

	// We print to stderr to avoid clobbering stdout, if this is not
	// run in web mode
	if cycles == defaultCycles {
		fmt.Fprintln(os.Stderr, "Using default cycles.")
	} else {
		fmt.Fprintf(os.Stderr, "Using non-default cycles of %d.\n", cycles)
	}

	freq := rand.Float64() * 3.0 // relative frequency of y oscillator
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0 // phase difference

	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)

		for t := 0.0; t < float64(cycles)*2.0*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5),
				blackIndex)
		}

		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}

	gif.EncodeAll(out, &anim) // NOTE: ignoring encoding errors
}
