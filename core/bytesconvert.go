package core

import (
	"bytes"
	"encoding/binary"
	"github.com/imroc/biu"
	"strconv"
	"strings"
)
/// int 转 byte[4]
func  IntToBytes4(value int)([]byte){
	src :=make([]byte,0)
	src =append (src,(byte)(value >> 24) & 0xFF)
	src =append (src,(byte)(value >> 16) & 0xFF)
	src =append (src,(byte)(value >> 8) & 0xFF)
	src =append (src,(byte)(value & 0xFF))
	return src
}

/// 将byte数组b转为一个整数,字节数组的低位是整型的低字节位,直接转;
func   ByteToInt(b []byte)(int) {
	var iOutcome uint = 0
	var bLoop byte
	 length := len(b)
	for  i := 0; i < length; i++{
		bLoop = b[i]
		iOutcome += (uint)(bLoop & 0xFF) << (8 * uint(i))
	}
	return int(iOutcome)
}
// 16进制字符串转16进制字节数组;
func StrToHexByte(hexString string) []byte {
	hexString = strings.Replace(hexString, " ", "", 0)
	if len(hexString)%2 != 0 {
		hexString += " "
	}
	var returnBytes = make([]byte, len(hexString)/2)
	for i := 0; i < len(returnBytes); i++ {
		var temp = Substr(hexString, i*2, 2)
		x, _ := strconv.ParseInt(temp, 16, 32)
		returnBytes[i] = byte(x)
	}
	return returnBytes
}
// 16进制字符串转16进制字节数组;(2B 高低位互换)
func StrToHexByte2(hexString string) []byte {
	hexString = strings.Replace(hexString, " ", "", 0)
	hexString = strings.Replace(hexString, "0x", "", 1)
	if len(hexString)%2 != 0 {
		hexString += " "
	}
	var returnBytes = make([]byte, len(hexString)/2)
	for i := 0; i < len(returnBytes); i++ {
		var temp = Substr(hexString, i*2, 2)
		x, _ := strconv.ParseInt(temp, 16, 32)
		returnBytes[len(returnBytes)-1-i] = byte(x)
	}
	return returnBytes
}
//byte字节数组转为十六进制字符串
func BytesToHexStr(bytes []byte) (hex string) {
	return BinaryStrToHexStr(biu.BytesToBinaryString(bytes[:]))
}

//bytes数组转换为十进制字符串
func BytesConvertIntArr(byteArr []byte) (ints string) {
	buffer := new(bytes.Buffer)
	for _, b := range byteArr {
		s := strconv.FormatInt(int64(b&0xff), 10)
		if len(s) == 1 {
			buffer.WriteString("0")
		}
		buffer.WriteString(s)
	}
	return buffer.String()
}

//bytes数组转换为十六进制字符串
func BytesConvertHexArr(byteArr []byte) (hex string) {
	buffer := new(bytes.Buffer)
	for _, b := range byteArr {
		s := strconv.FormatInt(int64(b&0xff), 16)
		if len(s) == 1 {
			buffer.WriteString("0")
		}
		buffer.WriteString(s)
	}
	return buffer.String()
}

//bytes数组转换为十六进制字符串(2B 高低位互换)
func BytesConvertHexArr2(byteArr []byte) (hex string) {
	buffer := new(bytes.Buffer)
	var normal =make([]byte,0)
	if len(byteArr)>0{
		for i:=len(byteArr)-1;i>=0;i--{
			normal=append(normal,byteArr[i])
		}
	}
	for _, b := range normal {
		s := strconv.FormatInt(int64(b&0xff), 16)
		if len(s) == 1 {
			buffer.WriteString("0")
		}
		buffer.WriteString(s)
	}
	return buffer.String()
}

