
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// started from https://github.com/asdf072/struct-create

package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gedex/inflector"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"log"
	"os"
	"os/exec"
	"strings"
)

var config Configuration

func configurationDefaults() (Configuration, error) {
	if len(os.Args) != 7 {
		return Configuration{}, errors.New("Usage " + os.Args[0] + " <dbtype> <dbuser> <dbpasswd> <dbname> <dbserver> <dbport>")
	}
	cfg := Configuration{
		DbType:     os.Args[1],
		DbUser:     os.Args[2],
		DbPassword: os.Args[3],
		DbName:     os.Args[4],
		DbServer:   os.Args[5],
		DbPort:     os.Args[6],
		PkgName:    "api",
		TagLabel:   "db",
	}
	return cfg, nil
}

type Configuration struct {
	DbType     string `json:"db_type"`
	DbUser     string `json:"db_user"`
	DbPassword string `json:"db_password"`
	DbName     string `json:"db_name"`
	DbServer   string `json:"db_server"`
	DbPort     string `json:"db_port"`
	// PkgName gives name of the package using the stucts
	PkgName string `json:"pkg_name"`
	// TagLabel produces tags commonly used to match database field names with Go struct members
	TagLabel string `json:"tag_label"`
}

type ColumnSchema struct {
	TableName              string
	ColumnName             string
	IsNullable             string
	DataType               string
	CharacterMaximumLength sql.NullInt64
	NumericPrecision       sql.NullInt64
	NumericScale           sql.NullInt64
	ColumnType             string
	ColumnKey              string
	ColumnForeignTable     string
	ColumForeignColumn     string
	PrimaryKey             bool

	GoType         string
	RequiredImport string
}

type FKSchema struct {
	ConstraintName    string
	TableName         string
	columnName        string
	ForeignTableName  string
	ForeignColumnName string
}

// \todo fix for compound keys
func primaryKey(schemas []ColumnSchema, table string) []ColumnSchema {
	var pkColumns []ColumnSchema
	for _, cs := range schemas {
		if cs.TableName != table || !cs.PrimaryKey {
			continue
		}
		pkColumns = append(pkColumns, cs)
	}

	if pkColumns == nil {
		//	return "Links." + formatName(cs.ColumnName) + "Link.ID"
		log.Fatal("Table " + table + " without primary key")
		panic("Table " + table + " without primary key")
	}

	return pkColumns
}

func writeFile(schemas []ColumnSchema, table string) (int, error) {
	file, err := os.Create("./generated/" + table + ".go")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	license := "// Licensed under the Apache License, Version 2.0 (the \"License\");\n"
	license += "// you may not use this file except in compliance with the License.\n"
	license += "// You may obtain a copy of the License at\n\n"
	license += "// http://www.apache.org/licenses/LICENSE-2.0\n\n"
	license += "// Unless required by applicable law or agreed to in writing, software\n"
	license += "// distributed under the License is distributed on an \"AS IS\" BASIS,\n"
	license += "// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.\n"
	license += "// See the License for the specific language governing permissions and\n"
	license += "// limitations under the License.\n\n"
	license += "// This file was initially generated by gen_to_start.go (add link), as a start\n"
	license += "// of the Traffic Ops golang data model\n\n"

	header := "package " + config.PkgName + "\n\n"
	header += "import (\n"
	header += "\"log\"\n"
	header += "\"github.com/jmoiron/sqlx\"\n"

	sString := structString(schemas, table)

	if strings.Contains(sString, "null.") {
		header += "null \"gopkg.in/guregu/null.v3\"\n"
	}
	header += "_ \"github.com/apache/trafficcontrol/traffic_ops/experimental/server/output_format\" // needed for swagger\n"
	if strings.Contains(sString, "time.") {
		header += "\"time\"\n"
	}
	header += "\"encoding/json\"\n"
	header += ")\n\n"

	hString := handleString(schemas, table)
	totalBytes, err := fmt.Fprint(file, license+header+sString+hString)
	if err != nil {
		log.Fatal(err)
	}

	return totalBytes, nil
}

