package main

import (
	"container/list"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/riete/frbl"
)

var (
	logPath         *string
	mySqlName       *string
	secret          *string
	webhook         *string
	excludeDb       *string
	onlyDb          *string
	maxMessageDelay *int
)

func flagParse() {
	logPath = flag.String("log-path", "", "mysql slow log file path")
	mySqlName = flag.String("mysql-name", "", "mysql name to set dingtalk title")
	secret = flag.String("secret", "", "dingtalk webhook secret")
	webhook = flag.String("webhook", "", "dingtalk webhook")
	onlyDb = flag.String("only-db", "", "db1,db2,... only send these databases slow log, seperated by ','")
	excludeDb = flag.String("exclude-db", "", "db1,db2,... exclude these databases slow log, seperated by ','")
	maxMessageDelay = flag.Int("max-message-delay", 100, "the max num of delay message if there are many slow log")
	flag.Parse()
	if *logPath == "" {
		panic("log-path is required")
	}
	if *mySqlName == "" {
		panic("sql-name is required")
	}
	if *secret == "" {
		panic("secret is required")
	}
	if *webhook == "" {
		panic("webhook is required")
	}
	flag.Parse()
}

func read(file frbl.FileReader) {
	for {
		if err := file.ReadLine(); err != nil {
			log.Println(err)
		}
		time.Sleep(time.Second)
	}
}

func logParse(file frbl.FileReader, sl SizedList) {
	var record SQLSlowRecord
	for c := range file.Content() {
		if strings.HasPrefix(c, "Tcp port") ||
			strings.HasPrefix(c, "SET timestamp") ||
			strings.HasPrefix(c, "Time                 Id") ||
			strings.HasSuffix(c, "started with:\n") {
			continue
		}
		if strings.HasPrefix(c, "# Time:") {
			record = SQLSlowRecord{DBName: record.DBName}
			record.ExecutionStartTime = strings.Split(strings.ReplaceAll(c, "# Time: ", ""), ".")[0]
		} else if strings.HasPrefix(c, "# User") {
			user := strings.Split(c, " ")
			record.HostAddress = fmt.Sprintf("%s %s %s %s", user[2], user[3], user[4], user[5])
		} else if strings.HasPrefix(c, "# Query_time") {
			data := strings.Split(c, " ")
			record.QueryTimes = data[2]
			record.LockTimes = data[5]
			record.ReturnRowCounts = data[7]
			record.ParseRowCounts = data[10]
		} else {
			if strings.HasPrefix(strings.ToLower(c), "use") {
				record.DBName = strings.ReplaceAll(strings.Split(c, " ")[1], ";\n", "")
			} else {
				record.SQLText += c
				if strings.HasSuffix(c, ";\n") {
					if checkDB(record.DBName) {
						sl.PushBack(record)
					}
				}
			}
		}
	}
}

func checkDB(db string) bool {
	if *onlyDb != "" {
		onlyDbs := strings.Split(*onlyDb, ",")
		return stringInA(db, onlyDbs)
	} else if *excludeDb != "" {
		excludeDbs := strings.Split(*excludeDb, ",")
		return !stringInA(db, excludeDbs)
	} else {
		return true
	}
}

func stringInA(s string, a []string) bool {
	check := make(map[string]bool)
	for _, i := range a {
		check[i] = true
	}
	return check[s]
}

func main() {
	flagParse()
	file := frbl.NewFileReader(*logPath)
	sl := SizedList{list: list.New(), max: *maxMessageDelay}
	go logParse(file, sl)
	go SendDingTalk(sl)
	read(file)
}