//二进制字节数数组转十六进制字符串
func BinaryStrToHexStr(str string) string {
	dataArr := strings.Split(str, " ")
	s := make([]string, len(dataArr))
	for i := range dataArr {
		if i == 0 || i == len(dataArr) {
			dataArr[i] = strings.Replace(dataArr[0], "[", "", 1)
		}
		s[i] = Btox(dataArr[i])
	}
	return strings.Join(s, "")
}

//二进制转十六进制
func Btox(b string) string {
	base, _ := strconv.ParseInt(b, 2, 10)
	buffer := new(bytes.Buffer)
	s := strconv.FormatInt(int64(base&0xff), 16)
	if len(s) == 1 {
		buffer.WriteString("0")
	}
	buffer.WriteString(s)
	// 转化为字符串
	//fmt.Println()
	return buffer.String() //strconv.FormatInt(base, 16)
}

//十六进制转二进制
func Xtob(x string) string {
	base, _ := strconv.ParseInt(x, 16, 10)
	return strconv.FormatInt(base, 2)
}

//二进制字节数数字转int数组
func BytesConvert(b []byte) string {
	s := make([]string, len(b))
	for i := range b {
		s[i] = strconv.Itoa(int(b[i]))
	}
	return strings.Join(s, ",")
}

//字节数组转换成int
func BytesToInt(b []byte) int {
	var temp uint = 0
	var bLoop byte
	for i := 0; i < len(b); i++ {
		bLoop = b[i]
		temp += (uint)(bLoop&0xFF) << (8 * uint(i))
	}
	return int(temp)
}

//字节数组转换成int
func BytesToUInt(b []byte) uint {
	var temp uint = 0
	var bLoop byte
	for i := 0; i < len(b); i++ {
		bLoop = b[i]
		temp += (uint)(bLoop&0xFF) << (8 * uint(i))
	}
	return temp
}

//字节数组转换成int
func BytesToInt2(b []byte) int64 {

	//var tmp int
	//for i := 0; i < len(b); i++ {
	//	tmp += int(b[i])
	//}
	var temp int64 = 0
	var bLoop byte
	for i := 0; i < len(b); i++ {
		bLoop = b[i]
		var j = strconv.FormatInt(int64((uint)(bLoop&0xFF)<<(8*uint(i))), 10)
		sum, _ := strconv.ParseInt(j, 16, 64)
		temp += sum //(uint)(bLoop&0xFF) << (8 * uint(i))
	}
	return temp
}

func BytesToInt3(b []byte) int64 {

	//var tmp int
	//for i := 0; i < len(b); i++ {
	//	tmp += int(b[i])
	//}
	var temp uint = 0
	var bLoop byte
	for i := 0; i < len(b); i++ {
		bLoop = b[i]
		//var j = strconv.FormatInt(int64((uint)(bLoop&0xFF) << (8 * uint(i))), 10)
		//sum, _ :=strconv.ParseInt(j, 16, 64)
		temp += (uint)(bLoop&0xFF) << (8 * uint(i))
	}
	return int64(temp)
}

//字节数组转换成int64
func BytesToInt64(b []byte) int64 {
	bytesBuffer := bytes.NewBuffer(b)
	var tmp int64
	binary.Read(bytesBuffer, binary.BigEndian, &tmp)
	return int64(tmp)
}

//Byte字节数组分割
func BytesSplit(bytes []byte, n int) (remains []byte, new []byte) {
	new = make([]byte, n)
	remains = make([]byte, len(bytes)-n)
	for i := 0; i < n; i++ {
		new[i] = bytes[i]
	}
	for i := 0; i < len(bytes)-n; i++ {
		remains[i] = bytes[n+i]
	}
	return remains, new
}

/// <summary>
/// 12位拼凑成3个字节，12*2 = 3*8;
/// 返回两个值;
/// </summary>
/// <param name="b"></param>
/// <param name="bit"></param>
/// <returns></returns>
func ByteToInts(b []byte) (intArr [2]int) {
	var values [2]int
	total := BytesToInt(b)
	values[0] = total & 0XFFF
	values[1] = total >> 12
	return values
}