// gen a list of columnames without id and last_updated
func colString(schemas []ColumnSchema, table string, prefix string, varName string) string {
	out := ""
	sep := ""
	for _, cs := range schemas {
		if cs.TableName == table && cs.ColumnName != "id" && cs.ColumnName != "last_updated" {
			out += varName + "+= \"" + sep + prefix + cs.ColumnName + "\"\n"
			sep = ","
		}
	}
	return out
}

func genInsertVarLines(schemas []ColumnSchema, table string) string {
	out := "sqlString := \"INSERT INTO " + table + "(\"\n"
	out += colString(schemas, table, "", "sqlString")
	out += "sqlString += \") VALUES (\"\n"
	out += colString(schemas, table, ":", "sqlString")
	out += "sqlString += \")\"\n"

	return out
}

func updString(schemas []ColumnSchema, table string, prefix string, varName string) string {
	out := ""
	sep := ""
	for _, cs := range schemas {
		if cs.TableName == table && cs.ColumnName != "id" {
			out += varName + "+= \"" + sep + cs.ColumnName + " = :" + cs.ColumnName + "\"\n"
			sep = ","
		}
	}
	return out
}

func genUpdateVarLines(schemas []ColumnSchema, table string) string {
	pk := primaryKey(schemas, table)
	out := "sqlString := \"UPDATE " + table + " SET \"\n"
	out += updString(schemas, table, "", "sqlString")
	out += "sqlString += \" WHERE "
	for _, col := range pk {
		if col.ColumForeignColumn == "" {
			out += col.ColumnName + "=:" + col.ColumnName + " AND "
		} else {
			out += col.ColumnName + "=:Links." + formatName(col.ColumnForeignTable) + "Link.ID" + " AND "
		}
	}
	out = out[:len(out)-5] // strip final " AND "
	out += "\"\n"
	return out
}

func hasLastUpdated(schemas []ColumnSchema, table string) bool {
	for _, cs := range schemas {
		if cs.TableName == table {
			if cs.ColumnName == "last_updated" {
				return true
			}
		}
	}
	return false
}

// func genApiPostDocChangeLines(schemas []ColumnSchema, table string) string {
// 	out := ""
// 	for _, cs := range schemas {
// 		if cs.TableName == table && cs.ColumnName != "id" && cs.ColumnName != "last_updated" {
// 			goType, _, err := goType(&cs)
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			goType = strings.Replace(goType, "null.", "", 1)
// 			goType = strings.Replace(goType, string(goType[0]), strings.ToLower(string(goType[0])), 1)
// 			if goType == "float" {
// 				goType = strings.Replace(goType, "float", "float64", 1)
// 			}
// 			nullable := "false"
// 			if cs.IsNullable == "YES" {
// 				nullable = "true"
// 			}
// 			out += fmt.Sprintf("// @Param %20s json %10s %7s \"%s description\"\n",
// 				formatName(cs.ColumnName), goType, nullable, cs.ColumnName)
// 		}
// 	}
// 	return out
// }

// func genApiPutDocChangeLines(schemas []ColumnSchema, table string) string {
// 	out := ""
// 	for _, cs := range schemas {
// 		if cs.TableName == table && cs.ColumnName != "id" && cs.ColumnName != "last_updated" {
// 			goType, _, err := goType(&cs)
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			goType = strings.Replace(goType, "null.", "", 1)
// 			goType = strings.Replace(goType, string(goType[0]), strings.ToLower(string(goType[0])), 1)
// 			if goType == "float" {
// 				goType = strings.Replace(goType, "float", "float64", 1)
// 			}
// 			nullable := "false"
// 			if cs.IsNullable == "YES" {
// 				nullable = "true"
// 			}
// 			out += fmt.Sprintf("// @Param %20s json %10s %7s \"%s description\"\n",
// 				formatName(cs.ColumnName), goType, nullable, cs.ColumnName)
// 		}
// 	}
// 	return out
// }

