package infra

var apiInitializerRegister *InitializeRegister = new(InitializeRegister)

// 注册web api初始化对象
func RegisterApi(api Initializer) {
	apiInitializerRegister.Register(api)
}

// 获取注册过的 web api 初始对象
func GetApiInitializers() []Initializer {
	return apiInitializerRegister.Initializers
}

type WebApiStarter struct {
	BaseStarter
}

func (w *WebApiStarter) Setup(ctx StarterContext) {
	for _, initializer := range GetApiInitializers() {
		initializer.Init()
	}
}
