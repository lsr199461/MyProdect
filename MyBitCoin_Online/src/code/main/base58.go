package main

import "math/big"
import (
	"bytes"
	"fmt"
	//"strconv"
)

//字母表格，最终会展示的字符
var b58Alphabet=[]byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")


func Base58Encode(input []byte)[]byte{
	var  result []byte
	x:=big.NewInt(0).SetBytes(input)//输入的数据存入二进制，

	base:=big.NewInt(int64(len(b58Alphabet)))//创建了一个大数
	zero:=big.NewInt(0)//创建了一个大数
	mod:=&big.Int{}//创建了一个大数
	for  x.Cmp(zero)!=0{
		x.DivMod(x,base,mod)////求余数，
		result=append(result,b58Alphabet[mod.Int64()])
	}
	ReverseBytes(result)
	for myb  :=range input{
		if myb==0x00{
			//不断追加
			result=append([]byte{b58Alphabet[0]},result...)
		}else{
			break
		}
	}

	return result
}
func Base58Decode(input []byte)[]byte{
	result:=big.NewInt(0) //初始化为0
	zeroBytes :=0 //记数
	for  b:=range input{ //循环
		if b==0x00{
			zeroBytes++//叠加
		}
	}
	payload:=input[zeroBytes:]//取出字节
	for _,b :=range payload{
		charIndex:=bytes.IndexByte(b58Alphabet,b)//字母表格
		result.Mul(result,big.NewInt(58))//乘法
		result.Add(result,big.NewInt(int64(charIndex)))//加法
	}
	decoded:=result.Bytes()//解码
	//叠加
	decoded=append(bytes.Repeat([]byte{byte(0x00)},zeroBytes),decoded...)
	return decoded
}

func mainTestBase58(){
	fmt.Println(Base58Encode([]byte("12345")))
	fmt.Println(Base58Decode(Base58Encode([]byte("12345"))))
	fmt.Printf("%s",   Base58Decode(Base58Encode([]byte("12345"))))
}