// getPkGoFuncParamString returns the param string for a Go func being generated,
// for the given primary keys.
func getPkGoFuncParamString(primaryKey []ColumnSchema) string {
	var s string
	for _, column := range primaryKey {
		s += formatNameLower(column.ColumnName) + " " + column.GoType + ", "
	}
	s = s[:len(s)-2] // strip the last ", "
	return s
}

func selfQueryStr(pk []ColumnSchema, table string) string {
	if len(pk) == 1 {
		return "concat('\" + API_PATH + \"" + table + "/', " + pk[0].ColumnName + ") as self"
	}
	s := "concat('\" + API_PATH + \"" + table + "'"
	for _, col := range pk {
		s += ", '/" + col.ColumnName + "/', " + col.ColumnName
	}
	s += ") as self"
	return s
}

func setStructPkFields(pk []ColumnSchema) string {
	var s string
	for _, col := range pk {
		if col.ColumForeignColumn == "" {
			s += "    arg." + formatName(col.ColumnName) + "= " + formatNameLower(col.ColumnName) + "\n"
		} else {
			s += "    arg.Links." + formatName(col.ColumnForeignTable) + "Link.ID = " + formatNameLower(col.ColumnName) + "\n"
		}
	}
	return s
}

func pkWhereStr(pk []ColumnSchema) string {
	s := "WHERE "
	for _, col := range pk {
		s += col.ColumnName + "=:"
		if col.ColumForeignColumn == "" {
			s += col.ColumnName
		} else {
			s += "Links." + formatName(col.ColumnForeignTable) + "Link.ID"
		}
		s += " AND "
	}
	s = s[:len(s)-5] // strip last " AND "
	return s
}

func setFkHALQueryStr(schemas []ColumnSchema, table string) string {
	var s string
	for _, col := range schemas {
		if col.TableName != table || col.ColumnForeignTable == "" {
			continue
		}
		s += "queryStr += \", concat('\" + API_PATH + \"" + col.ColumnForeignTable + "/', " + col.ColumnName + ") as "
		s += col.ColumnForeignTable + "_" + col.ColumForeignColumn + "_ref\"\n"
	}
	return s
}

