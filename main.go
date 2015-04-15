package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/docopt/docopt-go"
	_ "github.com/lib/pq"
)

func parseArgs() (args map[string]interface{}, err error) {
	usage := `ury

Usage:
  ury -h
  ury -v
  ury show search <query>
  ury show seasons <showid>
  ury show timeslots <showseasonid>

Options:
  -h --help                  Show this help message
  -v --version               Show version`

	args, err = docopt.Parse(usage, nil, true, "urytools 0.0", false)
	return
}

func main() {
	logger := log.New(os.Stdout, "", log.Lshortfile)
	args, err := parseArgs()
	if err != nil {
		logger.Fatal(err)
	}
	dburl, err := ioutil.ReadFile(".urydb")
	if err != nil {
		logger.Fatal(err)
	}
	db, err := sql.Open("postgres", string(dburl))
	if err != nil {
		logger.Fatal(err)
	}
	switch {
	case args["show"].(bool) && args["search"].(bool):
		rows, err := db.Query(`SELECT m.metadata_value as title, s.show_id as id FROM schedule.show s LEFT JOIN schedule.show_metadata m ON s.show_id=m.show_id AND m.metadata_key_id=2 WHERE m.metadata_value ILIKE '%' ||$1|| '%'`, args["<query>"])
		defer rows.Close()
		if err != nil {
			logger.Fatal(err)
		}
		var title string
		var id int
		for rows.Next() {
			if err := rows.Scan(&title, &id); err != nil {
				logger.Fatal(err)
			}
			logger.Println(fmt.Sprintf("%d: %s", id, title))
		}
	case args["show"].(bool) && args["seasons"].(bool):
		rows, err := db.Query(`SELECT s.show_season_id, t.descr, t.start FROM schedule.show_season s LEFT JOIN terms t ON t.termid=s.termid WHERE s.show_id=$1`, args["<showid>"])
		defer rows.Close()
		if err != nil {
			logger.Fatal(err)
		}
		var id int
		var term string
		var start time.Time
		for rows.Next() {
			if err := rows.Scan(&id, &term, &start); err != nil {
				logger.Fatal(err)
			}
			logger.Println(fmt.Sprintf("%d: %s %d", id, term, start.Year()))
		}
	}
}
