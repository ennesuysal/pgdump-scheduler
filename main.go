package main

import (
	"github.com/ennesuysal/pgdump-scheduler/parser"
	"github.com/ennesuysal/pgdump-scheduler/psql"
)

func main() {
	p, err := parser.Parse()
	if err != nil {
		panic(err)
	}

	psql.DumpAllHosts(p, ".")
}
