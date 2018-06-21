package main

import (
	"github.com/boltdb/bolt"
	"log"
	"fmt"
	"encoding/hex"
	"os"
	"crypto/ecdsa"
	"bytes"
	"errors"
	//"github.com/hyperledger/fabric/common/genesis"
	//"github.com/golang/net/html/atom"
	//"github.com/golang/net/html/atom"
)

var dbFile="blockchain.db"  //数据库文件名当前目录下
const blockBucket="blocks"  //名称，
const  genesisCoinbaseData="区块链交流QQ1114747523"


type  BlockChain struct {
	tip []byte  //二进制数据
	db *bolt.DB //数据库
}

//挖矿带来的交易，
func (blockchain *BlockChain)MineBlock(transactions []*Transaction)*Block{
	var lastHash [] byte //最后的哈希
	var  lastHeight int //最后的长度


	for _,tx:=range transactions{
		if blockchain.VertifyTransaction(tx)!=true{
			log.Panic("交易不正确，有错误")
		}
	}


	err:=blockchain.db.View(func(tx *bolt.Tx)error {
		bucket:=tx.Bucket([]byte (blockBucket)) //查看数据
		lastHash=bucket.Get([]byte("1"))//取出最后区块的哈希

		blockData:=bucket.Get(lastHash)//取出最后的区块数据
		block:=DeserializeBlock(blockData)//解码

		lastHeight=block.Height//抓取宽度
		return nil
	})
	if err!=nil{
		log.Panic(err)//处理错误
	}
	newBlock:=NewBlock(transactions,lastHash,lastHeight+1)//创建一个新的区块
	err=blockchain.db.Update(func (tx *bolt.Tx)error{
		bucket:=tx.Bucket([]byte(blockBucket))//取出索引
		err:=bucket.Put(newBlock.Hash,newBlock.Serialize())//存入数据库
		if err!=nil{
			log.Panic(err)//处理错误
		}
		err=bucket.Put([]byte("1"),newBlock.Hash)//压入保存最后一个哈希
		if err!=nil{
			log.Panic(err)//处理错误
		}
		blockchain.tip=newBlock.Hash //保存上一块的哈希
		return nil
	})

	return newBlock
}


//获取没使用输出的交易列表
func (blockchain *BlockChain)FindUnspentTransactions(pubkeyhash []byte)[]Transaction{
	var unspentTXs [] Transaction //交易事务
	spentTXOS:=make(map[string][]int)//开辟内存
	bci:=blockchain.Iterator() //迭代器
	for{
		block:=bci.next()//循环下一个
		for _,tx :=range block.Transactions{//循环每个交易
			txID:=hex.EncodeToString(tx.ID)//获取交易编号

		Outputs:
			for outindex,out:=range tx.Vout{//循环遍历输出
				if spentTXOS[txID]!=nil{
					for _,spentOut:=range spentTXOS[txID]{
						if spentOut==outindex{
							continue Outputs  //循环到不等
						}
					}
				}
				if out.IsLockedWithKey(pubkeyhash){
					unspentTXs=append(unspentTXs,*tx)//加入列表
				}
			}
			if tx.IsCoinBase()==false{
				for _,in  :=range tx.Vin{
					if in.UsesKey(pubkeyhash){//判断是否可以锁定
						inTxID:=hex.EncodeToString(in.Txid)//编码为字符串
						spentTXOS[inTxID]=append(spentTXOS[inTxID],in.Vout)
					}
				}
			}
		}
		if len(block.PrevBlockHash)==0{//最后一块，跳出
			break
		}
	}
	return  unspentTXs
}


