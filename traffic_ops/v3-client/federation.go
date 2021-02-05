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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

const APIFederations = apiBase + "/federations"

func (to *Session) FederationsWithHdr(header http.Header) ([]tc.AllDeliveryServiceFederationsMapping, ReqInf, error) {
	type FederationResponse struct {
		Response []tc.AllDeliveryServiceFederationsMapping `json:"response"`
	}
	data := FederationResponse{}
	inf, err := to.get(APIFederations, header, &data)
	return data.Response, inf, err
}

// Deprecated: Federations will be removed in 6.0. Use FederationsWithHdr.
func (to *Session) Federations() ([]tc.AllDeliveryServiceFederationsMapping, ReqInf, error) {
	return to.FederationsWithHdr(nil)
}

func (to *Session) AllFederationsWithHdr(header http.Header) ([]tc.AllDeliveryServiceFederationsMapping, ReqInf, error) {
	type FederationResponse struct {
		Response []tc.AllDeliveryServiceFederationsMapping `json:"response"`
	}
	data := FederationResponse{}
	inf, err := to.get(apiBase+"/federations/all", header, &data)
	return data.Response, inf, err
}

// Deprecated: AllFederations will be removed in 6.0. Use AllFederationsWithHdr.
func (to *Session) AllFederations() ([]tc.AllDeliveryServiceFederationsMapping, ReqInf, error) {
	return to.AllFederationsWithHdr(nil)
}

func (to *Session) AllFederationsForCDNWithHdr(cdnName string, header http.Header) ([]tc.AllDeliveryServiceFederationsMapping, ReqInf, error) {
	// because the Federations JSON array is heterogeneous (array members may be a AllFederation or AllFederationCDN), we have to try decoding each separately.
	type FederationResponse struct {
		Response []json.RawMessage `json:"response"`
	}
	data := FederationResponse{}
	inf, err := to.get(apiBase+"/federations/all?cdnName="+url.QueryEscape(cdnName), header, &data)
	if err != nil {
		return nil, inf, err
	}

	feds := []tc.AllDeliveryServiceFederationsMapping{}
	for _, raw := range data.Response {
		fed := tc.AllDeliveryServiceFederationsMapping{}
		if err := json.Unmarshal([]byte(raw), &fed); err != nil {
			// we don't actually need the CDN, but we want to return an error if we got something unexpected
			cdnFed := tc.AllFederationCDN{}
			if err := json.Unmarshal([]byte(raw), &cdnFed); err != nil {
				return nil, inf, errors.New("Traffic Ops returned an unexpected object: '" + string(raw) + "'")
			}
		}
		feds = append(feds, fed)
	}
	return feds, inf, nil
}

// Deprecated: AllFederationsForCDN will be removed in 6.0. Use AllFederationsForCDNWithHdr.
func (to *Session) AllFederationsForCDN(cdnName string) ([]tc.AllDeliveryServiceFederationsMapping, ReqInf, error) {
	return to.AllFederationsForCDNWithHdr(cdnName, nil)
}

func (to *Session) CreateFederationDeliveryServices(federationID int, deliveryServiceIDs []int, replace bool) (ReqInf, error) {
	req := tc.FederationDSPost{DSIDs: deliveryServiceIDs, Replace: &replace}
	resp := map[string]interface{}{}
	inf, err := to.post(apiBase+`/federations/`+strconv.Itoa(federationID)+`/deliveryservices`, req, nil, &resp)
	return inf, err
}

func (to *Session) GetFederationDeliveryServicesWithHdr(federationID int, header http.Header) ([]tc.FederationDeliveryServiceNullable, ReqInf, error) {
	type FederationDSesResponse struct {
		Response []tc.FederationDeliveryServiceNullable `json:"response"`
	}
	data := FederationDSesResponse{}
	inf, err := to.get(fmt.Sprintf("%s/federations/%d/deliveryservices", apiBase, federationID), header, &data)
	return data.Response, inf, err
}

