package psql

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/robfig/cron"
)

type Connection struct {
	host        string
	port        string
	username    string
	password    string
	serviceName string
	dbs         map[string]interface{}
}

func NewConnection(host string, port string, username string, password string, serviceName string, dbs map[string]interface{}) *Connection {
	return &Connection{
		host:        host,
		port:        port,
		username:    username,
		password:    password,
		serviceName: serviceName,
		dbs:         dbs,
	}
}

func dumpFull(c *Connection, dbName string, basePath string, mode string, tableName string) error {
	fmt.Println("Values: ", dbName, basePath, mode, tableName)
	cmdStr := []string{"-h", c.host, "-p", c.port, "-U", c.username, "-d", dbName}

	baseFName := basePath + "/" + c.serviceName + "/" + mode
	fileName := baseFName + "/" + dbName
	if mode == "table" {
		cmdStr = append(cmdStr, "--table", tableName)
		fileName += "_" + tableName
	}

	fileName += "_" + time.Now().Format("07-09-2017") + ".sql"

	cmd := exec.Command("/usr/local/bin/pg_dump", cmdStr...)
	cmd.Env = append(cmd.Env, "PGPASSWORD="+c.password)

	err := os.MkdirAll(baseFName, os.ModePerm)
	if err != nil {
		return err
	}

	outFile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer outFile.Close()

	cmd.Stdout = outFile
	cmd.Stderr = os.Stdout

	err = cmd.Run()

	if err != nil {
		return err
	}

	return nil
}

func DumpAllHosts(dbs map[string]map[string]interface{}, basePath string) error {
	c := cron.New()
	for service, db := range dbs {

		conn := NewConnection(db["host"].(string), db["port"].(string), db["username"].(string), db["password"].(string), service, db["dbs"].(map[string]interface{}))
		var tName string
		for dbName, attr := range conn.dbs {
			dbName := dbName
			attr := attr
			for mode, args := range attr.(map[string]interface{}) {
				mode := mode
				args := args
				if tn, ok := args.(map[string]interface{})["table-name"]; ok {
					tName = tn.(string)
				} else {
					tName = ""
				}
				tName := tName
				c.AddFunc(args.(map[string]interface{})["cron"].(string), func() {
					dumpFull(conn, dbName, basePath, mode, tName)
				})
			}
		}
	}
	c.Start()
	for {
		time.Sleep(3600 * time.Second)
	}
	return nil
}
