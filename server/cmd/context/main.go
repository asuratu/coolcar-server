package main

import (
    "context"
    "fmt"
    "time"
)

type paramKey struct {
}

func main() {
    // 创建一个带有key-value的context
    c := context.WithValue(context.Background(), paramKey{}, "abc")
    // 设置context在5秒之后结束
    c, cancel := context.WithTimeout(c, 10*time.Second)
    defer cancel()

    // 不加go的话，会等待mainTask执行完毕，才会执行下面的代码
    go mainTask(c)

    var cmd string
    for {
        fmt.Scan(&cmd)
        if cmd == "c" {
            cancel()
        }
    }
}

func mainTask(c context.Context) {
    fmt.Printf("mainTask started with param %q\n", c.Value(paramKey{}))
    // 这个是启动后台任务的正确做法
    go func() {
        c1, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()
        smallTask(c1, "task1", 9*time.Second)
    }()
    smallTask(c, "task2", 8*time.Second)
}

func smallTask(c context.Context, name string, d time.Duration) {
    fmt.Printf("%s started with param %q\n", name, c.Value(paramKey{}))
    select {
    // 如果context被取消，则会跳出select，被取消就会从c.Done()收到值
    case <-c.Done():
        fmt.Printf("%s canceled\n", name)
    // 如果context超时，则会跳出select
    case <-time.After(d):
        fmt.Printf("%s finished\n", name)
    }
}
