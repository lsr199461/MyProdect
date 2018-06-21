package main

import (
	"bytes"
	"encoding/gob"
	"log"
)
//输出
type TXOutput struct {
	Value  int  //，output保存了“币”（上面的Value）
	PubKeyHash  []byte //用脚本语言意味着比特币可以也作为智能合约平台,公钥
}

//输出锁住的标志
func (out *TXOutput)Lock(address []byte){
	pubkeyhash:=Base58Decode(address)//编码
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
type  TXoutputs struct{
	Outputs [] TXOutput
}


//对象到二进制
func  (outs * TXoutputs)Serialize()[]byte{
	var buff  bytes.Buffer //开辟内存
	enc:=gob.NewEncoder(&buff)//创建编码器
	err:=enc.Encode(outs)
	if err!=nil{
		log.Panic(err)//处理错误
	}
	return buff.Bytes()//返回二进制
}
////二进制到对象
func DeserializeOutputs(data[]byte)TXoutputs{
	var outputs TXoutputs
	dec :=gob.NewDecoder(bytes.NewReader(data))//解码对象
	err:=dec.Decode(&outputs)//解码操作
	if err!=nil{
		log.Panic(err)//处理错误
	}
	return outputs
}













