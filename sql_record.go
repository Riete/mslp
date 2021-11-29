package main

import (
	"container/list"
	"fmt"
	"log"
	"time"

	"github.com/riete/dingtalk"
)

type SQLSlowRecord struct {
	HostAddress        string
	DBName             string
	SQLText            string
	QueryTimes         string
	LockTimes          string
	ParseRowCounts     string
	ReturnRowCounts    string
	ExecutionStartTime string
}

func (record SQLSlowRecord) Message() string {
	executionTime, _ := time.Parse("2006-01-02T15:04:05", record.ExecutionStartTime)
	return fmt.Sprintf(
		"> 执行时间：%s\n\n> 客户端IP：%s\n\n> 数据库名：%s\n\n> 执行时长：%ss\n\n"+
			"> 锁定时长：%ss\n\n> 解析行数：%s\n\n> 返回行数：%s\n\n> SQL语句：%s\n\n",
		executionTime.Format("2006-01-02 15:04:05"),
		record.HostAddress,
		record.DBName,
		record.QueryTimes,
		record.LockTimes,
		record.ParseRowCounts,
		record.ReturnRowCounts,
		record.SQLText,
	)
}

func SendDingTalk(sl SizedList) {
	for {
		title := fmt.Sprintf("数据库[%s]新增慢SQL信息", *mySqlName)
		for i := sl.list.Front(); i != nil; i = i.Next() {
			dingtalk.SendDingTalkMarkdown(title, i.Value.(SQLSlowRecord).Message(), *webhook, *secret, false)
			sl.list.Remove(i)
			log.Println("total slow log message count to be sent is: ", sl.list.Len())
			time.Sleep(3 * time.Second)
		}
	}
}

type SizedList struct {
	list *list.List
	max  int
}

func (l *SizedList) PushBack(v interface{}) {
	if l.list.Len() < l.max {
		l.list.PushBack(v)
	} else {
		l.list.Remove(l.list.Front())
		l.list.PushBack(v)
	}
}
