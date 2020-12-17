package builder

import (
	"fmt"
	"gorm-plus/extension"
	"gorm-plus/utils"
	"os"
	"strings"
)

type Builder interface {
	Exec() error
}

type builder struct {
	Param utils.Param
	extension.Protocol
}

func GetBuilder(param utils.Param) Builder {
	var err error
	b := &builder{Param: param}
	b.Protocol, err = extension.GetProtocol(param.Protocol, param)
	if err != nil {
		panic(err)
	}
	return b
}

func (b *builder) Exec() error {
	_, err := b.GetConnection()
	if err != nil {
		return err
	}
	if len(b.Param.TableNames) == 0 {
		b.Param.TableNames, err = b.GetTableNames()
		if err != nil {
			fmt.Println("get tableNames error: " + err.Error())
			return err
		}
	}
	tableInfos, err := b.GetColumnWithTableNames(b.Param.TableNames)
	if err != nil {
		fmt.Println("get table column info error: " + err.Error())
		return err
	}
	return b.createFile(tableInfos)
}

var t = map[int]string{
	0: "\t",
	1: "\t",
	2: "\t\t",
	3: "\t\t\t",
	4: "\t\t\t\t",
	5: "\t\t\t\t\t",
}

func getT(total, l int) string {
	index := (total-l)/4
	if (total-l)%4 == 0 {
		index--
	}
	return t[index]
}

func (b *builder) createFile(infos []utils.TableInfo) error {
	if len(infos) == 0 {
		return nil
	}
	json := b.Param.Attachment[utils.JsonKey].(bool)
	targetDir := b.Param.Attachment[utils.TargetKey].(string)
	if targetDir[len(targetDir)-1] != '/' {
		targetDir += "/"
	}
	var targetFile *os.File
	var err error
	packageName := getPackageName(targetDir)
	if targetFile, err = os.Create(targetDir + packageName + utils.GoFileSuffix); err != nil {
		return err
	}
	tempStructStr := ""
	var structName, name1, name2, annotation string
	var importArr = make(map[string]bool)
	for _, info := range infos {
		_, structName = convertFiledName(info.Name)
		tempStructStr += "type " + structName + " struct {\n"
		for _, column := range info.Columns {
			if strings.HasPrefix(column.GoType, "sql.") {
				importArr["database/sql"] = true
			} else if strings.HasPrefix(column.GoType, "time.") {
				importArr["time"] = true
			}
			name1, name2 = convertFiledName(column.Name)
			annotation = fmt.Sprintf("gorm:\"column:%s%s\"", column.Name, column.Key)
			if json {
				annotation += fmt.Sprintf(" json:\"%s\"", name1)
			}

			tempStructStr += fmt.Sprintf("\t%s%s%s%s`%s`", name2, getT(20, len(column.Name)), column.GoType, getT(16, len(column.GoType)), annotation)
			if column.Comment != "" {
				tempStructStr += fmt.Sprintf("// %s\n", column.Comment)
			} else {
				tempStructStr += "\n"
			}
		}
		tempStructStr += "}\n\n"
		tempStructStr += "func (" + strings.ToLower(string(info.Name[0])) + " *" + structName + ") TableName() string {\n" +
			"\treturn \"" + info.Name + "\"\n" +
			"}\n\n"

	}
	temp := "package " + packageName + "\n\n"
	if len(importArr) == 1 {
		for k, _ := range importArr {
			temp += "import \"" + k + "\""
		}

	} else if len(importArr) > 1 {
		temp += "import (\n"
		for k, _ := range importArr {
			temp += "\t\"" + k + "\"\n"
		}
		temp += ")\n\n"
	}
	tempStructStr = temp + tempStructStr
	if _, err = targetFile.WriteString(tempStructStr); err != nil {
		return err
	}
	if targetFile != nil {
		if err = targetFile.Close(); err != nil {
			return err
		}
	}
	return nil
}

func getPackageName(dirPath string) string {
	rs := []rune(dirPath)
	rs = rs[0 : len(rs)-1]
	index := len(rs) - 1
	for i := index; i >= 0; i-- {
		if rs[i] == '/' {
			index = i + 1
			break
		}
	}
	return string(rs[index:])
}

func convertFiledName(dbColumnName string) (string, string) {
	bs := []byte(dbColumnName)
	if bs[0] >= '0' && bs[0] <= '9' {
		bs[0] += 49
		temp := bs[1:]
		bs = append(bs[0:1], '_')
		bs = append(bs, temp...)
	}
	i, j := 0, 0
	ok := false
	for j < len(bs) {
		if bs[j] != '_' {
			if j != i {
				bs[i], bs[j] = bs[j], bs[i]
				if ok && bs[i] >= 'a' && bs[i] <= 'z' {
					bs[i] -= 32
				}
				ok = false
			}
			i++
		} else{
			ok = true
		}
		j++
	}
	bs = bs[0:i]
	r1 := string(bs)
	if bs[0] <= 'Z' && bs[0] >= 'A' {
		bs[0] += 32
		return string(bs), r1
	} else if bs[0] >= 'a' && bs[0] <= 'z' {
		bs[0] -= 32
		return r1, string(bs)
	}
	return r1, r1
}
