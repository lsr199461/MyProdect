package main

import "fmt"
import "log"


func (cli *CLI)getBalance(address string,nodeID string){
	if !ValidateAddress(address){
		log.Panic("地址错误")
	}
	bc:=NewBlockChain(nodeID)//根据地址创建
	UTXOSet:=UTXOSet{bc}//创建UTXO
	defer  bc.db.Close()//延迟关闭数据库
	balance:=0
	pubkeyhash:=Base58Decode([]byte(address))//提取公钥
	pubkeyhash=pubkeyhash[1:len(pubkeyhash)-4]
	UTXOs:=UTXOSet.FindUTXO(pubkeyhash)//根据公钥查询

	for  _,out :=range UTXOs {
		balance += out.Value //取出金额
	}
	fmt.Printf("查询的金额如下%s :%d \n",address,balance)


}