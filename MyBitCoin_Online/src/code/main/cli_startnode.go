package main
import "fmt"
import (
	"log"
	//"github.com/ethereum/go-ethereum/p2p/simulations/adapters"
)

func (cli *CLI) startNode(nodeID,minerAddress string){
	fmt.Printf("开启一个节点%s\n",nodeID)
	if len(minerAddress)>0{
		if ValidateAddress(minerAddress){
			fmt.Printf("正在挖矿地址是这个%s",minerAddress)
		}else{
			log.Panic("错误的挖矿地址")
		}
	}
	StartSever(nodeID,minerAddress)//开启服务器



}
