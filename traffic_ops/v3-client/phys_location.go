/*

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package client

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

const (
	API_PHYS_LOCATIONS = apiBase + "/phys_locations"
)

// CreatePhysLocation creates a PhysLocation.
func (to *Session) CreatePhysLocation(pl tc.PhysLocation) (tc.Alerts, ReqInf, error) {
	if pl.RegionID == 0 && pl.RegionName != "" {
		regions, _, err := to.GetRegionByNameWithHdr(pl.RegionName, nil)
		if err != nil {
			return tc.Alerts{}, ReqInf{}, err
		}
		if len(regions) == 0 {
			return tc.Alerts{}, ReqInf{}, errors.New("no region with name " + pl.RegionName)
		}
		pl.RegionID = regions[0].ID
	}
	var alerts tc.Alerts
	reqInf, err := to.post(API_PHYS_LOCATIONS, pl, nil, &alerts)
	return alerts, reqInf, err
}

func (to *Session) UpdatePhysLocationByIDWithHdr(id int, pl tc.PhysLocation, header http.Header) (tc.Alerts, ReqInf, error) {
	route := fmt.Sprintf("%s/%d", API_PHYS_LOCATIONS, id)
	var alerts tc.Alerts
	reqInf, err := to.put(route, pl, header, &alerts)
	return alerts, reqInf, err
}

// Update a PhysLocation by ID
// Deprecated: UpdatePhysLocationByID will be removed in 6.0. Use UpdatePhysLocationByIDWithHdr.
func (to *Session) UpdatePhysLocationByID(id int, pl tc.PhysLocation) (tc.Alerts, ReqInf, error) {
	return to.UpdatePhysLocationByIDWithHdr(id, pl, nil)
}

func (to *Session) GetPhysLocationsWithHdr(params map[string]string, header http.Header) ([]tc.PhysLocation, ReqInf, error) {
	path := API_PHYS_LOCATIONS + mapToQueryParameters(params)
	var data tc.PhysLocationsResponse
	reqInf, err := to.get(path, header, &data)
	return data.Response, reqInf, err
}

// Returns a list of PhysLocations with optional query parameters applied
// Deprecated: GetPhysLocations will be removed in 6.0. Use GetPhysLocationsWithHdr.
func (to *Session) GetPhysLocations(params map[string]string) ([]tc.PhysLocation, ReqInf, error) {
	return to.GetPhysLocationsWithHdr(params, nil)
}

func (to *Session) GetPhysLocationByIDWithHdr(id int, header http.Header) ([]tc.PhysLocation, ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", API_PHYS_LOCATIONS, id)
	var data tc.PhysLocationsResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// GET a PhysLocation by the PhysLocation ID
// Deprecated: GetPhysLocationByID will be removed in 6.0. Use GetPhysLocationByIDWithHdr.
func (to *Session) GetPhysLocationByID(id int) ([]tc.PhysLocation, ReqInf, error) {
	return to.GetPhysLocationByIDWithHdr(id, nil)
}

func (to *Session) GetPhysLocationByNameWithHdr(name string, header http.Header) ([]tc.PhysLocation, ReqInf, error) {
	route := fmt.Sprintf("%s?name=%s", API_PHYS_LOCATIONS, url.QueryEscape(name))
	var data tc.PhysLocationsResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// GET a PhysLocation by the PhysLocation name
// Deprecated: GetPhysLocationByName will be removed in 6.0. Use GetPhysLocationByNameWithHdr.
func (to *Session) GetPhysLocationByName(name string) ([]tc.PhysLocation, ReqInf, error) {
	return to.GetPhysLocationByNameWithHdr(name, nil)
}

// DELETE a PhysLocation by ID
func (to *Session) DeletePhysLocationByID(id int) (tc.Alerts, ReqInf, error) {
	route := fmt.Sprintf("%s/%d", API_PHYS_LOCATIONS, id)
	var alerts tc.Alerts
	reqInf, err := to.del(route, nil, &alerts)
	return alerts, reqInf, err
}
