package errgroup

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestErrGroup_Wait(t *testing.T) {
	var number int32
	var ctx context.Context = context.Background()

	type args struct {
		eg IGroup
		f  fc
		n  int
	}
	tests := []struct {
		name      string
		args      args
		want      int32
		wantError bool
	}{
		{name: "continue", args: args{eg: NewContinue(ctx), n: 10, f: func(ctx context.Context) error { atomic.AddInt32(&number, 1); return nil }}, want: 10},
		{name: "cancel", args: args{eg: NewCancel(ctx), n: 10, f: func(ctx context.Context) error { atomic.AddInt32(&number, 1); return nil }}, want: 20},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < tt.args.n; i++ {
				tt.args.eg.Go(tt.args.f)
			}
			gotErr := tt.args.eg.Wait()
			if (gotErr != nil) != tt.wantError {
				t.Errorf("TestErrGroup_Wait gotErr: %v, wantErr: %v", gotErr, tt.wantError)
			}
			if number != tt.want {
				t.Errorf("TestErrGroup_Wait got: %v, want: %v", number, tt.want)
			}
		})
	}
}

func TestErrGroup_Error(t *testing.T) {
	var number int32
	var ctx context.Context = context.Background()
	var lock sync.Mutex

	type args struct {
		eg IGroup
		f  fc
		n  int
	}
	tests := []struct {
		name      string
		args      args
		want      int32
		wantError bool
	}{
		{name: "cancel", args: args{eg: NewCancel(ctx), n: 10, f: func(ctx context.Context) error {
			lock.Lock()
			defer lock.Unlock()
			time.Sleep(10 * time.Millisecond)
			select {
			case <-ctx.Done():
				return nil
			default:
				if number++; number == 5 {
					return fmt.Errorf("limit 5")
				}
			}
			return nil
		}}, want: 5, wantError: true},
		{name: "continue", args: args{eg: NewContinue(ctx), n: 10, f: func(ctx context.Context) error {
			lock.Lock()
			defer lock.Unlock()
			select {
			case <-ctx.Done():
				return nil
			default:
				if number++; number == 5 {
					return fmt.Errorf("limit 5")
				}
			}
			return nil
		}}, want: 10, wantError: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			number = 0
			for i := 0; i < tt.args.n; i++ {
				tt.args.eg.Go(tt.args.f)
			}
			gotErr := tt.args.eg.Wait()
			if (gotErr != nil) != tt.wantError {
				t.Errorf("TestErrGroup_Wait gotErr: %v, wantErr: %v", gotErr, tt.wantError)
			}
			if number != tt.want {
				t.Errorf("TestErrGroup_Wait got: %v, want: %v", number, tt.want)
			}
		})
	}
}

func TestErrGroup_Panic(t *testing.T) {
	var number int32
	var ctx context.Context = context.Background()
	var lock sync.Mutex

	type args struct {
		eg IGroup
		f  fc
		n  int
	}
	tests := []struct {
		name      string
		args      args
		want      int32
		wantError bool
	}{
		{name: "cancel", args: args{eg: NewCancel(ctx), n: 10, f: func(ctx context.Context) error {
			lock.Lock()
			defer lock.Unlock()
			time.Sleep(100 * time.Millisecond)
			select {
			case <-ctx.Done():
				return nil
			default:
				if number++; number == 5 {
					panic("limit 5")
				}
			}
			return nil
		}}, want: 5, wantError: true},
		{name: "continue", args: args{eg: NewContinue(ctx), n: 10, f: func(ctx context.Context) error {
			lock.Lock()
			defer lock.Unlock()
			select {
			case <-ctx.Done():
				return nil
			default:
				if number++; number == 5 {
					panic("limit 5")
				}
			}
			return nil
		}}, want: 10, wantError: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			number = 0
			for i := 0; i < tt.args.n; i++ {
				tt.args.eg.Go(tt.args.f)
			}
			gotErr := tt.args.eg.Wait()
			//t.Log(gotErr)
			if (gotErr != nil) != tt.wantError {
				t.Errorf("TestErrGroup_Wait gotErr: %v, wantErr: %v", gotErr, tt.wantError)
			}
			if number != tt.want {
				t.Errorf("TestErrGroup_Wait got: %v, want: %v", number, tt.want)
			}
		})
	}
}

func TestErrGroup_WithMaxProcess(t *testing.T) {
	var number int32
	var ctx context.Context = context.Background()

	type args struct {
		eg IGroup
		f  fc
		n  int
	}
	tests := []struct {
		name      string
		args      args
		want      int32
		wantError bool
	}{
		{name: "cancel", args: args{eg: NewCancel(ctx, WithMaxProcess(1)), n: 10, f: func(ctx context.Context) error {
			select {
			case <-ctx.Done():
				return nil
			default:
				if number++; number == 5 {
					return fmt.Errorf("limit 5")
				}
			}
			return nil
		}}, want: 5, wantError: true},
		{name: "continue", args: args{eg: NewContinue(ctx, WithMaxProcess(1)), n: 10, f: func(ctx context.Context) error {
			select {
			case <-ctx.Done():
				return nil
			default:
				if number++; number == 5 {
					return fmt.Errorf("limit 5")
				}
			}
			return nil
		}}, want: 10, wantError: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			number = 0
			for i := 0; i < tt.args.n; i++ {
				tt.args.eg.Go(tt.args.f)
			}
			gotErr := tt.args.eg.Wait()
			if (gotErr != nil) != tt.wantError {
				t.Errorf("TestErrGroup_Wait gotErr: %v, wantErr: %v", gotErr, tt.wantError)
			}
			if number != tt.want {
				t.Errorf("TestErrGroup_Wait got: %v, want: %v", number, tt.want)
			}
		})
	}
}
