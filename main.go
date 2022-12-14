package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/alexeyco/simpletable"

	helper "tjk-scrapper/helpers"
	jokey "tjk-scrapper/models"
)

func parseSuspendedJokey(jsonOutPtr *bool, tableOutPtr *bool) {
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

	jokeys := make([]jokey.Jokey, 0)

	// ajaxtbody element contains information about suspended jokey(s)
	// check each element if values are filled.
	doc.Find(".ajaxtbody").Each(func(i int, s *goquery.Selection) {
		s.Find(".sorgu-CezaliJokey-JokeyAdi").Each(func(i int, k *goquery.Selection) {
			name := strings.Trim(k.Text(), "\r\n")
			name = strings.ReplaceAll(name, "  ", "")
			if strings.TrimSpace(name) != "" {
				jokeys = append(jokeys, jokey.Jokey{Id: i + 1, Name: strings.TrimSpace(name)})
			} else {
				s.Next()
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
				jokeys[i].DueDate = strings.TrimSpace(endDate)
			}
		})
		s.Find(".sorgu-CezaliJokey-CezaNedeni").Each(func(i int, k *goquery.Selection) {
			description := strings.Trim(k.Text(), "\r\n")
			description = strings.ReplaceAll(description, "  ", "")
			if strings.TrimSpace(description) != "" {
				jokeys[i].BanReason = helper.WordWrap(strings.TrimSpace(description), 50)
			}
		})

		if *tableOutPtr {
			table := simpletable.New()

			table.Header = &simpletable.Header{
				Cells: []*simpletable.Cell{
					{Align: simpletable.AlignCenter, Text: "#"},
					{Align: simpletable.AlignCenter, Text: "Jokey Ad??"},
					{Align: simpletable.AlignCenter, Text: "Ceza Ba??lang????"},
					{Align: simpletable.AlignCenter, Text: "Ceza Biti??"},
					{Align: simpletable.AlignCenter, Text: "Nedeni"},
				},
			}

			for row := range jokeys {
				r := []*simpletable.Cell{
					{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%d", jokeys[row].Id)},
					{Text: jokeys[row].Name},
					{Align: simpletable.AlignLeft, Text: jokeys[row].SuspendDate},
					{Align: simpletable.AlignLeft, Text: jokeys[row].DueDate},
					{Align: simpletable.AlignLeft, Span: 1, Text: jokeys[row].BanReason},
				}

				table.Body.Cells = append(table.Body.Cells, r)
			}

			table.SetStyle(simpletable.StyleDefault)
			fmt.Println(table.String())
		}

		if *jsonOutPtr {
			fmt.Println(jokey.ToJson(jokeys))
		}
	})
}

func main() {
	jsonOutPtr := flag.Bool("json", false, "a bool")
	tableOutPtr := flag.Bool("table", false, "a bool")

	flag.Parse()
	if !*jsonOutPtr && !*tableOutPtr {
		fmt.Println("no parameters given. try --json or/and --table")
		os.Exit(1)
	}

	parseSuspendedJokey(jsonOutPtr, tableOutPtr)
	os.Exit(0) // sussess exit
}
