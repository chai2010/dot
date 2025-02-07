// Copyright © 2025 ChaiShushan <chaishushan{AT}gmail.com>.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/goccy/go-graphviz"
)

var (
	flagFormatDot  = flag.Bool("Tdot", false, "set dot format")
	flagFormatSvg  = flag.Bool("Tsvg", false, "set svg format")
	flagFormatPng  = flag.Bool("Tpng", false, "set png format")
	flagFormatJpg  = flag.Bool("Tjpg", false, "set jpg format")
	flagOutputFile = flag.String("o", "", "set output file")
)

func main() {
	flag.Parse()

	format := graphviz.SVG
	switch {
	case *flagFormatDot:
		format = graphviz.XDOT
	case *flagFormatSvg:
		format = graphviz.SVG
	case *flagFormatPng:
		format = graphviz.PNG
	case *flagFormatJpg:
		format = graphviz.JPG
	}

	if err := run(format, *flagOutputFile); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(format graphviz.Format, outputFile string) (e error) {
	ctx := context.Background()

	graph, err := readGraph()
	if err != nil {
		return err
	}
	g, err := graphviz.New(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err := graph.Close(); err != nil {
			e = err
		}
		if err := g.Close(); err != nil {
			e = err
		}
	}()

	if outputFile != "" {
		return g.RenderFilename(ctx, graph, format, outputFile)
	}

	return g.Render(ctx, graph, format, os.Stdout)
}

func readGraph() (*graphviz.Graph, error) {
	bytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, err
	}
	return graphviz.ParseBytes(bytes)
}
