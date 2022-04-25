package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

// Chapter defines the regex expressions to search for that denote the start and end of a
// chapter (e.g. Street, Street Touring) of the rulebook.
type Chapter struct {
	name        string
	number      string
	subChapters []string
	start       *regexp.Regexp
	end         *regexp.Regexp
}

// getSubChapters returns an array of sub chapters (e.g. 13.1, 13.2) that exist for a given
// chapter
func getSubChapters(rules, chapterNumber string) []string {
	subChapters := []string{}
	regexString := chapterNumber + `\.([0-9]+[.A-Z]*).*\.+[. ]([0-9]+)`
	tableOfContents := regexp.MustCompile(regexString)
	match := tableOfContents.FindAllStringSubmatch(rules, -1)
	// This means there probably aren't any subchapters
	if len(match) < 2 {
		return []string{}
	}
	for i := range match {
		subChapters = append(subChapters, fmt.Sprintf("%s.%s", chapterNumber, match[i][1]))
	}
	return subChapters
}

// readFile returns a Reader of a specific file
func readFile() *strings.Reader {
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

	return strings.NewReader(rulesString)
}

func getChapter(rules *strings.Reader, chapter Chapter) *io.SectionReader {
	rules.Seek(0, 0)
	startMatch := chapter.start.FindReaderIndex(rules)
	rules.Seek(0, 0)
	endMatch := chapter.end.FindReaderIndex(rules)
	rules.Seek(0, 0)
	length := endMatch[0] - startMatch[0]
	return io.NewSectionReader(rules, int64(startMatch[0]), int64(length))
}

func main() {
	rules := readFile()

	rulesBytes, err := io.ReadAll(rules)
	if err != nil {
		log.Fatal(err)
	}
	rules.Seek(0, 0)
	allChapters := []Chapter{
		{
			name:   "Street",
			number: "13",
			start:  regexp.MustCompile(`\n13\. STREET CATEGORY\n`),
			end:    regexp.MustCompile(`\n14\. STREET TOURING® CATEGORY\n`),
		},
		{
			name:   "Street Touring",
			number: "14",
			start:  regexp.MustCompile(`\n14\. STREET TOURING® CATEGORY\n`),
			end:    regexp.MustCompile(`\n15\. STREET PREPARED CATEGORY\n`),
		},
		{
			name:   "Street Prepared",
			number: "15",
			start:  regexp.MustCompile(`\n15\. STREET PREPARED CATEGORY\n`),
			end:    regexp.MustCompile(`\n16\. STREET MODIFIED CATEGORY\n`),
		},
		{
			name:   "Street Modified",
			number: "16",
			start:  regexp.MustCompile(`\n16\. STREET MODIFIED CATEGORY\n`),
			end:    regexp.MustCompile(`\n17\. PREPARED CATEGORY\n`),
		},
		{
			name:   "Prepared",
			number: "17",
			start:  regexp.MustCompile(`\n17\. PREPARED CATEGORY\n`),
			end:    regexp.MustCompile(`\n18\. MODIFIED CATEGORY\n`),
		},
		{
			name:   "Modified",
			number: "18",
			start:  regexp.MustCompile(`\n18\. MODIFIED CATEGORY\n`),
			end:    regexp.MustCompile(`\n19\. KART CATEGORY\n`),
		},
		{
			name:   "Solo Spec Coupe",
			number: "20",
			start:  regexp.MustCompile(`\n20\. SOLO® SPEC COUPE \(SSC\)\n`),
			end:    regexp.MustCompile(`\n21\. PROSOLO® NATIONAL SERIES RULES\n`),
		},
		{
			name:   "Extreme Street",
			number: "n/a",
			start:  regexp.MustCompile(`\nEXTREME STREET \(XS\)\n`),
			end:    regexp.MustCompile(`\nAPPENDIX C - SOLO® ROLL BAR STANDARDS\n`),
		},
	}
	for i := range allChapters {
		if allChapters[i].number != "n/a" {
			subChapters := getSubChapters(string(rulesBytes), allChapters[i].number)
			fmt.Println(subChapters)
		}
		sectionReader := getChapter(rules, allChapters[i])
		rulesBytes, err := io.ReadAll(sectionReader)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(allChapters[i].name)
		fmt.Println(string(rulesBytes))
		fmt.Println()
	}
}
