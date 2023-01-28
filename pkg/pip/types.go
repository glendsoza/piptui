package pip

import (
	"regexp"
	"strings"
)

var SPLIT_FILED_REGEX = regexp.MustCompile("(?P<field>.*?):(?P<value>.*)")

type Packages []*Package

type Package struct {
	raw         string
	Name        string
	Version     string
	Summary     string
	HomePage    string
	Author      string
	AuthorEmail string
	License     string
	Location    string
	Requires    []string
	RequiredBy  []string
}

func (m *Package) ToString() string {
	return m.raw
}

func ParseRawData(data string) (*Package, error) {
	m := &Package{raw: strings.TrimSpace(data)}
	field := ""
	value := ""
	for _, line := range strings.Split(m.raw, "\n") {
		fieldValuePair := SPLIT_FILED_REGEX.FindStringSubmatch(line)
		if len(fieldValuePair) != 3 {
			// handle multiline
			value = line + "\n"
		} else {
			field = strings.TrimSpace(fieldValuePair[1])
			value = strings.TrimSpace(fieldValuePair[2])
		}
		switch field {
		case "Name":
			m.Name += value
		case "Version":
			m.Version += value
		case "Summary":
			m.Summary += value
		case "Home-Page":
			m.HomePage += value
		case "Author":
			m.Author += value
		case "Author-email":
			m.AuthorEmail += value
		case "License":
			m.License += value
		case "Location":
			m.Location += value
		case "Requires":
			for _, v := range strings.Split(strings.TrimSpace(value), ",") {
				if trimmedV := strings.TrimSpace(v); trimmedV != "" {
					m.Requires = append(m.Requires, trimmedV)
				}
			}
		case "Required-by":
			for _, v := range strings.Split(strings.TrimSpace(value), ",") {
				if trimmedV := strings.TrimSpace(v); trimmedV != "" {
					m.RequiredBy = append(m.RequiredBy, trimmedV)
				}
			}
		}
	}
	return m, nil
}
