package util

import (
	"encoding/csv"
	"os"
)

type Converter interface {
	ToCSVRow() []string
}

type Generator struct {
	f *os.File
	w *csv.Writer
}

func New(filename string, options ...Option) Generator {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	w := csv.NewWriter(f)

	generator := Generator{
		f: f,
		w: w,
	}

	for _, option := range options {
		option(&generator)
	}

	return generator
}

type Option func(*Generator)

func WithHeader(header []string) Option {
	return func(g *Generator) {
		g.w.Write(header)
	}
}

func (g Generator) Write(data Converter) error {
	return g.w.Write(data.ToCSVRow())
}

func (g Generator) Done() {
	g.w.Flush()
	g.f.Close()
}
