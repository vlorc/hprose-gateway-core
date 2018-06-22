package option

func NewDefault() *GatewayOption {
	return &GatewayOption{Prefix: "*"}
}

func NewOption(o ...func(*GatewayOption)) *GatewayOption {
	return NewOptionWith(NewDefault(), o...)
}

func NewOptionWith(opt *GatewayOption, o ...func(*GatewayOption)) *GatewayOption {
	for _, v := range o {
		v(opt)
	}
	return opt
}
