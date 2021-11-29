# MSLP
Tool For Parse MySQL Slow Log File And Send Via DingTalk
# Usage
```
./mslp --help
Usage of ./mslp:
  -exclude-db string
        db1,db2,... exclude these databases slow log, seperated by ','
  -log-path string
        mysql slow log file path
  -max-message-delay int
        the max num of delay message if there are many slow log (default 100)
  -only-db string
        db1,db2,... only send these databases slow log, seperated by ','
  -secret string
        dingtalk webhook secret
  -sql-name string
        mysql name
  -webhook string
        dingtalk webhook
```
