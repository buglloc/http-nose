package httpfeature

import (
	"strings"
	"fmt"
	"errors"
	"math/rand"
	"strconv"
	"sort"
)

const concurrency = 30
var NotAlphaNumSyms []rune
var AlphaNumSyms []rune
var HttpMethods = []string{
	// RFC 2616
	"GET",	"HEAD",	"POST",	"PUT",	"DELETE", "OPTIONS", "TRACE", "CONNECT",
	// RFC 2518
	"PROPFIND", "PROPPATCH", "MKCOL", "COPY", "MOVE", "LOCK", "UNLOCK",
	// RFC 3744
	"ACL",
}


func init() {
	for i := 0; i < 0xFF; i++ {
		if i >= 0x30 && i <= 0x39 {
			// Nums
			AlphaNumSyms = append(AlphaNumSyms, rune(i))
			continue
		}

		if i >= 0x41 && i <= 0x5a {
			// Upper alpha
			AlphaNumSyms = append(AlphaNumSyms, rune(i))
			continue
		}

		if i >= 0x61 && i <= 0x7a {
			// Lower alpha
			AlphaNumSyms = append(AlphaNumSyms, rune(i))
			continue
		}

		NotAlphaNumSyms = append(NotAlphaNumSyms, rune(i))
	}
}

func TruncatingSprintf(str string, args ...interface{}) (string, error) {
	n := strings.Count(str, "%c")
	if n > len(args) {
		return "", errors.New("Unexpected string:" + str)
	}
	return fmt.Sprintf(str, args[:n]...), nil
}

func RandAlphanumString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = AlphaNumSyms[rand.Int63() % int64(len(AlphaNumSyms))]
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

func PrintableRunes(data []rune) string {
	if len(data) == 0 {
		return "none"
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i] < data[j]
	})

	result := make([]string, 0)
	max := len(data) - 1
	var first rune = -1
	continuousRune := false
	for i, v := range data {
		if first == -1 {
			first = v
		}

		if i != max && data[i+1]-v == 1 {
			continuousRune = true
			continue
		}

		if first == v {
			result = append(result, fmt.Sprintf("\\x%02X", first))
			first = -1
			continuousRune = false
		} else {
			sep := ","
			if continuousRune {
				sep = "-"
			}
			result = append(result, fmt.Sprintf("\\x%02X%s\\x%02X", first, sep, v))
			first = -1
			continuousRune = false
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