//获取所有没有使用的交易
func  (blockchain * BlockChain)FindUTXO()map[string]TXoutputs{
	UTXO:=make(map[string]TXoutputs)//新建序列
	spentTXOs:=make(map[string][]int)//花掉的交易
	bci:=blockchain.Iterator()//迭代器
	for{
		block :=bci.next()

		for _,tx :=range block.Transactions{
			txID:=hex.EncodeToString(tx.ID )//根据编号编码

		Outputs:
			for  outIdx,out:=range tx.Vout{
				if spentTXOs[txID]!=nil{
					for _,spendoutidx :=range spentTXOs[txID]{
						if spendoutidx==outIdx{
							continue Outputs
						}
					}
				}
				outs:=UTXO[txID]
				outs.Outputs=append(outs.Outputs,out)//叠加
				UTXO[txID]=outs//抓取编号赋值保存
			}
			if tx.IsCoinBase()==false{
				for _,in :=range  tx.Vin{
					inTxID:=hex.EncodeToString(in.Txid)//编码
					spentTXOs[inTxID]=append(spentTXOs[inTxID],in.Vout)//追加
				}
			}

		}



		if len(block.PrevBlockHash)==0{
			break
		}
	}


	return UTXO
}

//获取没有使用的输出以参考输入
func  (blockchain * BlockChain)FindSpendableOutputs(pubkeyhash[]byte,amount int)(int,map[string][]int){
	unspentOutputs:=make(map[string][]int) //输出
	unspentTxs:=blockchain.FindUnspentTransactions(pubkeyhash)//根据地址查找所有交易
	accmulated:=0//累计
Work:
	for  _,tx  :=range unspentTxs{
		txID:=hex.EncodeToString(tx.ID)//获取编号
		for outindex,out:=range tx.Vout{
			if out.IsLockedWithKey(pubkeyhash)&&accmulated<amount{
				accmulated+=out.Value //统计金额
			    unspentOutputs[txID]=append(unspentOutputs[txID],outindex)//序列叠加
			    if accmulated>=amount{
					break Work
				}
			}
		}
	}
	return  accmulated,unspentOutputs
}



//迭代器
func (block *BlockChain)Iterator()*BlockChainIterator{
	bcit:=&BlockChainIterator{block.tip,block.db}
	return bcit //根据区块链创建区块链迭代器
}



//判断数据库是否存在
func dbExists(dbFile string)bool{
	if _,err:=os.Stat(dbFile);os.IsNotExist(err){
		return false
	}
	return true
}


//新建一个区块链
func NewBlockChain(nodeID string)*BlockChain{
	dbFile=fmt.Sprintf("blockchain_%s.db",nodeID)
	if dbExists(dbFile)==false{
		fmt.Println("数据库不存在，创建一个先")
		os.Exit(1)
	}
	var tip []byte  //存储区块链的二进制数据
	db,err:=bolt.Open(dbFile,0600,nil)//打开数据库
	if err!=nil{
		log.Panic(err)//处理数据库打开错误
	}
	//处理数据更新
	err=db.Update(func (tx *bolt.Tx)error {
		bucket:=tx.Bucket([]byte(blockBucket))//按照名称打开数据库的表格
		tip=bucket.Get([]byte("1"))
		return nil
	})
	if err!=nil{
		log.Panic(err)//处理数据库更新错误
	}
	bc:=BlockChain{tip,db} //创建一个区块链
	return &bc
}
func  CreateBlockChain(address string,nodeID string)*BlockChain{
	dbFile:=fmt.Sprintf("blockchain_%s.db",nodeID)
	if  dbExists(dbFile){
		fmt.Println("数据库已经存在无需创建")
		os.Exit(1)
	}

	var tip []byte  //存储区块链的二进制数据


	cbtx:=NewCoinBaseTX(address,genesisCoinbaseData)//创建创世区块的事无交易
	genesis :=NewGenesisBlock(cbtx)//创建创世区块的快

	db,err:=bolt.Open(dbFile,0600,nil)//打开数据库
	if err!=nil{
		log.Panic(err)//处理数据库打开错误
	}
	err=db.Update( func (tx *bolt.Tx)error{

		bucket,err:=tx.CreateBucket([]byte(blockBucket))
		if err!=nil{
			log.Panic(err)//处理数据库打开错误
		}
		err=bucket.Put(genesis.Hash,genesis.Serialize())//存储
		if err!=nil{
			log.Panic(err)
		}
		err=bucket.Put([]byte("1"),genesis.Hash)
		if err!=nil{
			log.Panic(err)
		}
		tip=genesis.Hash
		return nil
	})
	if err!=nil{
		log.Panic(err)
	}

	bc:=BlockChain{tip,db} //创建一个区块链
	return &bc

}
//交易签名
func (blockchain *BlockChain)SignTransaction(tx *Transaction,privatekey ecdsa.PrivateKey){
	prevTXs:=make(map[string]Transaction)
	for _,vin :=range tx.Vin{
		preTx,err:=blockchain.FindTransaction(vin.Txid)
		if err!=nil{
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(preTx.ID)]=preTx
	}
	tx.Sign(privatekey,prevTXs)
}

