package errgroup

type IOption func(conf *config)

func WithMaxProcess(n int) IOption {
	return func(conf *config) { conf.maxProcess = n }
}

func WithIgnoreErr(f func(err error) bool) IOption {
	return func(conf *config) { conf.ignoreError = f }
}
