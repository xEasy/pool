autopool:
    自动伸缩池，指定最小和最大连接，按需增加

workerpool:
    工人令牌池，指定工人数和队列，池创建就已初始化工人数

Usage:

```golang

import (
  "fmt"
	"github.com/xEasy/pool/workerpool"
)

pool := workerpool.NewPool(WORKER_COUNT, JOB_QUEUE_LEGTH)

func helloWorld(w string) {
  fmt.Print("hello", w)
}

pool.WaitCount(1)
pool.Enqueue(helloWorld, "world")
pool.WaitAll()

# release workpool
pool.Release()

```