func handleString(schemas []ColumnSchema, table string) string {
	pk := primaryKey(schemas, table)
	updateLastUpdated := hasLastUpdated(schemas, table)

	out := ""
	out += "// @Title get" + formatName(table) + "ById\n"
	out += "// @Description retrieves the " + table + " information for a certain id\n"
	out += "// @Accept  application/json\n"
	out += "// @Param   id              path    int     false        \"The row id\"\n"
	out += "// @Success 200 {array}    " + formatName(table) + "\n"
	out += "// @Resource /api/2.0\n"
	out += "// @Router /api/2.0/" + table + "/{id} [get]\n"
	out += "func get" + inflector.Singularize(formatName(table)) + "(" + getPkGoFuncParamString(pk) + ", db *sqlx.DB) (interface{}, error) {\n"
	out += "    ret := []" + formatName(table) + "{}\n"
	out += "    arg := " + formatName(table) + "{}\n"
	out += setStructPkFields(pk)
	out += "    queryStr := \"select *, " + selfQueryStr(pk, table) + "\"\n"
	out += setFkHALQueryStr(schemas, table)
	out += "    queryStr += \" from " + table + " " + pkWhereStr(pk) + "\"\n"
	out += "    nstmt, err := db.PrepareNamed(queryStr)\n"
	out += "    err = nstmt.Select(&ret, arg)\n"
	out += "	if err != nil {\n"
	out += "	    log.Println(err)\n"
	out += "	    return nil, err\n"
	out += "	}\n"
	out += "    nstmt.Close()\n"
	out += "	return ret, nil\n"
	out += "}\n\n"

	out += "// @Title get" + formatName(table) + "s\n"
	out += "// @Description retrieves the " + table + "\n"
	out += "// @Accept  application/json\n"
	out += "// @Success 200 {array}    " + formatName(table) + "\n"
	out += "// @Resource /api/2.0\n"
	out += "// @Router /api/2.0/" + table + " [get]\n"
	out += "func get" + inflector.Pluralize(formatName(table)) + "(db *sqlx.DB) (interface{}, error) {\n"
	out += "    ret := []" + formatName(table) + "{}\n"
	out += "    queryStr := \"select *, " + selfQueryStr(pk, table) + "\"\n"
	out += setFkHALQueryStr(schemas, table)
	out += "queryStr += \" from " + table + "\"\n"
	out += "	err := db.Select(&ret, queryStr)\n"
	out += "	if err != nil {\n"
	out += "	   log.Println(err)\n"
	out += "	   return nil, err\n"
	out += "	}\n"
	out += "	return ret, nil\n"
	out += "}\n\n"

	out += "// @Title post" + formatName(table) + "\n"
	out += "// @Description enter a new " + table + "\n"
	out += "// @Accept  application/json\n"
	out += "// @Param                 Body body     " + formatName(table) + "   true \"" + formatName(table) + " object that should be added to the table\"\n"
	out += "// @Success 200 {object}    output_format.ApiWrapper\n"
	out += "// @Resource /api/2.0\n"
	out += "// @Router /api/2.0/" + table + " [post]\n"
	out += "func post" + inflector.Singularize(formatName(table)) + "(payload []byte, db *sqlx.DB) (interface{}, error) {\n"
	out += "	var v " + formatName(table) + "\n"
	out += "	err := json.Unmarshal(payload, &v)\n"
	out += "	if err != nil {\n"
	out += "		log.Println(err)\n"
	out += "    	return nil, err\n"
	out += "	}\n"
	out += genInsertVarLines(schemas, table)
	out += "    result, err := db.NamedExec(sqlString, v)\n"
	out += "    if err != nil {\n"
	out += "        log.Println(err)\n"
	out += "    	return nil, err\n"
	out += "    }\n"
	out += "    return result, err\n"
	out += "}\n\n"

	out += "// @Title put" + formatName(table) + "\n"
	out += "// @Description modify an existing " + table + "entry\n"
	out += "// @Accept  application/json\n"
	out += "// @Param   id              path    int     true        \"The row id\"\n"
	out += "// @Param                 Body body     " + formatName(table) + "   true \"" + formatName(table) + " object that should be added to the table\"\n"
	out += "// @Success 200 {object}    output_format.ApiWrapper\n"
	out += "// @Resource /api/2.0\n"
	out += "// @Router /api/2.0/" + table + "/{id}  [put]\n"
	out += "func put" + inflector.Singularize(formatName(table)) + "(" + getPkGoFuncParamString(pk) + ", payload []byte, db *sqlx.DB) (interface{}, error) {\n"
	out += "    var arg " + formatName(table) + "\n"
	out += "    err := json.Unmarshal(payload, &arg)\n"
	out += setStructPkFields(pk)
	out += "    if err != nil {\n"
	out += "    	log.Println(err)\n"
	out += "    	return nil, err\n"
	out += "    }\n"
	if updateLastUpdated {
		out += "    arg.LastUpdated = time.Now()\n"
	}
	out += genUpdateVarLines(schemas, table)
	out += "    result, err := db.NamedExec(sqlString, arg)\n"
	out += "    if err != nil {\n"
	out += "    	log.Println(err)\n"
	out += "    	return nil, err\n"
	out += "    }\n"
	out += "    return result, err\n"
	out += "}\n\n"

	out += "// @Title del" + formatName(table) + "ById\n"
	out += "// @Description deletes " + table + " information for a certain id\n"
	out += "// @Accept  application/json\n"
	out += "// @Param   id              path    int     false        \"The row id\"\n"
	out += "// @Success 200 {array}    " + formatName(table) + "\n"
	out += "// @Resource /api/2.0\n"
	out += "// @Router /api/2.0/" + table + "/{id} [delete]\n"
	out += "func del" + inflector.Singularize(formatName(table)) + "(" + getPkGoFuncParamString(pk) + ", db *sqlx.DB) (interface{}, error) {\n"
	out += "    arg := " + formatName(table) + "{}\n"
	out += setStructPkFields(pk)
	out += "    result, err := db.NamedExec(\"DELETE FROM " + table + " " + pkWhereStr(pk) + "\", arg)\n"
	out += "    if err != nil {\n"
	out += "    	log.Println(err)\n"
	out += "    	return nil, err\n"
	out += "    }\n"
	out += "    return result, err\n"
	out += "}\n\n"
	return out
}

