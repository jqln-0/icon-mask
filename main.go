package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
)

func Usage() {
	fmt.Printf("usage: %v [flags] theme baseimg outdir\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	// Register flags.
	flag.Usage = Usage
	filtThemesRaw := flag.String("extra-themes", "", "Comma seperated list of additional themes to filter.")
	categoriesRaw := flag.String("categories", "apps,pixmaps", "Comma seperated list of categories to include.")
	maxProcs := flag.Int("maxprocs", runtime.NumCPU(), "Number of threads to use.")

	// Parse flags.
	GTKInit(os.Args)
	flag.Parse()

	runtime.GOMAXPROCS(*maxProcs)

	var filtThemes []GTKTheme
	for _, theme := range strings.Split(*filtThemesRaw, ",") {
		if theme != "" {
			filtThemes = append(filtThemes, CreateTheme(theme))
		}
	}

	var categories []string
	for _, cat := range strings.Split(*categoriesRaw, ",") {
		categories = append(categories, cat)
	}

	// Check input args.
	if flag.NArg() != 3 {
		flag.Usage()
		return
	}

	// Load the source theme.
	mainTheme := CreateTheme(flag.Arg(0))

	// Prepare the drawer.
	drawer := CreateMaskDrawer(flag.Arg(1), mainTheme)

	// Generate a chan containing the icons we need to generate.
	source := GenerateIconNames(mainTheme)
	filtered := CreateThemeFilter(mainTheme, source)
	for _, theme := range filtThemes {
		filtered = CreateThemeFilter(theme, filtered)
	}
	properties := CreatePropertyFiller(mainTheme, filtered)
	out := CreateCategoryFilter(categories, properties)

	// Generate and save an image for each icon.
	for icon := range out {
		drawer.CreateIcons(icon, flag.Arg(2))
	}
}
