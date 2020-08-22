
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
	"log"
	"time"
)

type CachegroupsParameters struct {
	CreatedAt time.Time                  `db:"created_at" json:"createdAt"`
	Links     CachegroupsParametersLinks `json:"_links" db:-`
}

type CachegroupsParametersLinks struct {
	Self            string          `db:"self" json:"_self"`
	CachegroupsLink CachegroupsLink `json:"cachegroups" db:-`
	ParametersLink  ParametersLink  `json:"parameters" db:-`
}

// @Title getCachegroupsParametersById
// @Description retrieves the cachegroups_parameters information for a certain id
// @Accept  application/json
// @Param   id              path    int     false        "The row id"
// @Success 200 {array}    CachegroupsParameters
// @Resource /api/2.0
// @Router /api/2.0/cachegroups_parameters/{id} [get]
func getCachegroupsParameter(cachegroup string, parameterId int64, db *sqlx.DB) (interface{}, error) {
	ret := []CachegroupsParameters{}
	arg := CachegroupsParameters{}
	arg.Links.CachegroupsLink.ID = cachegroup
	arg.Links.ParametersLink.ID = parameterId
	queryStr := "select *, concat('" + API_PATH + "cachegroups_parameters', '/cachegroup/', cachegroup, '/parameter_id/', parameter_id) as self"
	queryStr += ", concat('" + API_PATH + "cachegroups/', cachegroup) as cachegroups_name_ref"
	queryStr += ", concat('" + API_PATH + "parameters/', parameter_id) as parameters_id_ref"
	queryStr += " from cachegroups_parameters WHERE cachegroup=:Links.CachegroupsLink.ID AND parameter_id=:Links.ParametersLink.ID"
	nstmt, err := db.PrepareNamed(queryStr)
	err = nstmt.Select(&ret, arg)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	nstmt.Close()
	return ret, nil
}

// @Title getCachegroupsParameterss
// @Description retrieves the cachegroups_parameters
// @Accept  application/json
// @Success 200 {array}    CachegroupsParameters
// @Resource /api/2.0
// @Router /api/2.0/cachegroups_parameters [get]
func getCachegroupsParameters(db *sqlx.DB) (interface{}, error) {
	ret := []CachegroupsParameters{}
	queryStr := "select *, concat('" + API_PATH + "cachegroups_parameters', '/cachegroup/', cachegroup, '/parameter_id/', parameter_id) as self"
	queryStr += ", concat('" + API_PATH + "cachegroups/', cachegroup) as cachegroups_name_ref"
	queryStr += ", concat('" + API_PATH + "parameters/', parameter_id) as parameters_id_ref"
	queryStr += " from cachegroups_parameters"
	err := db.Select(&ret, queryStr)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return ret, nil
}

// @Title postCachegroupsParameters
// @Description enter a new cachegroups_parameters
// @Accept  application/json
// @Param                 Body body     CachegroupsParameters   true "CachegroupsParameters object that should be added to the table"
// @Success 200 {object}    output_format.ApiWrapper
// @Resource /api/2.0
// @Router /api/2.0/cachegroups_parameters [post]
func postCachegroupsParameter(payload []byte, db *sqlx.DB) (interface{}, error) {
	var v CachegroupsParameters
	err := json.Unmarshal(payload, &v)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	sqlString := "INSERT INTO cachegroups_parameters("
	sqlString += "cachegroup"
	sqlString += ",parameter_id"
	sqlString += ",created_at"
	sqlString += ") VALUES ("
	sqlString += ":cachegroup"
	sqlString += ",:parameter_id"
	sqlString += ",:created_at"
	sqlString += ")"
	result, err := db.NamedExec(sqlString, v)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return result, err
}

// @Title putCachegroupsParameters
// @Description modify an existing cachegroups_parametersentry
// @Accept  application/json
// @Param   id              path    int     true        "The row id"
// @Param                 Body body     CachegroupsParameters   true "CachegroupsParameters object that should be added to the table"
// @Success 200 {object}    output_format.ApiWrapper
// @Resource /api/2.0
// @Router /api/2.0/cachegroups_parameters/{id}  [put]
func putCachegroupsParameter(cachegroup string, parameterId int64, payload []byte, db *sqlx.DB) (interface{}, error) {
	var arg CachegroupsParameters
	err := json.Unmarshal(payload, &arg)
	arg.Links.CachegroupsLink.ID = cachegroup
	arg.Links.ParametersLink.ID = parameterId
	if err != nil {
		log.Println(err)
		return nil, err
	}
	sqlString := "UPDATE cachegroups_parameters SET "
	sqlString += "cachegroup = :cachegroup"
	sqlString += ",parameter_id = :parameter_id"
	sqlString += ",created_at = :created_at"
	sqlString += " WHERE cachegroup=:Links.CachegroupsLink.ID AND parameter_id=:Links.ParametersLink.ID"
	result, err := db.NamedExec(sqlString, arg)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return result, err
}

// @Title delCachegroupsParametersById
// @Description deletes cachegroups_parameters information for a certain id
// @Accept  application/json
// @Param   id              path    int     false        "The row id"
// @Success 200 {array}    CachegroupsParameters
// @Resource /api/2.0
// @Router /api/2.0/cachegroups_parameters/{id} [delete]
func delCachegroupsParameter(cachegroup string, parameterId int64, db *sqlx.DB) (interface{}, error) {
	arg := CachegroupsParameters{}
	arg.Links.CachegroupsLink.ID = cachegroup
	arg.Links.ParametersLink.ID = parameterId
	result, err := db.NamedExec("DELETE FROM cachegroups_parameters WHERE cachegroup=:Links.CachegroupsLink.ID AND parameter_id=:Links.ParametersLink.ID", arg)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return result, err
}