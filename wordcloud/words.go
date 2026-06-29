package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
)

// sizedWord is a word with its resolved point size, ready for layout.
type sizedWord struct {
	text  string
	size  float64
	color string // per-word color override; empty means "use palette"
}

// parseFrequency reads "count word" lines (the output of a pipeline such as
//
//	cat *.txt | tr ' ' '\012' | sort | uniq -c | sort -n
//
// ) from r and returns them as count-bearing WordSpecs. A bare "word" line with
// no leading count is treated as count 1. Blank lines are ignored.
func parseFrequency(r io.Reader) ([]WordSpec, error) {
	var specs []WordSpec
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		fields := strings.Fields(sc.Text())
		if len(fields) == 0 {
			continue
		}
		if len(fields) >= 2 {
			if n, err := strconv.Atoi(fields[0]); err == nil {
				specs = append(specs, WordSpec{Text: fields[1], Count: n})
				continue
			}
		}
		// No leading count: treat the first token as a word with count 1.
		specs = append(specs, WordSpec{Text: fields[0], Count: 1})
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("read frequency input: %w", err)
	}
	return specs, nil
}

// buildWords resolves each spec to a concrete point size. A spec with Size > 0
// uses that size verbatim. Remaining specs are sized from their Count using
// sqrt scaling between minSize and maxSize, which compresses large counts so a
// single dominant word does not dwarf the rest. Specs with neither size nor a
// positive count fall back to minSize. Words with empty text are dropped.
func buildWords(specs []WordSpec, minSize, maxSize float64) []sizedWord {
	minCount, maxCount := math.MaxInt, math.MinInt
	for _, s := range specs {
		if s.Size > 0 || s.Text == "" {
			continue
		}
		c := s.Count
		if c < minCount {
			minCount = c
		}
		if c > maxCount {
			maxCount = c
		}
	}

	scale := func(count int) float64 {
		if count <= 0 || maxCount <= minCount {
			// Degenerate range (single distinct count, or no counts): use maxSize
			// so frequency-only clouds are not all rendered at the minimum.
			return maxSize
		}
		frac := float64(count-minCount) / float64(maxCount-minCount)
		return minSize + (maxSize-minSize)*math.Sqrt(frac)
	}

	var words []sizedWord
	for _, s := range specs {
		if s.Text == "" {
			continue
		}
		size := s.Size
		if size <= 0 {
			size = scale(s.Count)
		}
		words = append(words, sizedWord{text: s.Text, size: size, color: s.Color})
	}
	return words
}