/// <summary>
/// 8个1位拼凑成1个字节， 1*8 = 1*8；
/// 返回8个值;
/// ex: 0x55;
/// </summary>
/// <param name="b"></param>
func ByteToInts1(b byte) (intArr [8]int) {
	var values [8]int
	str := biu.ByteToBinaryString(b)

	values[0], _ = strconv.Atoi(Substr(str, 7, 1))
	values[1], _ = strconv.Atoi(Substr(str, 6, 1))
	values[2], _ = strconv.Atoi(Substr(str, 5, 1))
	values[3], _ = strconv.Atoi(Substr(str, 4, 1))
	values[4], _ = strconv.Atoi(Substr(str, 3, 1))
	values[5], _ = strconv.Atoi(Substr(str, 2, 1))
	values[6], _ = strconv.Atoi(Substr(str, 1, 1))
	values[7], _ = strconv.Atoi(Substr(str, 0, 1))
	return values
}

//字符串截取
func Substr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}
	return string(rs[start:end])
}
func Uint16ToBytes(n uint16) []byte {
	return []byte{
		byte(n),
		byte(n >> 8),
	}
}

/// <summary>
/// 4个字节表示时间，第一个字节的高两位位标志位，不算做数据;
/// </summary>
/// <param name="b"></param>
/// <returns></returns>
func BytesToTime(b []byte) (time int64) {
	t := 0
	t = int(b[0]&0x3F) << 24
	t += int(b[1]) << 16
	t += int(b[2]) << 8
	t += int(b[3])
	t *= 60
	return int64(t)
}

/// <summary>
/// 时间unix转成byte[4]
/// </summary>
/// <param name="value"></param>
/// <returns></returns>
func UnixTobyte( unixtime int) ( []byte) {
	var byteArr=make([]byte,0)
	unixtime /= 60
	buffer := new(bytes.Buffer)
	binary.Write(buffer,binary.LittleEndian,uint16((unixtime >> 24) | 0X80))
	byteArr=append(byteArr,buffer.Bytes()[0])

	buffer = new(bytes.Buffer)
	binary.Write(buffer,binary.LittleEndian,uint16(unixtime >> 16))
	byteArr=append(byteArr,buffer.Bytes()[0])

	buffer = new(bytes.Buffer)
	binary.Write(buffer,binary.LittleEndian,uint16(unixtime >> 8))
	byteArr=append(byteArr,buffer.Bytes()[0])

	buffer = new(bytes.Buffer)
	binary.Write(buffer,binary.LittleEndian,uint16(unixtime))
	byteArr=append(byteArr,buffer.Bytes()[0])
	return byteArr
}

/// 时间unix转成byte[4],高低位互换
func  UnixTobyte2(unixtime int)( []byte) {
	var byteArr=make([]byte,0)
	unixtime /= 60
	buffer := new(bytes.Buffer)
	binary.Write(buffer,binary.LittleEndian,uint16((unixtime >> 24) | 0X80))
	byteArr=append(byteArr,buffer.Bytes()[0])

	buffer = new(bytes.Buffer)
	binary.Write(buffer,binary.LittleEndian,uint16(unixtime >> 16))
	byteArr=append(byteArr,buffer.Bytes()[0])

	buffer = new(bytes.Buffer)
	binary.Write(buffer,binary.LittleEndian,uint16(unixtime >> 8))
	byteArr=append(byteArr,buffer.Bytes()[0])

	buffer = new(bytes.Buffer)
	binary.Write(buffer,binary.LittleEndian,uint16(unixtime))
	byteArr=append(byteArr,buffer.Bytes()[0])
	var normal =make([]byte,0)
	if len(byteArr)>0{
		for i:=len(byteArr)-1;i>=0;i--{
			normal=append(normal,byteArr[i])
		}
	}
	return normal
}
//BytesCombine 多个[]byte数组合并成一个[]byte
func BytesCombine(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(""))
}
