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

func dumpFull(c *Connection, dbName string, basePath string, mode string) error {
	//fmt.Printf("%s Job runned for %s!\n", mode, dbName)
	cmdStr := []string{"-h", c.host, "-p", c.port, "-U", c.username}
	cmd := exec.Command("/usr/local/bin/pg_dumpall", cmdStr...)
	cmd.Env = append(cmd.Env, "PGPASSWORD="+c.password)

	outFile, err := os.Create(basePath + "/" + c.serviceName + "/" + dbName + ".sql")
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
		err := os.MkdirAll(basePath+"/"+service, os.ModePerm)
		if err != nil {
			return err
		}
		conn := NewConnection(db["host"].(string), db["port"].(string), db["username"].(string), db["password"].(string), service, db["dbs"].(map[string]interface{}))
		for dbName, attr := range conn.dbs {
			for mode, cron := range attr.(map[string]interface{}) {
				mode := mode
				c.AddFunc(cron.(string), func() {
					fmt.Println(mode)
					dumpFull(conn, dbName, basePath, mode)
				})
			}
		}
	}
	c.Start()
	time.Sleep(3600 * time.Second)
	return nil
}
