package junitexport

import (
	"encoding/xml"
	"os"
	"qa-int-automation/tester"
)

type Testsuites struct {
	XMLName    xml.Name    `xml:"testsuites"`
	Testsuites []Testsuite `xml:"testsuite"`
}

type Testsuite struct {
	Name      string     `xml:"name,attr"`
	Tests     int        `xml:"tests,attr"`
	Failures  int        `xml:"failures,attr"`
	Testcases []Testcase `xml:"testcase"`
}

type Testcase struct {
	Name    string   `xml:"name,attr"`
	Class   string   `xml:"classname,attr"`
	Time    string   `xml:"time,attr"`
	Failure *Failure `xml:"failure,omitempty"`
}

type Failure struct {
	Message string `xml:"message,attr"`
	Type    string `xml:"type,attr"`
	Body    string `xml:",chardata"`
}

func ExportToJUnit(results []tester.Result, filename string) error {
	var testcases []Testcase
	failures := 0

	for _, r := range results {
		tc := Testcase{
			Name:  r.Name,
			Class: "qa-int-automation",
			Time:  "0.1", // Default time, can be enhanced
		}
		if !r.Passed {
			failures++
			tc.Failure = &Failure{
				Message: r.Message,
				Type:    "Failure",
				Body:    r.Response,
			}
		}
		testcases = append(testcases, tc)
	}

	ts := Testsuite{
		Name:      "QA Integration Tests",
		Tests:     len(results),
		Failures:  failures,
		Testcases: testcases,
	}

	suites := Testsuites{Testsuites: []Testsuite{ts}}

	if err := os.MkdirAll("results", 0755); err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := xml.NewEncoder(file)
	encoder.Indent("", "  ")
	return encoder.Encode(suites)
}
