package errgroup

type IOption func(IGroup)

func WithMaxProcess(n int) IOption {
	return func(g IGroup) { g.maxProcess(n) }
}
