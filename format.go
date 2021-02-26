package main

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

// GetNextVersion gets next version according to the tag of the format.
func GetNextVersion(tag string) (string, error) {
	if tag == "" {
		return "v1.0.0", nil
	}

	tags := strings.Split(tag, ".")
	len := len(tags)

	// semantic version(e.g. v1.2.3)
	if len > 2 {
		patch, err := strconv.Atoi(tags[len-1])
		if err != nil {
			return "", err
		}

		tags[len-1] = strconv.Itoa(patch + 1)
		return strings.Join(tags, "."), nil
	}

	// date version(e.g. 20180525.1)
	const layout = "20060102"
	today := time.Now().Format(layout)

	dateRe := regexp.MustCompile(`(.*)(\d{8})(\.)(.+)`)
	if match := dateRe.FindStringSubmatch(tag); match != nil {
		if match[2] == today {
			minor, err := strconv.Atoi(match[4])
			if err != nil {
				return "", err
			}

			next := strconv.Itoa(minor + 1)
			return match[1] + today + "." + next, nil
		}
		return match[1] + today + "." + "1", nil
	}
	return today + ".1", nil
}

// GetReleaseNote formats merge-commits list.
func GetReleaseNote(tag string, list string) string {
	return "Release " + tag + "\n\n" + "## " + tag + "\n" + list
}