func structString(schemas []ColumnSchema, table string) string {
	idColumnSchemas := primaryKey(schemas, table)
	idColumnSchema := idColumnSchemas[0] // debug
	idColumnType := idColumnSchema.GoType

	out := "type " + formatName(table) + " struct{\n"
	linkMap := make(map[string]int)
	for i, cs := range schemas {
		if cs.TableName == table {
			// fmt.Println(cs)
			if cs.ColumForeignColumn == "" {
				out = out + "\t" + formatName(cs.ColumnName) + " " + cs.GoType
				if len(config.TagLabel) > 0 {
					out = out + "\t`" + config.TagLabel + ":\"" + cs.ColumnName + "\" json:\"" + formatNameLower(cs.ColumnName) + "\"`"
				}
				out = out + "\n"
			} else {
				// fmt.Println(cs, ">"+cs.ColumForeignColumn+"<")
				// out = out + "\t" + formatName(cs.ColumnName) + " >>>> " + cs.GoType
				linkMap[cs.ColumnForeignTable] = i
			}
		}
	}
	out += "\tLinks " + formatName(table) + "Links `json:\"_links\" db:-`\n"
	out = out + "}\n\n"

	out += "type " + formatName(table) + "Links struct {\n"
	out += "\tSelf string `db:\"self\" json:\"_self\"`\n"

	for fk, _ := range linkMap {
		typeName := formatName(fk)
		if strings.HasSuffix(typeName, "Cachegroup") {
			typeName = "Cachegroup"
		}
		out += "\t\t" + formatName(fk) + "Link " + typeName + "Link `json:\"" + fk + "\" db:-`\n"
	}
	out += "} \n\n"

	for index, cs := range schemas {
		if cs.ColumnForeignTable == table {
			out += "type " + formatName(table) + "Link struct { \n"
			out += "\tID  " + idColumnType + " `db:\"" + inflector.Singularize(table) + "\" json:\"" + schemas[index].ColumForeignColumn + "\"`\n"
			out += "\tRef string `db:\"" + schemas[index].ColumnForeignTable + "_" + schemas[index].ColumForeignColumn + "_ref\" json:\"_ref\"`\n"
			out += "}\n\n"
			break
		}
	}
	return out
}

