package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/logrusorgru/aurora"

	"github.com/sbowman/pgstats"
)

// DB defines the environment variable expected with the database URI.
const DB = "DB_URI"

func main() {
	if len(os.Args) < 2 {
		help()
		os.Exit(1)
	}

	uri := os.Getenv(DB)
	if uri == "" {
		_, _ = fmt.Fprintf(os.Stderr, "Please set the %s environment variable to the PostgreSQL connection string,\ne.g. DB=postgres://postgres@localhost/mydb\n", DB)
		os.Exit(1)
	}

	db, err := sqlx.Connect("postgres", uri)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Unable to connect to the database: %s\n", err)
		os.Exit(1)
	}

	cmd, ok := pgstats.Commands[os.Args[1]]
	if !ok {
		help()
		os.Exit(1)
	}

	results, err := cmd.Fn(db, os.Args[2:]...)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	doc, _ := json.MarshalIndent(results, "", "  ")
	fmt.Println(string(doc))
}

func help() {
	fmt.Println()
	fmt.Println(aurora.BrightRed("Please provide one of the following commands:"))

	var keys []string
	for key := range pgstats.Commands {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		cmd := pgstats.Commands[key]

		fmt.Println()
		fmt.Println(aurora.Underline(aurora.Green(fmt.Sprintf("%s", cmd.Name))))
		fmt.Println()
		fmt.Println(aurora.Faint(aurora.Italic(wrap(cmd.Description, 78))))
	}

	fmt.Println()
}

func wrap(value string, width int) string {
	var lines []string
	var current strings.Builder

	words := strings.Split(value, " ")
	for _, word := range words {
		if current.Len() == 0 && len(word) > width {
			width = len(word)
			lines = append(lines, word)
			continue
		}

		if current.Len()+len(word) > width {
			lines = append(lines, current.String())
			current.Reset()
			current.WriteString(word)
			continue
		}

		if strings.HasSuffix(word, "\n\n") {
			current.WriteString(word[:len(word)-2])
			lines = append(lines, current.String())
			lines = append(lines, "")
			current.Reset()
			continue
		}

		if current.Len() > 0 {
			current.WriteString(" ")
		}

		current.WriteString(word)
	}

	if current.Len() > 0 {
		lines = append(lines, current.String())
	}

	return strings.Join(lines, "\n")
}
