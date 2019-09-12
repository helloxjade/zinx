package iface

type IMsgHandle interface {
	AddRouter(uint32, IRouter)
	DoMsgHandler(IRequest)
}