// GetFederationDeliveryServices Returns a given Federation's Delivery Services
// Deprecated: GetFederationDeliveryServices will be removed in 6.0. Use GetFederationDeliveryServicesWithHdr.
func (to *Session) GetFederationDeliveryServices(federationID int) ([]tc.FederationDeliveryServiceNullable, ReqInf, error) {
	return to.GetFederationDeliveryServicesWithHdr(federationID, nil)
}

// DeleteFederationDeliveryService Deletes a given Delivery Service from a Federation
func (to *Session) DeleteFederationDeliveryService(federationID, deliveryServiceID int) (tc.Alerts, ReqInf, error) {
	route := fmt.Sprintf("%s/federations/%d/deliveryservices/%d", apiBase, federationID, deliveryServiceID)
	var alerts tc.Alerts
	reqInf, err := to.del(route, nil, &alerts)
	return alerts, reqInf, err
}

// GetFederationUsers Associates the given Users' IDs to a Federation
func (to *Session) CreateFederationUsers(federationID int, userIDs []int, replace bool) (tc.Alerts, ReqInf, error) {
	req := tc.FederationUserPost{IDs: userIDs, Replace: &replace}
	var alerts tc.Alerts
	inf, err := to.post(fmt.Sprintf("%s/federations/%d/users", apiBase, federationID), req, nil, &alerts)
	return alerts, inf, err
}

func (to *Session) GetFederationUsersWithHdr(federationID int, header http.Header) ([]tc.FederationUser, ReqInf, error) {
	type FederationUsersResponse struct {
		Response []tc.FederationUser `json:"response"`
	}
	data := FederationUsersResponse{}
	inf, err := to.get(fmt.Sprintf("%s/federations/%d/users", apiBase, federationID), header, &data)
	return data.Response, inf, err
}

// GetFederationUsers Returns a given Federation's Users
// Deprecated: GetFederationUsers will be removed in 6.0. Use GetFederationUsersWithHdr.
func (to *Session) GetFederationUsers(federationID int) ([]tc.FederationUser, ReqInf, error) {
	return to.GetFederationUsersWithHdr(federationID, nil)
}

// DeleteFederationUser Deletes a given User from a Federation
func (to *Session) DeleteFederationUser(federationID, userID int) (tc.Alerts, ReqInf, error) {
	route := fmt.Sprintf("%s/federations/%d/users/%d", apiBase, federationID, userID)
	var alerts tc.Alerts
	reqInf, err := to.del(route, nil, &alerts)
	return alerts, reqInf, err
}

// AddFederationResolverMappingsForCurrentUser adds Federation Resolver mappings to one or more
// Delivery Services for the current user.
func (to *Session) AddFederationResolverMappingsForCurrentUser(mappings tc.DeliveryServiceFederationResolverMappingRequest) (tc.Alerts, ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.post(APIFederations, mappings, nil, &alerts)
	return alerts, reqInf, err
}

// DeleteFederationResolverMappingsForCurrentUser removes ALL Federation Resolver mappings for ALL
// Federations assigned to the currently authenticated user, as well as deleting ALL of the
// Federation Resolvers themselves.
func (to *Session) DeleteFederationResolverMappingsForCurrentUser() (tc.Alerts, ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.del(APIFederations, nil, &alerts)
	return alerts, reqInf, err
}

// ReplaceFederationResolverMappingsForCurrentUser replaces any and all Federation Resolver mappings
// on all Federations assigned to the currently authenticated user. This will first remove ALL
// Federation Resolver mappings for ALL Federations assigned to the currently authenticated user, as
// well as deleting ALL of the Federation Resolvers themselves. In other words, calling this is
// equivalent to a call to DeleteFederationResolverMappingsForCurrentUser followed by a call to
// AddFederationResolverMappingsForCurrentUser .
func (to *Session) ReplaceFederationResolverMappingsForCurrentUser(mappings tc.DeliveryServiceFederationResolverMappingRequest) (tc.Alerts, ReqInf, error) {
	var alerts tc.Alerts
	reqInf, err := to.put(APIFederations, mappings, nil, &alerts)
	return alerts, reqInf, err
}
