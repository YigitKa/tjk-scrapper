package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/alexeyco/simpletable"

	jokey "hello/models"
)

// Move another file
func parseSuspendedJokey() {
	// Request the HTML page.
	res, err := http.Get("https://www.tjk.org/TR/YarisSever/Query/Page/CezaliJokey")
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	// Check status code
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: try make instead of array
	jokeys := []jokey.Jokey{}

	// TODO: need refactor
	// ajaxtbody element contains information about suspended jokey(s)
	// check each element if values are filled.
	doc.Find(".ajaxtbody").Each(func(i int, s *goquery.Selection) {
		s.Find(".sorgu-CezaliJokey-JokeyAdi").Each(func(i int, k *goquery.Selection) {
			name := strings.Trim(k.Text(), "\r\n")
			name = strings.ReplaceAll(name, "  ", "")
			if strings.TrimSpace(name) != "" {
				jokeys = append(jokeys, jokey.Jokey{Id: i + 1, Name: strings.TrimSpace(name)})
			} else {
				// TODO: if name empty. no need to continue. go next record
			}
		})
		s.Find(".sorgu-CezaliJokey-BaslangicTarihi").Each(func(i int, k *goquery.Selection) {
			startDate := strings.Trim(k.Text(), "\r\n")
			startDate = strings.ReplaceAll(startDate, "  ", "")
			if strings.TrimSpace(startDate) != "" {
				jokeys[i].SuspendDate = strings.TrimSpace(startDate)
			}
		})
		s.Find(".sorgu-CezaliJokey-BitisTarihi").Each(func(i int, k *goquery.Selection) {
			endDate := strings.Trim(k.Text(), "\r\n")
			endDate = strings.ReplaceAll(endDate, "  ", "")
			if strings.TrimSpace(endDate) != "" {
				jokeys[i].AllowDate = strings.TrimSpace(endDate)
			}
		})
		s.Find(".sorgu-CezaliJokey-CezaNedeni").Each(func(i int, k *goquery.Selection) {
			description := strings.Trim(k.Text(), "\r\n")
			description = strings.ReplaceAll(description, "  ", "")
			if strings.TrimSpace(description) != "" {
				jokeys[i].Description = wordWrap(strings.TrimSpace(description), 50)
			}
		})
		table := simpletable.New()

		table.Header = &simpletable.Header{
			Cells: []*simpletable.Cell{
				{Align: simpletable.AlignCenter, Text: "#"},
				{Align: simpletable.AlignCenter, Text: "Jokey Adı"},
				{Align: simpletable.AlignCenter, Text: "Ceza Başlangıç"},
				{Align: simpletable.AlignCenter, Text: "Ceza Bitiş"},
				{Align: simpletable.AlignCenter, Text: "Nedeni"},
			},
		}

		for row := range jokeys {
			r := []*simpletable.Cell{
				{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%d", jokeys[row].Id)},
				{Text: jokeys[row].Name},
				{Align: simpletable.AlignLeft, Text: jokeys[row].SuspendDate},
				{Align: simpletable.AlignLeft, Text: jokeys[row].AllowDate},
				{Align: simpletable.AlignLeft, Span: 1, Text: jokeys[row].Description},
			}

			table.Body.Cells = append(table.Body.Cells, r)
		}

		table.SetStyle(simpletable.StyleDefault)
		fmt.Println(table.String())
	})
}

// TODO: move another file
func wordWrap(text string, lineWidth int) (wrapped string) {
	words := strings.Fields(strings.TrimSpace(text))
	if len(words) == 0 {
		return text
	}
	wrapped = words[0]
	spaceLeft := lineWidth - len(wrapped)
	for _, word := range words[1:] {
		if len(word)+1 > spaceLeft {
			wrapped += "\n" + word
			spaceLeft = lineWidth - len(word)
		} else {
			wrapped += " " + word
			spaceLeft -= 1 + len(word)
		}
	}
	return
}

// TODO check args
func main() {
	parseSuspendedJokey()
}
