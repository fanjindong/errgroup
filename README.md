# errgroup

[中文文档](./README_zh.md)

The `errgroup` package is an extension of the `sync/errgroup` package, and its core is to manage the behavior of a set of `gotoutines`. 
Control the behavior of other `goroutines` in the group when `error` or `panic` occurs: cancel all executions or continue all executions.

![image](./images/BatEaredFox_EN-AU12936466242_1920x1080.jpg)

## Install

`go get -u github.com/fanjindong/errgroup`

## Fast Start

`errgroup` contains three common methods

1. `NewContinue` At this point, all tasks will not be cancelled because one task failed

```go
eg := &errgroup.NewContinue(ctx)
eg.Go(func (ctx context.Context) {
    // NOTE: In this case, ctx is the ctx passed by NewContinue
    // do something
})
eg.Wait()
```

2. `NewCancel` If one task fails, **all pending or ongoing tasks** will be cancelled

```go
eg := errgroup.NewCancel(ctx)
eg.Go(func (ctx context.Context) {
    // NOTE: At this point ctx is derived from the ctx passed in WithContext
    // do something
})
eg.Wait()
```

3. Setting the maximum number of parallelism SetMaxProcess works for all the above uses

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
// NOTE: At this point, WithMaxProcess is set to 2,
// Added task1, task2, task3 ,
// Initially, only two tasks will run, and the last task will not be executed until one task is complete.
```


