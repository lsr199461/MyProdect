package main

import (
	"bytes"
	"code/tools"
)
//输出
type TXOutput struct {
	Value  int  //，output保存了“币”（上面的Value）
	PubKeyHash  []byte //用脚本语言意味着比特币可以也作为智能合约平台,公钥
}

//输出锁住的标志
func (out *TXOutput)Lock(address []byte){
	pubkeyhash:=tools.Base58Decode(address)//编码
	pubkeyhash=pubkeyhash[1:len(pubkeyhash)-4]//截取有效哈希
	out.PubKeyHash=pubkeyhash//锁住，无法再被修改
}
//监测是否被key锁住
func (out *TXOutput)IsLockedWithKey(pubkeyHAsh []byte)bool{
	return  bytes.Compare(out.PubKeyHash,pubkeyHAsh)==0
}
//创造一个输出
func NewTXOUTput(value int,address string)*TXOutput{
	txo:=&TXOutput{value,nil}//输出
	txo.Lock([]byte(address))//锁住
	return txo
}