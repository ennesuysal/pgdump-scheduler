package main

import (
	"github.com/ennesuysal/pgdump-scheduler/parser"
	"github.com/ennesuysal/pgdump-scheduler/psql"
)

func main() {
	// c := pg.NewConnection("127.0.0.1", "5432", "postgres", "mysecretpassword")
	// pg.DumpAll(c, "test.txt")
	p, _ := parser.Parse()
	psql.DumpAllHosts(p, ".")
}
