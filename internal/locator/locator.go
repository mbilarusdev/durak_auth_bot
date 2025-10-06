package locator

var Instance Locator

type Locator interface {
	Register(name string, impl any)
	Get(name string) any
}

type AuthBotLocator struct {
	registerMap map[string]any
}

func Setup() {
	Instance = new(AuthBotLocator)
	Instance.(*AuthBotLocator).registerMap = make(map[string]any)
}

func (locator *AuthBotLocator) Register(name string, impl any) {
	locator.registerMap[name] = impl
}

func (locator *AuthBotLocator) Get(name string) any {
	return locator.registerMap[name]
}
