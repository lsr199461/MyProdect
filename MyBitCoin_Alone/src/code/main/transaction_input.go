package main

import "bytes"

type TXInput struct {
	Txid []byte  //Txid存储了交易的id，
	Vout  int  //Vout则保存该交易的中一个output索引
    Signature []byte  //签名
    PubKey []byte  //公钥
}
//key监测一下地址与交易
func (in *TXInput)UsesKey(pubKeyHash []byte) bool{
	lockinghash:=HashPubkey(in.PubKey)
	return bytes.Compare(lockinghash,pubKeyHash)==0
}



