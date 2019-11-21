
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// This file was initially generated by gen_to_start.go (add link), as a start
// of the Traffic Ops golang data model

package api

import (
	"encoding/json"
	_ "github.com/apache/trafficcontrol/traffic_ops/experimental/server/output_format" // needed for swagger
	"github.com/jmoiron/sqlx"
	null "gopkg.in/guregu/null.v3"
	"log"
	"time"
)

type Log struct {
	Id        int64       `db:"id" json:"id"`
	Level     null.String `db:"level" json:"level"`
	Message   string      `db:"message" json:"message"`
	Username  string      `db:"username" json:"username"`
	Ticketnum null.String `db:"ticketnum" json:"ticketnum"`
	CreatedAt time.Time   `db:"created_at" json:"createdAt"`
	Links     LogLinks    `json:"_links" db:-`
}

type LogLinks struct {
	Self string `db:"self" json:"_self"`
}

// @Title getLogById
// @Description retrieves the log information for a certain id
// @Accept  application/json
// @Param   id              path    int     false        "The row id"
// @Success 200 {array}    Log
// @Resource /api/2.0
// @Router /api/2.0/log/{id} [get]
func getLog(id int64, db *sqlx.DB) (interface{}, error) {
	ret := []Log{}
	arg := Log{}
	arg.Id = id
	queryStr := "select *, concat('" + API_PATH + "log/', id) as self"
	queryStr += " from log WHERE id=:id"
	nstmt, err := db.PrepareNamed(queryStr)
	err = nstmt.Select(&ret, arg)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	nstmt.Close()
	return ret, nil
}

// @Title getLogs
// @Description retrieves the log
// @Accept  application/json
// @Success 200 {array}    Log
// @Resource /api/2.0
// @Router /api/2.0/log [get]
func getLogs(db *sqlx.DB) (interface{}, error) {
	ret := []Log{}
	queryStr := "select *, concat('" + API_PATH + "log/', id) as self"
	queryStr += " from log"
	err := db.Select(&ret, queryStr)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return ret, nil
}

// @Title postLog
// @Description enter a new log
// @Accept  application/json
// @Param                 Body body     Log   true "Log object that should be added to the table"
// @Success 200 {object}    output_format.ApiWrapper
// @Resource /api/2.0
// @Router /api/2.0/log [post]
func postLog(payload []byte, db *sqlx.DB) (interface{}, error) {
	var v Log
	err := json.Unmarshal(payload, &v)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	sqlString := "INSERT INTO log("
	sqlString += "level"
	sqlString += ",message"
	sqlString += ",username"
	sqlString += ",ticketnum"
	sqlString += ",created_at"
	sqlString += ") VALUES ("
	sqlString += ":level"
	sqlString += ",:message"
	sqlString += ",:username"
	sqlString += ",:ticketnum"
	sqlString += ",:created_at"
	sqlString += ")"
	result, err := db.NamedExec(sqlString, v)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return result, err
}

// @Title putLog
// @Description modify an existing logentry
// @Accept  application/json
// @Param   id              path    int     true        "The row id"
// @Param                 Body body     Log   true "Log object that should be added to the table"
// @Success 200 {object}    output_format.ApiWrapper
// @Resource /api/2.0
// @Router /api/2.0/log/{id}  [put]
func putLog(id int64, payload []byte, db *sqlx.DB) (interface{}, error) {
	var arg Log
	err := json.Unmarshal(payload, &arg)
	arg.Id = id
	if err != nil {
		log.Println(err)
		return nil, err
	}
	sqlString := "UPDATE log SET "
	sqlString += "level = :level"
	sqlString += ",message = :message"
	sqlString += ",username = :username"
	sqlString += ",ticketnum = :ticketnum"
	sqlString += ",created_at = :created_at"
	sqlString += " WHERE id=:id"
	result, err := db.NamedExec(sqlString, arg)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return result, err
}

// @Title delLogById
// @Description deletes log information for a certain id
// @Accept  application/json
// @Param   id              path    int     false        "The row id"
// @Success 200 {array}    Log
// @Resource /api/2.0
// @Router /api/2.0/log/{id} [delete]
func delLog(id int64, db *sqlx.DB) (interface{}, error) {
	arg := Log{}
	arg.Id = id
	result, err := db.NamedExec("DELETE FROM log WHERE id=:id", arg)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return result, err
}