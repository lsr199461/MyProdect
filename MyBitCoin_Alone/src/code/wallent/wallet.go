package wallent

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"log"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"code/tools"
	"bytes"
)

const (
	Version            = byte(0x00)    //钱包版本
	walletFile        = "wallet.dat" //钱包文件
	AddressCheckSumlen = 4             //检测地址长度
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey //钱包私钥
	PublicKey  []byte           //收款地址
}

//创建一个钱包
func NewWallet() *Wallet {
	private, public := NewKyePair()
	wallet := Wallet{
		private,
		public,
	}
	return &wallet
}

//创建公私钥
func NewKyePair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()                              //加密算法
	private, err := ecdsa.GenerateKey(curve, rand.Reader) //生成Key
	if err != nil {
		log.Panic(err)
	}
	public := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, public
}

//公钥校验
func CheckSum(payload []byte) []byte {
	firstSha := sha256.Sum256(payload)
	secondSha := sha256.Sum256(firstSha[:])
	return secondSha[:AddressCheckSumlen]
}

//公钥Hash处理
func HashPubKey(pubkey []byte) []byte {
	pubSha := sha256.Sum256(pubkey)
	r160hash := ripemd160.New() //创建一个算法对象
	_, err := r160hash.Write(pubSha[:])
	if err != nil {
		log.Panic(err)
	}
	pubr160 := r160hash.Sum(nil)
	return pubr160
}

//抓取钱包地址
func (w *Wallet) GetAddress() []byte {
	pubKeyHash := HashPubKey(w.PublicKey) //取得Hash值
	versionPayload := append([]byte{Version}, pubKeyHash...)
	checksum := CheckSum(versionPayload)              //检测版本和公钥
	allPayLoad := append(versionPayload, checksum...) //叠加校验
	address := tools.Base58Encode(allPayLoad)         //编码
	return address
}

//校验钱包地址
func ValidAddress(address string) bool {
	pubHash := tools.Base58Decode([]byte(address)) //解码
	actualchecksum := pubHash[len(pubHash)-AddressCheckSumlen:]
	version := pubHash[0] //取得钱包版本
	pubHash = pubHash[1 : len(pubHash)-AddressCheckSumlen]
	targetCheckSum := CheckSum(append([]byte{version}, pubHash...))
	return bytes.Compare(actualchecksum,targetCheckSum)==0
}
