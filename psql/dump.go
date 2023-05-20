package psql

import (
	"fmt"
	"os"
	"os/exec"
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

func dumpFull(c *Connection, basePath string) error {
	cmdStr := []string{"-h", c.host, "-p", c.port, "-U", c.username}
	cmd := exec.Command("/usr/local/bin/pg_dumpall", cmdStr...)
	cmd.Env = append(cmd.Env, "PGPASSWORD="+c.password)

	for dbName, attr := range c.dbs {
		for db, values := range attr.(map[string]interface{}) {
			fmt.Println(db, values)
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
	}

	return nil

}

func DumpAllHosts(dbs map[string]map[string]interface{}, basePath string) error {

	for service, db := range dbs {
		err := os.MkdirAll(basePath+"/"+service, os.ModePerm)
		if err != nil {
			return err
		}
		dumpFull(NewConnection(db["host"].(string), db["port"].(string), db["username"].(string), db["password"].(string), service, db["dbs"].(map[string]interface{})), basePath)
	}
	return nil
}
