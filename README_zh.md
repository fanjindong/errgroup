# errgroup

`errgroup`包是`sync/errgroup`包的扩展，其核心是管理一组`gotoutine`的行为。当发生`error`or`panic`时，控制组内其他`goroutine`的行为： 取消所有执行、继续所有执行。

![image](./images/BatEaredFox_EN-AU12936466242_1920x1080.jpg)

## Install

`go get -u github.com/fanjindong/errgroup`

## Fast Start

`errgroup` 包含三种常用方式

1. `NewContinue` 此时不会因为一个任务失败导致所有任务被 cancel

```go
eg := &errgroup.NewContinue(ctx)
eg.Go(func (ctx context.Context) {
    // NOTE: NOTE: 此时 ctx 为 NewContinue 传递的 ctx
    // do something
})
eg.Wait()
```

2. `NewCancel` 此时如果有一个任务失败会导致所有**未进行或进行中**的任务被 cancel

```go
eg := errgroup.NewCancel(ctx)
eg.Go(func (ctx context.Context) {
    // NOTE: 此时 ctx 是从 WithContext 传递的 ctx 派生出的 ctx
    // do something
})
eg.Wait()
```

3. 设置最大并行数 SetMaxProcess 对以上使用方式均起效

```go
eg := errgroup.NewCancel(ctx, WithMaxProcess(2))
// task1
eg.Go(func(ctx context.Context) {
    fmt.Println("task1")
})
// task2
eg.Go(func (ctx context.Context) {
    fmt.Println("task2")
})
// task3
eg.Go(func (ctx context.Context) {
    fmt.Println("task3")
})
eg.Wait()
// NOTE: 此时设置的 WithMaxProcess 为 2, 
// 添加了三个任务 task1, task2, task3 ,
// 最初只有2个task运行，直到一个task完成，最后一个task才会开始执行。
```


