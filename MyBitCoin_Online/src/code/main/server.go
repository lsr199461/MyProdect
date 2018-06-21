package main

import (
	"fmt"
	"bytes"
	"net"
	"encoding/gob"
	"log"
	"io"
	"encoding/hex"
	"io/ioutil"
)

const protocol = "tcp"   //安全保障的网络协议
const nodeVersion = 1    //版本
const commandlength = 12 //命令行长度

var nodeAddress string                     //节点地址
var miningAddress string                   //挖矿地址
var knowNodes = []string{"localhost:3000"} //已经知道的节点
var blocksInTransit = [][]byte{}
var mempool = make(map[string]Transaction) //内存池

type addr struct {
	Addrlist [] string //节点
}
type block struct {
	AddrFrom string //来源地址
	Block    []byte //块
}
type getblocks struct {
	//来源地址
	AddrFrom string
}
type getdata struct {
	AddrFrom string //来源
	Type     string //类型
	ID       [] byte
}
type inv struct {
	AddrFrom string //来源
	Type     string //l类型
	Items    [][]byte
}
type tx struct {
	AddFrom     string //来源
	Transaction []byte
}
type verzion struct {
	Version    int //版本参数
	BestHeight int
	AddrFrom   string
}

//字节到命令
func bytesToCommand(bytes []byte) string {
	var command []byte
	for _, b := range bytes {
		if b != 0x0 {
			command = append(command, b) //增加命令的字节

		}
	}
	return fmt.Sprintf("%s", command)
}

//命令到字节
func commandToBytes(command string) [] byte {
	var bytes [commandlength] byte
	for index, char := range command {
		bytes[index] = byte(char) //字节转化为索引
	}
	return bytes[:]
}

//提取命令
func extractCommand(request []byte) []byte {
	return request[:commandlength]
}

//请求块
func requestBlocks() {
	for _, node := range knowNodes { //给所有已经知道的节点发送请求
		sendGetBlocks(node)
	}
}

//发送块
func sendBlock(addr string, bc *Block) {
	data := block{nodeAddress, bc.Serialize()}             //构造模块
	payload := gobEncode(data)                             //追加的数据处理
	request := append(commandToBytes("block"), payload...) //定制请求
	sendData(addr, request)                                //发送数据
}

//发送地址
func sendaddr(address string) {
	nodes := addr{knowNodes}                              //已经知道的所有节点
	nodes.Addrlist = append(nodes.Addrlist, nodeAddress)  //追加当前节点
	payload := gobEncode(nodes)                           //增加解码的节点
	request := append(commandToBytes("addr"), payload...) //创建请求
	sendData(address, request)                            //发送数据
}

//发送数据
func sendData(addr string, data []byte) {
	conn, err := net.Dial(protocol, addr) //建立TCP网络连接对象
	defer conn.Close()                            //延迟关闭
	if err != nil {
		fmt.Printf("%s 地址不可到达\n", addr)
		for _, node := range knowNodes {
			if node != addr {
				knowNodes = append(knowNodes, addr) //刷新节点
			}
		}
	}
	_, err = io.Copy(conn, bytes.NewReader(data)) //拷贝数据，发送
	if err != nil {
		log.Panic(err) //处理错误
	}
}

//发送请求
func sendInv(address, kind string, items [][]byte) {
	inventory := inv{nodeAddress, kind, items}           //库存数据
	payload := gobEncode(inventory)                      //历史数据
	request := append(commandToBytes("inv"), payload...) //网络请求
	sendData(address, request)                           //发送数据

}

//发送请求多个模块
func sendGetBlocks(address string) {
	payload := gobEncode(getblocks{nodeAddress}) //解码地址
	request := append(commandToBytes("getblocks"), payload...)
	sendData(address, request) //发送数据与请求
}

//发送请求的数据
func sendGetData(address, kind string, id []byte) {
	payload := gobEncode(getdata{nodeAddress, kind, id}) //解码地址
	request := append(commandToBytes("getdata"), payload...)
	sendData(address, request) //发送数据与请求
}

//发送一个交易
func sendTx(addr string, tnx *Transaction) {
	data := tx{nodeAddress, tnx.Serialize()} //处理数据
	payload := gobEncode(data)               //编码
	request := append(commandToBytes("tx"), payload...)
	sendData(addr, request) //发送数据与请求
}

//发送版本信息
func sendVersion(addr string, bc *BlockChain) {
	bestHeight := bc.GetBestHeight() //最后一个区块的height
	payload := gobEncode(verzion{nodeVersion, bestHeight, nodeAddress})
	request := append(commandToBytes("version"), payload...)
	sendData(addr, request) //发送数据与请求
}

