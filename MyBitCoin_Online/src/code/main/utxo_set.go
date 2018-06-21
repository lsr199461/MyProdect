package main
import (
	"log"//日志
	"encoding/hex"//十六进制
	"github.com/boltdb/bolt"//使用数据库
	//"github.com/astaxie/beego"
)


const utxoBucket="chainstate"//存储状态


//二次封装区块链
type UTXOSet struct {
	blockchain *BlockChain
}
//输出查找并返回未曾使用的输出
func (utxo UTXOSet)FindSpendableOutpus(publickeyhash []byte,amount int)(int,map[string][]int){
	unspentOutputs:=make(map[string][]int)//处理输出
	accumulated:=0 //累计的金额
	db:=utxo.blockchain.db //调用数据库

	//查询数据
	err:=db.View(func (tx *bolt.Tx)error{
		bucket:=tx.Bucket([]byte(utxoBucket))//查询数据
		cur:=bucket.Cursor()//当前的游标
		for key,value :=cur.First();key!=nil;key,value=cur.Next(){//循环
			txID:=hex.EncodeToString(key)//编号
			outs:=DeserializeOutputs(value)//解码
			for outIdx,out  :=range  outs.Outputs{
				//判断是否锁住，判断金额
				if out.IsLockedWithKey(publickeyhash) && accumulated<amount{
					accumulated+=out.Value//叠加金额
					unspentOutputs[txID]=append(unspentOutputs[txID],outIdx)//叠加序列
				}
			}
		}

		return nil

	})
	if err!=nil{
		log.Panic(err)//输出错误
	}
	return accumulated,	unspentOutputs //返回数据
}

//查找UTXO，按照公钥查询
func  (utxo  UTXOSet)FindUTXO(publickeyHash []byte)[]TXOutput{
	var UTXOs  []TXOutput
	db:=utxo.blockchain.db //取出数据库，进行查询
	err:=db.View(func (tx *bolt.Tx)error{
		bucket:=tx.Bucket([]byte(utxoBucket))//查询数据
		cur:=bucket.Cursor()//当前的游标
		for key,value :=cur.First();key!=nil;key,value=cur.Next() { //循环
			outs:=DeserializeOutputs(value)//反序列化数据库的数据
			for _,out:=range outs.Outputs{
				if out.IsLockedWithKey(publickeyHash){//判断是否锁住
					UTXOs=append(UTXOs,out) //数据叠加
				}
			}

		}


		return nil
	})
	if err!=nil{
		log.Panic(err)//输出错误
	}
	return UTXOs

}
//统计交易
func (utxo UTXOSet)CountTransactions()int{
	db:=utxo.blockchain.db//引用数据库
	counter:=0
	err:=db.View(func(tx *bolt.Tx) error {
		bucket:=tx.Bucket([]byte(utxoBucket))//查询数据
		cur:=bucket.Cursor()//当前的游标
		for  k,_:=cur.First();k!=nil;k,_=cur.Next(){
			counter++ //叠加
		}

		return nil
	})

	if err!=nil{
		log.Panic(err)//输出错误
	}
	return counter
}
//重建索引
func (utxo UTXOSet)Reindex(){
	db:=utxo.blockchain.db//数据库
	buckername:=[]byte(utxoBucket)//数据
	err:= db.Update(func(tx *bolt.Tx) error {
		err:=tx.DeleteBucket(buckername)//删除
		if err!=nil && err!= bolt.ErrBucketNotFound{
			log.Panic(err)
		}
		_,err=tx.CreateBucket(buckername)//新建
		if err!=nil{
			log.Panic(err)
		}
		return  nil
	})
	if err!=nil{
		log.Panic(err)
	}
	UTXO:=utxo.blockchain.FindUTXO()//数据查找
	err=db.Update(func(tx *bolt.Tx) error {
		bucket :=tx.Bucket(buckername)//取出数据

		for txID,outs:=range UTXO{
			key,err:=hex.DecodeString(txID)
			if err!=nil{
				log.Panic(err)
			}
			err=bucket.Put(key,outs.Serialize())
			if err!=nil{
				log.Panic(err)
			}
		}
		return nil
	})


}
//刷新数据
func (utxo UTXOSet)Update(block *Block){
	db:=utxo.blockchain.db//取出数据库
	err:=db.Update(func(tx *bolt.Tx) error {
		bucket :=tx.Bucket([]byte (utxoBucket))//取出数据库的对象数据
		for _,tx :=range block.Transactions{//循环遍历所有的交易
			if tx.IsCoinBase()==false{//取出非挖矿
				for _,vin :=range tx.Vin{
					updateOuts:= TXoutputs{}//创建集合
					outsBytes:=bucket.Get(vin.Txid)//取出数据
					outs :=DeserializeOutputs(outsBytes)//解码二进制数据
					for outIdx ,out :=range  outs.Outputs{
						if outIdx!=vin.Vout{
							updateOuts.Outputs=append(updateOuts.Outputs,out)//序列叠加
						}

					}
					if len(updateOuts.Outputs)==0{
						err :=bucket.Delete(vin.Txid)//，删除
						if err!=nil{
							log.Panic(err)
						}
					} else{
						err:= bucket.Put(vin.Txid,updateOuts.Serialize())//处理错误
						if err!=nil{
							log.Panic(err)
						}
					}

				}


			}
			newOutputs:=TXoutputs{}
			for _,out :=range tx.Vout{
				newOutputs.Outputs=append(newOutputs.Outputs,out)//处理好了叠加
			}
			err:=bucket.Put(tx.ID,newOutputs.Serialize())
			if err!=nil{
				log.Panic(err)//输出错误
			}

		}


		return nil
	})
	if err!=nil{
		log.Panic(err)//输出错误
	}
}
