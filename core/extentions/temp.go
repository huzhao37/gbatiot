package extensions

/// <summary>
/// 温度 有符号数转化;
/// </summary>
/// <param name="source"></param>
/// <returns></returns>
func TempTranster(source int) float32 {
	var des float32 = float32(source)
	var temp uint = uint(source)
	if (temp & 0x800) != 0 {
		temp = temp | 0xfffff000
	}
	des =float32(int32(temp))
	return des
}
