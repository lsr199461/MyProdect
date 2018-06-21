package main

import "fmt"
import "log"

func (cli *CLI) createBlockChain(address string,nodeID string){
	if !ValidateAddress(address){
		log.Panic("地址错误")
	}
	bc:=CreateBlockChain(address,nodeID)//创建一个区块链
	defer bc.db.Close()

	UTXOSet:=UTXOSet{bc}//创建UTXO集合
	UTXOSet.Reindex()
	fmt.Println("创建成功")
}
