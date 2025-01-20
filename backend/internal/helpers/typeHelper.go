package helpers

import (
	"strconv"
	"strings"
	"time"
)

func ParseDateTime(datetime string) (time.Time, error) {
	parsedDateTime, err := time.Parse("2006-01-02T15:04:05.999", datetime)
	if err != nil {
		return time.Time{}, err
	}

	return parsedDateTime, nil
}

func normalizeEurope(old string) string {
	count := strings.Count(old, ".")
	s := strings.Replace(old, ",", ".", -1)
	return strings.Replace(s, ".", "", count)

}
func normalizeAmericanBritain(old string) string {
	return strings.Replace(old, ",", "", -1)
}

func reverse(a []byte) string {
	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}
	return string(a)
}

func ConvertStringFloatToFloat(fs string) (f float64, err error) {
	rev := reverse([]byte(fs))
	point := strings.Index(rev, ".")
	comma := strings.Index(rev, ",")

	if comma == -1 {
		f, err = strconv.ParseFloat(fs, 64)
		if err != nil {
			return
		}
	}

	if point < comma {
		f, err = strconv.ParseFloat(normalizeAmericanBritain(fs), 64)
		if err != nil {
			return
		}
	} else {
		f, err = strconv.ParseFloat(normalizeEurope(fs), 64)
		if err != nil {
			return
		}
	}
	return
}
