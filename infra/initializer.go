package infra

// 初始化接口

type Initializer interface {
	// 用于对象实例化后的初始化操作
	Init()
}

type InitializeRegister struct {
	Initializers []Initializer
}

// 初注册初始化对象
func (r *InitializeRegister) Register(initializer Initializer) {
	r.Initializers = append(r.Initializers, initializer)
}
