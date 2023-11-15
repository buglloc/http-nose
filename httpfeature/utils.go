package httpfeature

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"strings"
)

const concurrency = 16

var NotAlphaNumSyms []byte
var AlphaNumSyms []byte
var HttpMethods = []string{
	// RFC 2616
	"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS", "TRACE", "CONNECT",
	// RFC 2518
	"PROPFIND", "PROPPATCH", "MKCOL", "COPY", "MOVE", "LOCK", "UNLOCK",
	// RFC 3744
	"ACL",
}

func init() {
	for i := 0; i < 0xFF; i++ {
		if i >= 0x30 && i <= 0x39 {
			// Nums
			AlphaNumSyms = append(AlphaNumSyms, byte(i))
			continue
		}

		if i >= 0x41 && i <= 0x5a {
			// Upper alpha
			AlphaNumSyms = append(AlphaNumSyms, byte(i))
			continue
		}

		if i >= 0x61 && i <= 0x7a {
			// Lower alpha
			AlphaNumSyms = append(AlphaNumSyms, byte(i))
			continue
		}

		NotAlphaNumSyms = append(NotAlphaNumSyms, byte(i))
	}
}

func Symf(str string, sym byte) string {
	return strings.ReplaceAll(str, "{sym}", string([]byte{sym}))
}

func RandAlphanumString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = AlphaNumSyms[rand.Int63()%int64(len(AlphaNumSyms))]
	}
	return string(b)
}

func PrintableStrings(data []string) string {
	if len(data) == 0 {
		return "none"
	}

	result := make([]string, len(data))
	for i, v := range data {
		result[i] = strconv.Quote(v)
	}
	return strings.Join(result, ", ")
}

func PrintableSymbols(data []byte) string {
	if len(data) == 0 {
		return "none"
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i] < data[j]
	})

	var result []string
	max := len(data) - 1
	first := -1
	continuousByte := false
	for i, v := range data {
		if first == -1 {
			first = int(v)
		}

		if i != max && data[i+1]-v == 1 {
			continuousByte = true
			continue
		}

		if first == int(v) {
			result = append(result, fmt.Sprintf("\\x%02X", first))
			first = -1
			continuousByte = false
		} else {
			sep := ","
			if continuousByte {
				sep = "-"
			}
			result = append(result, fmt.Sprintf("\\x%02X%s\\x%02X", first, sep, v))
			first = -1
			continuousByte = false
		}
	}
	if first != -1 {
		result = append(result, fmt.Sprintf("\\x%02X", first))
	}
	return strings.Join(result, ",")
}

func PrintableBool(data bool) string {
	if data {
		return "Yes"
	}
	return "No"
}

func PrintableBytes(bytes int) string {
	const (
		BYTE     = 1.0
		KILOBYTE = 1024 * BYTE
		MEGABYTE = 1024 * KILOBYTE
		GIGABYTE = 1024 * MEGABYTE
		TERABYTE = 1024 * GIGABYTE
	)

	unit := ""
	value := float32(bytes)

	switch {
	case bytes >= TERABYTE:
		unit = "T"
		value = value / TERABYTE
	case bytes >= GIGABYTE:
		unit = "G"
		value = value / GIGABYTE
	case bytes >= MEGABYTE:
		unit = "M"
		value = value / MEGABYTE
	case bytes >= KILOBYTE:
		unit = "K"
		value = value / KILOBYTE
	case bytes >= BYTE:
		unit = "B"
	case bytes == 0:
		return "0"
	}

	stringValue := fmt.Sprintf("%.1f", value)
	return fmt.Sprintf("%s%siB", stringValue, unit)
}

func GoldenSectionSearch(beforeX, afterX, prec float64, testFunc func(x float64) float64) float64 {
	spacing := (3.0 - math.Sqrt(5)) / 2.0

	var x [4]float64
	var y [4]float64

	x[0], x[3] = beforeX, afterX
	gapSize := afterX - beforeX
	x[1], x[2] = x[0]+spacing*gapSize, x[3]-spacing*gapSize

	for i := 0; i < 4; i++ {
		y[i] = testFunc(x[i])
	}

	stepCount := int(math.Ceil(math.Log(prec/gapSize) / math.Log(1-spacing)))
	for i := 0; i < stepCount; i++ {
		// If we have narrowed down the region to machine
		// precision, we cannot go any further.
		if x[1] <= x[0] || x[2] <= x[1] || x[3] <= x[2] {
			break
		}

		if y[1] < y[2] {
			x[3] = x[2]
			diff := x[3] - x[0]
			x[1] = x[0] + diff*spacing
			x[2] = x[3] - diff*spacing
			y[2], y[3] = y[1], y[2]
			y[1] = testFunc(x[1])
		} else {
			x[0] = x[1]
			diff := x[3] - x[0]
			x[1] = x[0] + diff*spacing
			x[2] = x[3] - diff*spacing
			y[0], y[1] = y[1], y[2]
			y[2] = testFunc(x[2])
		}
	}

	return (x[3] + x[0]) / 2
}
