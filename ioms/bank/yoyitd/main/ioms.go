package main

import "tradebank/ioms/bank/yoyitd"

func main() {
	m := yoyitd.NewServer()

	m.InitServer()
	m.Start()
}
