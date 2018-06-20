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
	fmt.Println("mywallet 显示我的所有钱包")
	fmt.Println("getbalance -address 你输入的地址  根据地址查询金额")
	fmt.Println("firstblock -address 你输入你要挖创世块的钱包地址")
	fmt.Println("send  -from  From -to To -amount  Amount 转账 ")
	fmt.Println("showchain 显示区块链")
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
	createblockchaincmd:=flag.NewFlagSet("createblockchain",flag.ExitOnError)
	sendcmd:=flag.NewFlagSet("send",flag.ExitOnError)
	showchaincmd:=flag.NewFlagSet("showchain",flag.ExitOnError)

	getbalanceaddress:=getbalancecmd.String("address","","查询地址")
	createblockaddress:=createblockchaincmd.String("address","","查询地址")
	sendfrom:=sendcmd.String("from","","谁给的")
	sendto:=sendcmd.String("to","","给谁的")
	sendamount:=sendcmd.Int("amount",0,"金额")





	switch os.Args[1]{
	case "getbalance":
		err:=getbalancecmd.Parse(os.Args[2:])
		if err!=nil{
			log.Panic(err)//处理错误
		}
	case"firstblock":
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
	default:
		cli.printUsage()
		os.Exit(1)
	}
	if getbalancecmd.Parsed(){
		if  *getbalanceaddress==""{
			getbalancecmd.Usage()
			os.Exit(1)
		}
		cli.getBalance(*getbalanceaddress)//查询
	}
	if createblockchaincmd.Parsed(){
		if *createblockaddress==""{
			createblockchaincmd.Usage()
			os.Exit(1)
		}
		cli.createBlockChain(*createblockaddress)//创建区块链
	}
	if sendcmd.Parsed(){
		if *sendfrom=="" || *sendto=="" ||*sendamount<=0{
			sendcmd.Usage()
			os.Exit(1)
		}
		cli.send(*sendfrom,*sendto,*sendamount)
	}
	if showchaincmd.Parsed(){
		cli.showBlockChain()//显示区块链
	}
	if createwalletcmd.Parsed(){
		cli.createWallet()//显示区块链
	}

	if listaddressescmd.Parsed(){
		cli.listAddresses()//显示区块链
	}




}