func getSchema() ([]ColumnSchema, []string, error) {
	columns := []ColumnSchema{}
	tables := []string{}
	database := "information_schema"
	if config.DbType == "mysql" {
		connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=True", config.DbUser, config.DbPassword, config.DbServer, config.DbPort, database)
		conn, err := sql.Open(config.DbType, connStr)
		if err != nil {
			return nil, nil, err
		}
		defer conn.Close()

		q := "SELECT TABLE_NAME, COLUMN_NAME, IS_NULLABLE, DATA_TYPE, " +
			"CHARACTER_MAXIMUM_LENGTH, NUMERIC_PRECISION, NUMERIC_SCALE, COLUMN_TYPE, " +
			"COLUMN_KEY FROM COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION"
		rows, err := conn.Query(q, config.DbName)
		if err != nil {
			return nil, nil, err
		}

		for rows.Next() {
			cs := ColumnSchema{}
			err := rows.Scan(&cs.TableName, &cs.ColumnName, &cs.IsNullable, &cs.DataType,
				&cs.CharacterMaximumLength, &cs.NumericPrecision, &cs.NumericScale,
				&cs.ColumnType, &cs.ColumnKey)
			if err != nil {
				return nil, nil, err
			}

			// this could be moved to a wrapper function
			goType, requiredImport, err := goType(&cs)
			if err != nil {
				return nil, nil, err
			}
			cs.GoType = goType
			cs.RequiredImport = requiredImport

			columns = append(columns, cs)
		}
		if err := rows.Err(); err != nil {
			return nil, nil, err
		}

		q = "select TABLE_NAME from tables WHERE TABLE_SCHEMA = ? AND table_type='BASE TABLE'"
		rows, err = conn.Query(q, config.DbName)
		if err != nil {
			return nil, nil, err
		}

		for rows.Next() {
			var tableName string
			err := rows.Scan(&tableName)
			if err != nil {
				return nil, nil, err
			}
			tables = append(tables, tableName)
		}

	} else if config.DbType == "postgres" {
		connStr := fmt.Sprintf("dbname=%s user=%s password=%s sslmode=disable host=%s port=%s", config.DbName, config.DbUser, config.DbPassword, config.DbServer, config.DbPort)
		conn, err := sql.Open(config.DbType, connStr)
		if err != nil {
			return nil, nil, err
		}
		defer conn.Close()

		q := "SELECT TABLE_NAME, COLUMN_NAME, IS_NULLABLE, DATA_TYPE, " +
			"CHARACTER_MAXIMUM_LENGTH, NUMERIC_PRECISION, NUMERIC_SCALE " +
			"FROM information_schema.COLUMNS ORDER BY TABLE_NAME, ORDINAL_POSITION"
		rows, err := conn.Query(q)
		if err != nil {
			return nil, nil, err
		}

		for rows.Next() {
			cs := ColumnSchema{}
			err := rows.Scan(&cs.TableName, &cs.ColumnName, &cs.IsNullable, &cs.DataType,
				&cs.CharacterMaximumLength, &cs.NumericPrecision, &cs.NumericScale)
			cs.ColumForeignColumn = ""
			cs.ColumnForeignTable = ""
			if err != nil {
				return nil, nil, err
			}

			// this could be moved to a wrapper function
			goType, requiredImport, err := goType(&cs)
			if err != nil {
				return nil, nil, err
			}
			cs.GoType = goType
			cs.RequiredImport = requiredImport

			columns = append(columns, cs)
		}
		if err := rows.Err(); err != nil {
			return nil, nil, err
		}

		q = "select TABLE_NAME from information_schema.tables where table_type='BASE TABLE' and table_schema='public';" // TODO make schema param
		rows, err = conn.Query(q)
		if err != nil {
			return nil, nil, err
		}

		for rows.Next() {
			var tableName string
			err := rows.Scan(&tableName)
			if err != nil {
				return nil, nil, err
			}
			tables = append(tables, tableName)
		}

		// this query could probably be combined into one of the previous ones by someone smarter than me.
		q = `SELECT
    			tc.constraint_name, tc.table_name, kcu.column_name, 
    			ccu.table_name AS foreign_table_name,
    			ccu.column_name AS foreign_column_name 
			FROM 
    			information_schema.table_constraints AS tc 
    			JOIN information_schema.key_column_usage AS kcu
      			ON tc.constraint_name = kcu.constraint_name
    			JOIN information_schema.constraint_column_usage AS ccu
      			ON ccu.constraint_name = tc.constraint_name
			WHERE constraint_type = 'FOREIGN KEY'`
		rows, err = conn.Query(q)
		if err != nil {
			return nil, nil, err
		}

		for rows.Next() {
			fk := FKSchema{}
			err := rows.Scan(&fk.ConstraintName, &fk.TableName, &fk.columnName, &fk.ForeignTableName, &fk.ForeignColumnName)
			if err != nil {
				return nil, nil, err
			}
			for i, _ := range columns {
				if columns[i].ColumnName == fk.columnName && columns[i].TableName == fk.TableName {
					fmt.Println("Setting fk " + fk.ForeignTableName + "." + fk.ForeignColumnName + " for " + columns[i].TableName + "." + columns[i].ColumnName)
					columns[i].ColumnForeignTable = fk.ForeignTableName
					columns[i].ColumForeignColumn = fk.ForeignColumnName
					break
				}
			}
		}

		q = `select tc.table_name, ccu.column_name from information_schema.table_constraints as tc join information_schema.constraint_column_usage as ccu on ccu.constraint_name = tc.constraint_name where constraint_type = 'PRIMARY KEY';`
		rows, err = conn.Query(q)
		if err != nil {
			return nil, nil, err
		}

		for rows.Next() {
			var table string
			var column string
			err := rows.Scan(&table, &column)
			if err != nil {
				return nil, nil, err
			}

			for i, _ := range columns {
				if columns[i].ColumnName != column || columns[i].TableName != table {
					continue
				}
				columns[i].PrimaryKey = true
			}
		}
	}
	return columns, tables, nil
}

