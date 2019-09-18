package identypes
import (
	"crypto/sha256"


)
/*
tx的交易形式有三种：
类型是string
普通tx：
	tx=内容
跨片tx：
	relayTx,A,B=内容
	addTx,A,B = 内容
	checkpoint=
*/
type TX struct{
	Txtype string
	Sender string
	Receiver string
	ID       [sha256.Size]byte
	Content []string
} 
 //这格式只在relay tx相关的内容中使用

