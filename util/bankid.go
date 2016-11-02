package util

import (
	"fmt"
)

var id map[string]int64
var name map[int64]string

func ID(name string) (int64, error) {
	if i, ok := id[name]; ok {
		return i, nil
	}
	return 0, fmt.Errorf("name %s not found", name)
}
func Name(id int64) (string, error) {
	if n, ok := name[id]; ok {
		return n, nil
	}
	return "", fmt.Errorf("id %d not found", id)
}

func init() {
	id = make(map[string]int64)
	name = make(map[int64]string)

	id["ICBC"] = 1
	name[1] = "ICBC"
	id["ABC"] = 2
	name[2] = "ABC"
	id["BOC"] = 3
	name[3] = "BOC"
	id["CCB"] = 4
	name[4] = "CCB"
	id["BC"] = 5
	name[5] = "BC"
	id["CITIC"] = 6
	name[6] = "CITIC"
	id["CEB"] = 7
	name[7] = "CEB"
	id["HXB"] = 8
	name[8] = "HXB"
	id["CMBC"] = 9
	name[9] = "CMBC"
	id["CGBC"] = 10
	name[10] = "CGBC"
	id["SZFB"] = 11
	name[11] = "SZFB"
	id["CMB"] = 12
	name[12] = "CMB"
	id["CIB"] = 13
	name[13] = "CIB"
	id["SPDB"] = 14
	name[14] = "SPDB"
	id["BOS"] = 16
	name[16] = "BOS"
	id["BOT"] = 17
	name[17] = "BOT"
	id["BOH"] = 18
	name[18] = "BOH"
	id["BON"] = 19
	name[19] = "BON"
	id["GZBC"] = 20
	name[20] = "GZBC"
	id["PINGAN"] = 21
	name[21] = "PINGAN"
	id["BOD"] = 22
	name[22] = "BOD"
	id["163EPAY"] = 23
	name[23] = "163EPAY"
	id["SINAPAY"] = 24
	name[24] = "SIANPAY"
	id["YONGYI"] = 25
	name[25] = "YONGYI"
	id["YONGYI_TD"] = 26
	name[26] = "YONGYI_TD"
	id["EPAYLINKS"] = 27
	name[27] = "EPAYLINKS"
	id["RETENTION"] = 28
	name[28] = "RETENTION"
	id["SIAN_SD"] = 30
	name[30] = "SINA_SD"
}
