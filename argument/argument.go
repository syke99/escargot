package argument

type Arguments struct {
	args map[string]*any
}

func NewArguments(key string, value any) *Arguments {
	args := Arguments{}

	argMap := make(map[string]*any)

	args.args = argMap

	return &args
}

func (a *Arguments) GetArgsLength() int {
	return len(a.args)
}

func (a *Arguments) GetArgsSlice() []*any {
	args := make([]*any, len(a.args))

	for _, v := range a.args {
		args = append(args, v)
	}

	return args
}

func (a *Arguments) GetArg(key string) any {
	arg := a.args[key]
	return &arg
}

func (a *Arguments) SetArg(key string, value any) {
	a.args[key] = &value
}
