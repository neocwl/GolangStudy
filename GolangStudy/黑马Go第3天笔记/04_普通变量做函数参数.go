package main

import(
	"fmt"
)

func main() {
	a,b:=10,20
	fmt.Printf("a=%d,b=%d\n",a,b)
	fmt.Printf("&a=%v,&b=%v\n",&a,&b)
	//通过一个函数交换a和b的内容
	swap(a,b)//变量本身传递，值传递（站在变量角度）
	fmt.Printf("a=%d,b=%d\n",a,b)
}
func swap(a, b int) {
	a,b=b,a
	fmt.Printf("swap: a=%d,b=%d\n",a,b)
	fmt.Printf("&a=%v,&b=%v\n",&a,&b)

}
