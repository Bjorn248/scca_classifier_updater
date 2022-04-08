package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func findAppendixPage(rules, class string) int {
	pageNumber := 0
	class = strings.ReplaceAll(class, "(", `\(`)
	class = strings.ReplaceAll(class, ")", `\)`)
	regexString := class + `.*\.+[. ]([0-9]+)`
	tableOfContents := regexp.MustCompile(regexString)
	match := tableOfContents.FindStringSubmatch(rules)
	if len(match) < 2 {
		fmt.Println("Could not find page for " + class + " in table of contents")
		os.Exit(1)
	}
	pageNumber, err := strconv.Atoi(match[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return pageNumber
}

// readFile returns the contents of a file as a string
func readFile() string {
	filePath := "rules.txt"
	rules, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Encountered error reading file", filePath)
		os.Exit(1)
	}

	rulesString := string(rules)

	// Unsure why but the pdftotext output contains these characters
	// perhaps due to incorrect parsing?
	rulesString = strings.ReplaceAll(rulesString, "ﬀ", "ff")

	return rulesString
}

func main() {
	rules := readFile()
	allClasses := map[string]int{
		"13. STREET CATEGORY":          0,
		"14. STREET TOURING® CATEGORY": 0,
		"15. STREET PREPARED CATEGORY": 0,
		"16. STREET MODIFIED CATEGORY": 0,
		"17. PREPARED CATEGORY":        0,
		"18. MODIFIED CATEGORY":        0,
		"20. SOLO® SPEC COUPE (SSC)":   0,
		"EXTREME STREET (XS)":          0,
	}

	for class := range allClasses {
		allClasses[class] = findAppendixPage(rules, class)
	}
	fmt.Println(allClasses)
}
