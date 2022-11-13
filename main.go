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
	subChapters []SubChapter
	reader      *io.SectionReader
	start       *regexp.Regexp
	end         *regexp.Regexp
}

// SubChapter holds the name, number, and body of a subchapter of the rules (e.g. 13.2 Bodywork)
type SubChapter struct {
	name   string
	number string
	reader *io.SectionReader
}

// getSubChapters returns an array of sub chapters (e.g. 13.1, 13.2) that exist for a given
// chapter
func getSubChapters(rules, chapterNumber string) []SubChapter {
	subChapters := []SubChapter{}
	regexString := chapterNumber + `\.([0-9]+[.A-Z]*) ([^\.\n]*)\.+[\. ]([0-9]+)`
	tableOfContents := regexp.MustCompile(regexString)
	match := tableOfContents.FindAllStringSubmatch(rules, -1)
	// This means there probably aren't any subchapters
	if len(match) < 2 {
		return subChapters
	}
	for i := range match {
		subChapters = append(subChapters,
			SubChapter{
				number: fmt.Sprintf("%s.%s", chapterNumber, match[i][1]),
				name:   match[i][2],
			},
		)
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

	// standardize double quotes
	rulesString = strings.ReplaceAll(rulesString, "“", `"`)
	rulesString = strings.ReplaceAll(rulesString, "”", `"`)

	remove := regexp.MustCompile(`\n\f`)
	rulesString = string(remove.ReplaceAll([]byte(rulesString), []byte{}))

	return strings.NewReader(rulesString)
}

// findSubChapterBody populates the reader field of each subchapter with the body
// of that subchapter
func findSubChapterBody(chapter Chapter, chapterText []byte) []SubChapter {
	subChapters := chapter.subChapters
	reader := strings.NewReader(string(chapterText))
	for i, subChapter := range subChapters {
		reader.Seek(0, 0)
		var length int
		seekToEnd := false
		if i == len(subChapters)-1 {
			seekToEnd = true
		}
		startRegexString := `(?i)` + regexp.QuoteMeta(subChapter.number) + ` ` + regexp.QuoteMeta(subChapter.name)
		startRegex := regexp.MustCompile(startRegexString)
		startMatch := startRegex.FindReaderIndex(reader)
		reader.Seek(0, 0)
		if startMatch != nil {
			if seekToEnd == true {
				length = reader.Len() // jump to end of reader, this is the last element
			} else {
				endRegexString := `(?i)` + regexp.QuoteMeta(subChapters[i+1].number) + ` ` + regexp.QuoteMeta(subChapters[i+1].name)
				endRegex := regexp.MustCompile(endRegexString)
				endMatch := endRegex.FindReaderIndex(reader)
				reader.Seek(0, 0)
				if endMatch != nil {
					length = endMatch[0] - startMatch[0]
				}
			}
			sectionReader := io.NewSectionReader(reader, int64(startMatch[0]), int64(length))
			subChapters[i].reader = sectionReader
			subchapter, err := io.ReadAll(sectionReader)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(subChapter.number + " " + subChapter.name)
			fmt.Println(string(subchapter))
		}
	}
	return subChapters
}

func getChapterReader(rules *strings.Reader, chapter Chapter) *io.SectionReader {
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
			allChapters[i].subChapters = subChapters
		}
		chapterReader := getChapterReader(rules, allChapters[i])
		allChapters[i].reader = chapterReader

		chapterText, err := io.ReadAll(chapterReader)
		if err != nil {
			fmt.Println("error reading chapter text")
			os.Exit(1)
		}

		// remove all form feed (i.e. ) chapter title lines
		if allChapters[i].number != "n/a" {
			remove := regexp.MustCompile(`\n\f` + allChapters[i].number + `\. .+\n`)
			chapterText = remove.ReplaceAll(chapterText, []byte{})
		}

		// remove all page number text
		remove := regexp.MustCompile(`(?i)([0-9]+ — )*202[0-2] SCCA® NATIONAL SOLO® RULES( )*(— [0-9]+)*`)
		chapterText = remove.ReplaceAll(chapterText, []byte{})

		if allChapters[i].number != "n/a" && len(allChapters[i].subChapters) > 0 {
			allChapters[i].subChapters = findSubChapterBody(allChapters[i], chapterText)
		}
	}
	fmt.Printf("%+v", allChapters)
}
