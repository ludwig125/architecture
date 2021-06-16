# architecture

## how to run this code

```
cd cleanArchitectureWebAPI
$DB_TYPE=sqlite go run $(ls -1 *.go | grep -v _test.go)
```

参考：testを除いてgo runする方法
ref: https://stackoverflow.com/questions/23695448/how-to-run-all-go-files-within-current-directory-through-the-command-line-mult
