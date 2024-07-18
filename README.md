# rotatefile
a simple Go log rotation library

一个简单的 Go 日志轮换库。

## Install 安装

```bash
go get -u github.com/baagod/rotatefile
```

## Example 用例

```go
file, _ := New("logs/day.log", PerMinute)

for i := 0; i < 1000; i++ {
    now := time.Now().Format(time.DateTime)
    _, _ = file.Write([]byte(fmt.Sprintf("%d: %s\n", i, now)))
    fmt.Println(now)
    time.Sleep(time.Second)
}
```