//模块的句柄
func handleBlock(request []byte, bc *BlockChain) {
	var buff bytes.Buffer               //二进制数据内存
	var payload block                   //地址
	buff.Write(request[commandlength:]) //取出数据
	dec := gob.NewDecoder(&buff)        //解码器
	err := dec.Decode(&payload)         //解码器
	if err != nil {
		log.Panic(err) //s数据处理
	}
	blockData := payload.Block           //区块的数据
	block := DeserializeBlock(blockData) //解码
	fmt.Printf("收到一个新的区块\n")
	bc.AddBlock(block)
	fmt.Printf("增加一个区块%x\n", block.Hash)
	if len(blocksInTransit) > 0 {
		blockhash := blocksInTransit[0]
		sendGetData(payload.AddrFrom, "block", blockhash) //发送请求
		blocksInTransit = blocksInTransit[1:]
	} else {
		UTXOSet := UTXOSet{bc}
		UTXOSet.Reindex() //重建索引
	}

}

//读取网络地址
func handleaddr(request []byte) {
	var buff bytes.Buffer //二进制数据内存
	var payload addr      //地址

	buff.Write(request[commandlength:]) //取出数据
	dec := gob.NewDecoder(&buff)        //解码器
	err := dec.Decode(&payload)         //解码器
	if err != nil {
		log.Panic(err) //s数据处理
	}
	knowNodes = append(knowNodes, payload.Addrlist...) //追加已知列表
	fmt.Printf("已经有了%d个节点", len(knowNodes))
	requestBlocks() //请求区块数据
}

//请求的版本
func handleInv(request []byte, bc *BlockChain) {
	var buff bytes.Buffer //二进制数据内存
	var payload inv       //地址inv

	buff.Write(request[commandlength:]) //取出数据
	dec := gob.NewDecoder(&buff)        //解码器
	err := dec.Decode(&payload)         //解码器
	if err != nil {
		log.Panic(err) //s数据处理
	}
	fmt.Printf("收到库存 %d %s ", len(payload.Items), payload.Type)

	if payload.Type == "block" {
		blocksInTransit = payload.Items                   //历史抓取的区块
		blockhash := payload.Items[0]                     //区块哈希
		sendGetData(payload.AddrFrom, "block", blockhash) //发送请求数据

		newInTransit := [][]byte{} //字节二维数组
		for _, b := range blocksInTransit {
			if bytes.Compare(b, blockhash) != 0 {
				newInTransit = append(newInTransit, b) //加入新的区块的哈希
			}
		}
		blocksInTransit = newInTransit //同步区块

	}

	if payload.Type == "tx" {
		txID := payload.Items[0] //编号
		if mempool[hex.EncodeToString(txID)].ID == nil {
			sendGetData(payload.AddrFrom, "tx", txID) //发起请求的交易
		}
	}

}

//抓取多个区块
func handleGetBlocks(request []byte, bc *BlockChain) {

	var buff bytes.Buffer //处理进制
	var payload getblocks //获取区块
	buff.Write(request[commandlength:])
	dec := gob.NewDecoder(&buff) //解码器
	err := dec.Decode(&payload)  //解码器
	if err != nil {
		log.Panic(err) //s数据处理
	}
	blocks := bc.GetBlockHashes()
	sendInv(payload.AddrFrom, "block", blocks)

}

//抓取数据
func handleGetData(request []byte, bc *BlockChain) {
	var buff bytes.Buffer //处理进制
	var payload getdata   //获取区块
	buff.Write(request[commandlength:])
	dec := gob.NewDecoder(&buff) //解码器
	err := dec.Decode(&payload)  //解码器
	if err != nil {
		log.Panic(err) //s数据处理
	}
	if payload.Type == "block" {
		block, err := bc.GetBlock([]byte(payload.ID )) //抓取一个区块
		if err != nil {
			return
		}
		sendBlock(payload.AddrFrom, &block) //发送区块
	}
	if payload.Type == "tx" {
		txID := hex.EncodeToString(payload.ID) //处理交易
		tx := mempool[txID]                    //内存池
		sendTx(payload.AddrFrom, &tx)          //发送交易
	}
}

