package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
	"fmt"
	"encoding/hex"
	"crypto/ecdsa"
	"crypto/rand"
	"strings"
	"crypto/elliptic"
	"math/big"
)

const  subsidy=1000  //奖励，矿工挖矿给予的奖励
//输入





//交易，编号，输入，输出
type Transaction struct {
	ID []byte
	Vin [] TXInput
	Vout [] TXOutput
}
//序列化，对象到二进制，
func (tx Transaction)Serialize()[]byte{
	var encoded  bytes.Buffer
	enc:=gob.NewEncoder(&encoded)//编码器
	err:=enc.Encode(tx)//编码
	if err!=nil{
		log.Panic(err)
	}
	return encoded.Bytes()//返回二进制数据
}
//反序列化，二进制到对象
func DeserializeTransaction(data []byte)Transaction{
	var transaction Transaction
	decoder:=gob.NewDecoder(bytes.NewReader(data))//解码器
	err:=decoder.Decode(&transaction)//解码
	if err!=nil{
		log.Panic(err)
	}
	return transaction
}


//对于交易事务进行哈希
func (tx *Transaction)Hash()[]byte{
	var hash[32]byte
	txCopy :=*tx
	txCopy.ID=[]byte{}
	hash=sha256.Sum256(txCopy.Serialize())//取得二进制进行哈希
	return hash[:]
}
//签名
func(tx *Transaction)Sign(privateKey ecdsa.PrivateKey,prevTXs map[string]Transaction){
	if tx.IsCoinBase(){
		return //如果挖矿返回。无需签名
	}
	for _,vin :=range tx.Vin{
		if prevTXs[hex.EncodeToString(vin.Txid)].ID == nil{
			//log.Panic("以前的交易不正确")
		}
	}
	txCopy:=tx.TrimmedCopy()//拷贝没有私钥等等的副本
	for inID,vin:=range txCopy.Vin{
		//设置签名为空与公钥
		prevTx:=prevTXs[hex.EncodeToString(vin.Txid)]
		txCopy.Vin[inID].Signature=nil
		txCopy.Vin[inID].PubKey=prevTx.Vout[vin.Vout].PubKeyHash
		txCopy.ID=txCopy.Hash()
		txCopy.Vin[inID].PubKey=nil

		//datatoSign:=fmt.Sprintf("%x\n",txCopy)//要签名的数据

		r,s,err:=ecdsa.Sign(rand.Reader,&privateKey,txCopy.ID)
		if err != nil{
			log.Panic(err)
		}
		signature:=append(r.Bytes(),s.Bytes()...)
		tx.Vin[inID].Signature=signature


	}


}
//用于签名的交易事务，裁剪的副本
func (tx *Transaction)TrimmedCopy()Transaction{
	var inputs []TXInput
	var outputs []TXOutput
	for _,vin:=range tx.Vin{
		inputs=append(inputs,TXInput{vin.Txid,vin.Vout,nil,nil})
	}
	for _,vout:=range tx.Vout{
		outputs=append(outputs,TXOutput{vout.Value,vout.PubKeyHash})

	}
	txCopy:=Transaction{tx.ID,inputs,outputs}
	return txCopy

}
//把对象作为字符串展示出来
func  (tx Transaction) String()string{
	var lines []string
	lines=append(lines,fmt.Sprintf("Transaction %x",tx.ID))
	for i,input :=range tx.Vin{
		lines=append(lines,fmt.Sprintf("input %d",i))
		lines=append(lines,fmt.Sprintf("TXID %x",input.Txid))
		lines=append(lines,fmt.Sprintf("OUT %d",input.Vout))
		lines=append(lines,fmt.Sprintf("Signature %x",input.Signature))
		lines=append(lines,fmt.Sprintf("Pubkey %x",input.PubKey))
	}
	for i,output :=range tx.Vout{
		lines=append(lines,fmt.Sprintf("out %d",i))
		lines=append(lines,fmt.Sprintf("value %d",output.Value))
		lines=append(lines,fmt.Sprintf("out %x",output.PubKeyHash))
	}
	return strings.Join(lines,"\n")
}
//签名认证
func  (tx *Transaction)Verify(prevTXs map[string]Transaction)bool{
	if tx.IsCoinBase(){
		return true //如果挖矿返回。无需签名
	}
	for _,vin :=range tx.Vin{
		if prevTXs[hex.EncodeToString(vin.Txid)].ID==nil{
			log.Panic("之前交易是错误的")
		}
	}
	txCopy:=tx.TrimmedCopy()//拷贝
	curve:=elliptic.P256()//加密算法
	for inID,vin :=range tx.Vin{
		prevTx:=prevTXs[hex.EncodeToString(vin.Txid)]//前缀
		txCopy.Vin[inID].Signature=nil
		txCopy.Vin[inID].PubKey=prevTx.Vout[vin.Vout].PubKeyHash//设置公钥

		r:=big.Int{}
		s:=big.Int{}
		siglen:=len(vin.Signature)//统计签名长度
		r.SetBytes(vin.Signature[:(siglen/2)])
		s.SetBytes(vin.Signature[(siglen/2):])

		x:=big.Int{}
		y:=big.Int{}
		keylen:=len(vin.PubKey)
		x.SetBytes(vin.PubKey[:(keylen/2)])
		y.SetBytes(vin.PubKey[(keylen/2):])
		datatoVerify:=fmt.Sprintf("%x\n",txCopy)//校验
		rawPubkey:=ecdsa.PublicKey{curve,&x,&y}
		if ecdsa.Verify(&rawPubkey,[]byte(datatoVerify),&r,&s)==false{
			return false
		}
		txCopy.Vin[inID].PubKey=nil
	}


	return true

}



//检查交易事务是否为coinbase，挖矿得来的奖励币
func( tx *Transaction) IsCoinBase()bool{
	return len(tx.Vin)==1 && len(tx.Vin[0].Txid)==0 && tx.Vin[0].Vout==-1
}
//设置交易ID，从二进制数据中
func (tx *Transaction)SetID(){
	var  encoded bytes.Buffer //开辟内存
	var hash[32] byte  //哈希数组
	enc:=gob.NewEncoder(&encoded)//解码对象
	err:=enc.Encode(tx)//解码
	if err!=nil{
		log.Panic(err)
	}
	hash=sha256.Sum256(encoded.Bytes())//计算哈希
	tx.ID=hash[:]//设置ID
}
//挖矿交易
func NewCoinBaseTX(to ,data string)*Transaction{
	if data==""{
		data=fmt.Sprintf(" 奖励给 %s",to)
	}
	txin:=TXInput{[]byte{},-1,nil,[]byte(data)}//输入奖励
	txout:=NewTXOUTput(subsidy,to)
	tx:=Transaction{nil,[]TXInput{txin},[]TXOutput{*txout}}//交易
	return &tx
}



