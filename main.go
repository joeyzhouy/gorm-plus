package main

import (
	"github.com/joeyzhouy/gorm-plus/builder"
	_ "github.com/joeyzhouy/gorm-plus/builder/protocol"
	"github.com/joeyzhouy/gorm-plus/utils"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"strings"
)

var (
	app            = kingpin.New("gorm-plus", "A command-line for generate golang entity from database.")
	generate       = app.Command("generate", "create golang entity.")
	userFlag       = app.Flag("u", "user for connect to database, default: root").Default("root").String()
	passwordFlag   = app.Flag("p", "password for connect to database, default: root").Default("root").String()
	portFlag       = app.Flag("port", "specify port, default: 3306").Default("3306").Int()
	dbFlag         = app.Flag("d", "specify database name").Default("dbName").String()
	jsonFlag       = app.Flag("json", "add json annotation").Default("true").Bool()
	hostFlag       = app.Flag("host", "specify host").Default("127.0.0.1").String()
	targetFlag     = app.Flag("target", "specify target dir").String()
	tableNamesFlag = app.Flag("tables", "specify tables to generate, separator \",\"default all tables in database").String()
	protocolFlag   = app.Flag("protocol", "specify protocol, ex: mysql").Default("mysql").String()
)

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case generate.FullCommand():
		if *targetFlag == "" {
			println("target dir no specify")
			return
		}
		param := utils.Param{
			Protocol: *protocolFlag,
			UserName: *userFlag,
			Password: *passwordFlag,
			Port:     *portFlag,
			DBName:   *dbFlag,
			Host:     *hostFlag,
		}
		tables := make([]string, 0)
		if len(*tableNamesFlag) > 0 {
			for _, val := range strings.Split(*tableNamesFlag, ",") {
				tables = append(tables, strings.TrimSpace(val))
			}
		}
		param.TableNames = tables
		param.Attachment = map[string]interface{}{
			utils.JsonKey:   *jsonFlag,
			utils.TargetKey: *targetFlag,
		}
		err := builder.GetBuilder(param).Exec()
		if err != nil {
			println(err.Error())
		}
	}
}
