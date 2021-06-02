package errgroup

import (
	"context"
	"fmt"
)

func fakeRunTask(ctx context.Context) error {
	return nil
}

func ExampleGroup_group() {
	g := Group{}
	g.Go(func(context.Context) error {
		return fakeRunTask(context.Background())
	})
	g.Go(func(context.Context) error {
		return fakeRunTask(context.Background())
	})
	if err := g.Wait(); err != nil {
		// handle err
	}
}

func ExampleGroup_ctx() {
	g := WithContext(context.Background())
	g.Go(func(ctx context.Context) error {
		return fakeRunTask(ctx)
	})
	g.Go(func(ctx context.Context) error {
		return fakeRunTask(ctx)
	})
	if err := g.Wait(); err != nil {
		// handle err
	}
}

func ExampleGroup_cancel() {
	g := WithCancel(context.Background())
	g.Go(func(ctx context.Context) error {
		return fakeRunTask(ctx)
	})
	g.Go(func(ctx context.Context) error {
		return fakeRunTask(ctx)
	})
	if err := g.Wait(); err != nil {
		// handle err
	}
}

func ExampleGroup_maxproc() {
	g := Group{}
	// set max concurrency
	g.SetMaxProcess(2)
	g.Go(func(ctx context.Context) error {
		return fakeRunTask(context.Background())
	})
	g.Go(func(ctx context.Context) error {
		return fakeRunTask(context.Background())
	})
	if err := g.Wait(); err != nil {
		// handle err
	}
}

func ExampleGroup_demo() {
	c := context.Background()
	eg := WithContext(c)

	var res1 map[string]string
	eg.Go(func(ctx context.Context) error {
		// do something
		// 结果绑定至指针变量
		res1["message"] = "success"
		return nil
	})

	type S struct {
		Code int
	}
	var res2 *S
	eg.Go(func(ctx context.Context) error {
		// do something
		// 结果绑定至指针变量
		res2 = &S{Code: 1}
		return nil
	})

	if err := eg.Wait(); err != nil {
		panic(err)
	}

	fmt.Sprintln(res1["message"]) // out: success
	fmt.Sprintln(res2.Code)       // out: 1
	return
}
