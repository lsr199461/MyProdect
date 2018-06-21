package main

import (
	"fmt"
	"os"
	"flag"
	"log"
)

//命令行接口
type CLI struct {
	blockchain *BlockChain
}










//用法
func (cli *CLI)printUsage(){
	fmt.Println("用法如下")
	fmt.Println("newwallet 创建钱包")
	fmt.Println("mywallet 显示所有账户")
	fmt.Println("getbalance -address 你输入的地址  根据地址查询金额")
	fmt.Println("blockchain -address 你输入的地址  根据地址创建区块链")
	fmt.Println("send  -from  From -to To -amount  Amount 转账  -mine  ")
	fmt.Println("showchain 显示区块链")
	fmt.Println("reindexutxo 重建索引")
	fmt.Println("startnode  -miner  ADDR 开启一个节点")
	fmt.Println("nodeID  -port  设置端口号")
}
func (cli *CLI)validateArgs(){
	if len(os.Args)<2{
		cli.printUsage()//显示用法
		os.Exit(1)
	}
}




func (cli *CLI)Run(){
	cli.validateArgs()//校验

	//处理命令行参数
	listaddressescmd:=flag.NewFlagSet("mywallet",flag.ExitOnError)
	createwalletcmd:=flag.NewFlagSet("newwallet",flag.ExitOnError)
	getbalancecmd:=flag.NewFlagSet("getbalance",flag.ExitOnError)
	createblockchaincmd:=flag.NewFlagSet("blockchain",flag.ExitOnError)
	sendcmd:=flag.NewFlagSet("send",flag.ExitOnError)
	showchaincmd:=flag.NewFlagSet("showchain",flag.ExitOnError)
	reindexutxocmd:=flag.NewFlagSet("reindexutxo",flag.ExitOnError)
	startnodecmd:=flag.NewFlagSet("startnode",flag.ExitOnError)


	getbalanceaddress:=getbalancecmd.String("address","","查询地址")
	createblockaddress:=createblockchaincmd.String("address","","查询地址")
	sendfrom:=sendcmd.String("from","","谁给的")
	sendto:=sendcmd.String("to","","给谁的")
	sendamount:=sendcmd.Int("amount",0,"金额")
	sendmine:=sendcmd.Bool("mine",false,"是否立刻挖矿")
	startnodeminer:=startnodecmd.String("miner","","开启是否挖矿")





	switch os.Args[1]{

	case "getbalance":
		err:=getbalancecmd.Parse(os.Args[2:])
		if err!=nil{
			log.Panic(err)//处理错误
		}
	case"blockchain":
		err:=createblockchaincmd.Parse(os.Args[2:])
		if err!=nil{
			log.Panic(err)//处理错误
		}
	case"send":
		err:=sendcmd.Parse(os.Args[2:])
		if err!=nil{
			log.Panic(err)//处理错误
		}
	case"newwallet":
		err:=createwalletcmd.Parse(os.Args[2:])
		if err!=nil{
			log.Panic(err)//处理错误
		}
	case"mywallet":
		err:=listaddressescmd.Parse(os.Args[2:])
		if err!=nil{
			log.Panic(err)//处理错误
		}
	case"showchain":
		err:=showchaincmd.Parse(os.Args[2:])
		if err!=nil{
			log.Panic(err)//处理错误
		}
	case"reindexutxo":
		err:=	reindexutxocmd.Parse(os.Args[2:])
		if err!=nil{
			log.Panic(err)//处理错误
		}
	case"startnode":
		err:=startnodecmd.Parse(os.Args[2:])
		if err!=nil{
			log.Panic(err)//处理错误
		}

	default:
		cli.printUsage()
		os.Exit(1)
	}

	var nodeID ="3001"



	if getbalancecmd.Parsed(){
		if  *getbalanceaddress==""{
			getbalancecmd.Usage()
			os.Exit(1)
		}
		cli.getBalance(*getbalanceaddress,nodeID)//查询
	}
	if createblockchaincmd.Parsed(){
		if *createblockaddress==""{
			createblockchaincmd.Usage()
			os.Exit(1)
		}
		cli.createBlockChain(*createblockaddress,nodeID)//创建区块链
	}
	if sendcmd.Parsed(){
		if *sendfrom=="" || *sendto=="" ||*sendamount<=0{
			sendcmd.Usage()
			os.Exit(1)
		}
		cli.send(*sendfrom,*sendto,*sendamount,nodeID,*sendmine)
	}
	if showchaincmd.Parsed(){
		cli.showBlockChain(nodeID)//显示区块链
	}
	if createwalletcmd.Parsed(){
		fmt.Printf("创建钱包")
		cli.createWallet(nodeID)//显示区块链
	}

	if listaddressescmd.Parsed(){
		cli.listAddresses(nodeID)//显示区块链
	}
	if reindexutxocmd.Parsed(){
		cli.reindexUTXOP(nodeID)//重建索引
	}
	if startnodecmd.Parsed(){
		if nodeID == "" {
			startnodecmd.Usage()
			os.Exit(1)
		}
		cli.startNode(nodeID,*startnodeminer)
	}



}
