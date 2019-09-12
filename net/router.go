package net

import (
	"fmt"
	iface2 "github.com/helloxjade/zinx/iface"
)

//有的router不希望有prehandle或posthandle
type Router struct {
}

//绑定3个方法。处理业务做准备
func (r *Router) PreHandle(req iface2.IRequest) {
	fmt.Println("PreHandle called")
}

//真正的处理业务
func (r *Router) Handle(req iface2.IRequest) {
	fmt.Println("Handle called")
}

//处理业务之后做清理
func (r *Router) PostHandle(req iface2.IRequest) {
	fmt.Println("PostHandle called")
}
