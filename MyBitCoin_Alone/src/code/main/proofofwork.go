package main

import (
	"math"
	"math/big"
	"bytes"
	"fmt"
	"crypto/sha256"

	"code/tools"
)

var (
	maxNonce=math.MaxInt64 //最大的64位整数
)
const  targetBits=18 //对比的位数.位数越高，难度越大，时间越长

type ProofOfWork struct {
	block *Block //区块
	target * big.Int  //存储计算哈希对比的特定整数
}

//创建一个工作量证明的挖矿对象
func NewProofOfWork(block *Block)*ProofOfWork{
	target:=big.NewInt(1)//初始化目标整数
	target.Lsh(target,uint(256-targetBits))//数据转换
	pow:=&ProofOfWork{block,target}//创建对象
	return pow
}
//准备数据进行挖矿计算
func (pow * ProofOfWork ) prepareData(nonce int)[]byte{
	data:=bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,//上一块哈希
			pow.block.HashTransactions(),//当前数据
			tools.IntToHex(pow.block.Timestamp),//时间十六进制
			tools.IntToHex(int64(targetBits)),//位数，十六进制
			tools.IntToHex(int64(nonce)),//保存工作量的nonce
		},[]byte{},
	)
	return data
}
//挖矿执行
func  (pow * ProofOfWork ) Run()(int,[]byte){
	var  hashInt big.Int
	var hash [32]byte
	nonce :=0
	//fmt.Printf("当前挖矿计算的区块数据%s",pow.block.Data)
	for nonce<maxNonce{
		data :=pow.prepareData(nonce)//准备好的数据
		hash=sha256.Sum256(data) //计算出哈希
		fmt.Printf("\r%x",hash)//打印显示哈希
		hashInt.SetBytes(hash[:])//获取要对比的数据
		if hashInt.Cmp(pow.target)==-1{ //挖矿的校验
			break
		}else{
			nonce++
		}



	}
	return nonce,hash[:]//nonce解题的答案，hash当前哈希

}
//校验挖矿是不是真的成功
func (pow * ProofOfWork ) Validate()bool{
	var  hashInt big.Int
	data :=pow.prepareData(pow.block.Nonce)//准备好的数据
	hash:=sha256.Sum256(data) //计算出哈希
	hashInt.SetBytes(hash[:])//获取要对比的数据
	isValid:= hashInt.Cmp(pow.target)==-1//校验数据
	return isValid


}
