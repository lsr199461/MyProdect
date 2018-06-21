package main

import "fmt"

func (cli *CLI) reindexUTXOP(nodeID string){
	blockchain:=NewBlockChain(nodeID)//UTXO统计交易
	UTXOSet:=UTXOSet{blockchain}
	UTXOSet.Reindex()


	count:=UTXOSet.CountTransactions()
	fmt.Printf("已经有%d次交易在UTXO集合\n",count)
}