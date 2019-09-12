package iface

type IServer interface {
	Start()
	Stop()
	Serve()
	AddRouter(uint32, IRouter)
	GetConnMgr() IconnManager
	RegistStartHookFunc(func(IConnection))
	RegistStopHookFunc(func(IConnection))
	CallStartHook(IConnection)
	CallStopHookFunc(IConnection)
}