func (blockchain *BlockChain)FindTransaction(ID []byte)(Transaction,error){
	bci:=blockchain.Iterator()
	for{
		block:=bci.next()
		for _,tx :=range block.Transactions{
			if bytes.Compare(tx.ID,ID)==0{
				return *tx,nil
			}
		}
		if len(block.PrevBlockHash)==0{
			break
		}
	}
	return Transaction{},nil
}
func (blockchain *BlockChain)VertifyTransaction(tx *Transaction)bool{
	prevTxs:=make(map[string]Transaction)
	for _,vin:=range tx.Vin{
		prevTx,err:=blockchain.FindTransaction(vin.Txid)//查找交易
		if err!=nil{
			log.Panic(err)
		}
		prevTxs[hex.EncodeToString(prevTx.ID)]=prevTx
	}
	return tx.Verify(prevTxs)
}

//抓取最后一个区块用于同步
func (blockchain *BlockChain)GetBestHeight()int {
	var lastBlock Block//最后一个区块
	err :=blockchain.db.View(func(tx *bolt.Tx) error {
		bucket :=tx.Bucket([]byte (blockBucket))//取出数据库的数据对象
		lastHash:=bucket.Get([]byte("1"))//取得最后的哈希
		blockdata:=bucket.Get(lastHash)//取得上一个哈希
		lastBlock=*DeserializeBlock(blockdata)//解码区块数据
		return nil
	})
	if err!=nil{
		log.Panic(err)
	}
	return lastBlock.Height
}
//增加模块
func (blockchain*BlockChain)AddBlock(block *Block){
	err:=blockchain.db.Update(func(tx *bolt.Tx) error {
		bucket :=tx.Bucket([]byte (blockBucket))//抓取区块索引
		blockInDb:=bucket.Get(block.Hash)//判断区块是否存在
		if blockInDb!=nil{
			return nil
		}
		blockData :=block.Serialize()//序列化
		err:=bucket.Put(block.Hash,blockData)//压入数据
		if err!=nil{
			log.Panic(err)
		}

		lastHash:=bucket.Get([]byte("1"))//取出数据
		lastBlockdata:=bucket.Get(lastHash)//取得最后一个区块
		lastBlock:=DeserializeBlock(lastBlockdata)//反序列化上一个区块

		if block.Height >lastBlock.Height{//判断区块链的宽度。
			err=bucket.Put( []byte("1"),block.Hash)//压入哈希
			if err!=nil{
				log.Panic(err)
			}
			blockchain.tip=block.Hash

		}




		return nil
	})
	if err!=nil{
		log.Panic(err)
	}

}


func (blockchain * BlockChain)GetBlockHashes()[][]byte{
	var  blocks [][] byte
	bci :=blockchain.Iterator()
	for {
		block:=bci.next()
		blocks =append(blocks,block.Hash)//查找过程
		if len(block.PrevBlockHash)==0{
			break
		}
	}
	return blocks
}

//区块连中查找区块
func (blockchain * BlockChain)GetBlock(blockhash []byte)(Block,error){

	var bc Block
	err:=blockchain.db.View(func(tx *bolt.Tx) error {
		bucket :=tx.Bucket([]byte(blockBucket))//取出数据库的数据对象
		blockdata:=bucket.Get(blockhash)//取出数据
		if blockdata==nil{
			return errors.New("没有找到区块")
		}
		bc=*DeserializeBlock(blockdata)
		return nil
	})

	if err!=nil{
		return bc,err
	}
	return bc,nil
}