func formatName(name string) string {
	parts := strings.Split(name, "_")
	newName := ""
	for _, p := range parts {
		if len(p) < 1 {
			continue
		}
		newName = newName + strings.Replace(p, string(p[0]), strings.ToUpper(string(p[0])), 1)
	}
	return newName
}

func formatNameLower(name string) string {
	newName := formatName(name)
	newName = strings.Replace(newName, string(newName[0]), strings.ToLower(string(newName[0])), 1)
	return newName
}

// goType returns the Go type name, the text of any required import for the given type,
// and an error if no Go type has been defined for the database type.
func goType(col *ColumnSchema) (string, string, error) {
	var requiredImport string
	if col.IsNullable == "YES" {
		requiredImport = "database/sql"
	}
	var gt string
	switch col.DataType {
	case "name", "regproc", "\"char\"", "oid", "pg_node_tree", "ARRAY", "timestamp with time zone", "xid", "bytea", "pg_lsn", "abstime", "anyarray", "interval": // these could be made smarter Go types
		fallthrough
	case "inet": // TODO(make a Go struct type?)
		fallthrough
	case "char", "varchar", "enum", "text", "longtext", "mediumtext", "tinytext", "character varying":
		if col.IsNullable == "YES" {
			gt = "null.String"
		} else {
			gt = "string"
		}
	case "blob", "mediumblob", "longblob", "varbinary", "binary":
		gt = "[]byte"
	case "date", "time", "datetime", "timestamp", "tstamp", "timestamp without time zone":
		gt, requiredImport = "time.Time", "time"
	case "tinyint", "smallint", "int", "mediumint", "bigint", "numeric", "integer":
		if col.IsNullable == "YES" {
			gt = "null.Int"
		} else {
			gt = "int64"
		}
	case "float", "decimal", "double", "double precision", "real":
		if col.IsNullable == "YES" {
			gt = "null.Float"
		} else {
			gt = "float64"
		}
	case "boolean":
		if col.IsNullable == "YES" {
			gt = "null.Bool"
		} else {
			gt = "bool"
		}
	default:
		return "", "", errors.New("No compatible datatype (" + col.DataType + ") for " + col.TableName + "." + col.ColumnName + " found")
	}
	return gt, requiredImport, nil
}

func printUsage(err error) {
	fmt.Println(err.Error())
}

func main() {
	var err error
	config, err = configurationDefaults()
	if err != nil {
		printUsage(err)
		return
	}

	columns, tables, err := getSchema()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(tables)
	for _, table := range tables {
		if table == "goose_db_version" {
			continue
		}
		bytes, err := writeFile(columns, table)
		if err != nil {
			log.Fatal(err)
		}
		cmd := exec.Command("go", "fmt", "./generated/"+table+".go")
		err = cmd.Run()
		if err != nil {
			log.Fatalf("gofmt error: %v", err)
		}
		fmt.Printf("%s: Ok %d\n", table, bytes)
	}
}