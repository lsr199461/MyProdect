package wallent

import (
	"fmt"
	"bytes"
	"encoding/gob"
	"crypto/elliptic"
	"log"
	"io/ioutil"
	"os"
)

type Wallets struct {
	wallets map[string]*Wallet //一个字符串对应一个钱包
}

//创建钱包
func NewWallets() (error, *Wallets) {
	wallets := Wallets{}
	wallets.wallets = make(map[string]*Wallet)
	err := wallets.LoadFromFole()
	return err, &wallets
}

//创建
func (ws *Wallets) CreateWallet() string {
	wallet := NewWallet()
	address := fmt.Sprintf("%s", wallet.GetAddress())
	ws.wallets[address] = wallet
	return address
}

//获取所有钱包
func (ws *Wallets) GetAllWallet() []string {
	var address []string
	for _, addr := range address {
		address = append(address, addr)
	}
	return address
}

//获取一个钱包
func (ws *Wallets) GetWallet(address string) Wallet {
	return *ws.wallets[address]
}

//保存钱包到文件
func (ws *Wallets) SaveToFile(nodeID string) {
	var content bytes.Buffer
	file := walletFile            //生成文件地址
	gob.Register(elliptic.P256()) //加密算法
	encode := gob.NewEncoder(&content)
	err := encode.Encode(ws)
	if err != nil {
		log.Panic(err)
	}
	err = ioutil.WriteFile(file, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}

//从文件中读取钱包
func (ws *Wallets) LoadFromFole() error {
	myfile := walletFile
	_, err := os.Stat(myfile)
	if os.IsNotExist(err) {
		return err
	}
	fileContent, err := ioutil.ReadFile(myfile) //读取文件
	if err != nil {
		log.Panic(err)
	}
	//读二进制文件并解析
	var wallets Wallets
	gob.Register(elliptic.P256()) //加密
	decode := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decode.Decode(&wallets)
	if err != nil {
		log.Panic(err)
	}
	ws.wallets = wallets.wallets
	return err
}