//抓取交易
func handleTx(request []byte, bc *BlockChain) {
	var buff bytes.Buffer
	var payload tx

	buff.Write(request[commandlength:])
	dec := gob.NewDecoder(&buff) //解码器
	err := dec.Decode(&payload)  //解码器
	if err != nil {
		log.Panic(err) //s数据处理
	}

	txData := payload.Transaction           //交易数据
	tx := DeserializeTransaction(txData)    //解码交易数据
	mempool[hex.EncodeToString(tx.ID)] = tx //处理交易

	//if nodeAddress == knowNodes[0] {
	//	for _, node := range knowNodes {
	//		if node != nodeAddress && node != payload.AddFrom {
	//			sendInv(node, "tx", [][]byte{tx.ID}) //发送库存
	//		}
	//	}
	//}
	if len(mempool) >= 1 && len(miningAddress) > 0 {
		fmt.Printf("收到交易自行挖矿")
	MineTransactions:
		var txs []*Transaction //交易列表
		for id := range mempool {
			tx := mempool[id] //取得交易
			if bc.VertifyTransaction(&tx) { //校验交易是不是伪造
				txs = append(txs, &tx) //追加交易列表
			}
		}
		if len(txs) == 0 {
			fmt.Println("没有任何交易，等待新的交易加入")
			return
		}
		cbTx := NewCoinBaseTX(miningAddress, "") //创建一个地址，为这个地址挖矿
		txs = append(txs, cbTx)                  //叠加效应

		newBlock := bc.MineBlock(txs) //挖矿
		UTXOSet := UTXOSet{bc}
		UTXOSet.Reindex() //重建索引
		fmt.Printf("新的区块已经挖掘到\n")
		for _, tx := range txs {
			txID := hex.EncodeToString(tx.ID) //交易编号
			delete(mempool, txID)             //删除内存池
		}
		for _, node := range knowNodes {
			if node != nodeAddress {
				sendInv(node, "block", [][]byte{newBlock.Hash}) //挖矿成功广播
			}
		}
		if len(mempool) > 0 {
			goto MineTransactions
		}

	}

}

//处理版本
func handleVersion(request []byte, bc *BlockChain) {
	var buff bytes.Buffer
	var payload verzion
	buff.Write(request[commandlength:])
	dec := gob.NewDecoder(&buff) //解码器
	err := dec.Decode(&payload)  //解码器
	if err != nil {
		log.Panic(err) //s数据处理
	}
	mybestHeight := bc.GetBestHeight()        //抓取最好的宽度
	foreignerBestHeight := payload.BestHeight //抓取最好的宽度

	//版本同步
	if mybestHeight < foreignerBestHeight {
		sendGetBlocks(payload.AddrFrom)
	} else if mybestHeight > foreignerBestHeight {
		sendVersion(payload.AddrFrom, bc)
	}

	if !nodeIsKnow(payload.AddrFrom) {
		knowNodes = append(knowNodes, payload.AddrFrom) //判断节点是否已经知道
	}

}

//处理网络链接
func handleConnection(conn net.Conn, bc *BlockChain) {
	request, err := ioutil.ReadAll(conn) //处理所有网络连接
	if err != nil {
		log.Panic(err)
	}
	command := bytesToCommand(request[:commandlength])
	fmt.Printf("收到命令%s\n", command)

	switch command {
	case "addr":
		handleaddr(request)
	case "block":
		handleBlock(request, bc)
	case "inv":
		handleInv(request, bc)
	case "getblocks":
		handleGetBlocks(request, bc)
	case "getdata":
		handleGetData(request, bc)
	case "tx":
		handleTx(request, bc)
	case "version":
		handleVersion(request, bc)
	default:
		fmt.Printf("未知命令，垃圾:\n", command)
	}

	conn.Close()

}

//开启服务器
func StartSever(nodeID, minerAddress string) {
	nodeAddress = fmt.Sprintf("localhost:%s", nodeID)
	miningAddress = minerAddress                 //挖矿地址
	In, err := net.Listen(protocol, nodeAddress) //监听
	if err != nil {
		log.Panic(err)
	}
	defer In.Close()
	bc := NewBlockChain(nodeID)
	if nodeAddress != knowNodes[0] {
		sendVersion(knowNodes[0], bc)
	}
	for {
		conn, err := In.Accept() //接收请求
		if err != nil {
			log.Panic(err)
		}
		go handleConnection(conn, bc) //异步处理
	}

}

func gobEncode(data interface{}) []byte {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff) //编码器
	err := enc.Encode(data)      //编码
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes() //字节
}

func nodeIsKnow(addr string) bool { //判断一个节点是不是已经知道的节点
	for _, node := range knowNodes {
		if node == addr {
			return true
		}
	}
	return false
}
