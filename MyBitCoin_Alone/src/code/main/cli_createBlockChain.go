package main

import "fmt"
import "log"

func (cli *CLI) createBlockChain(address string){
	if !ValidateAddress(address){
		log.Panic("地址错误")
	}
	bc:=CreateBlockChain(address)//创建一个区块链
	bc.db.Close()
	fmt.Println("创建成功")
}
