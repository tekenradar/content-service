package v1

import (
	"errors"
	"strconv"
	"strings"

	"github.com/influenzanet/study-service/pkg/studyengine"
	"github.com/influenzanet/study-service/pkg/types"

	cstypes "github.com/tekenradar/content-service/pkg/types"
)

func studyEventToTBMapData(event studyengine.ExternalEventPayload) (tickBiteMapData cstypes.TickBiteMapData, err error) {
	item, err := findResponseItem(event.Response.Responses, "TBLoc.Q2")
	if err != nil {
		return cstypes.TickBiteMapData{}, err
	}

	lat, err := parseResponseValueAsFloat(item, "lat")
	if err != nil {
		return cstypes.TickBiteMapData{}, err
	}

	lng, err := parseResponseValueAsFloat(item, "lng")
	if err != nil {
		return cstypes.TickBiteMapData{}, err
	}

	rtype, err := getReportType(event.Response.Key)
	if err != nil {
		return cstypes.TickBiteMapData{}, err
	}

	time := event.Response.SubmittedAt

	return cstypes.TickBiteMapData{
		Time: time,
		Lat:  lat,
		Lng:  lng,
		Type: rtype}, nil

}

func getReportType(key string) (Rtype string, err error) {
	if strings.Contains(key, "TB") {
		return "TB", nil
	} else if strings.Contains(key, "EM") {
		return "EM", nil
	} else if strings.Contains(key, "Fever") {
		return "FE", nil
	} else if strings.Contains(key, "LB") {
		return "Other", nil
	} else if strings.Contains(key, "Chronic") {
		return "Other", nil
	} else {
		return "", errors.New("Could not allocate type value")
	}
}

func parseResponseValueAsFloat(mapItem []types.ResponseItem, name string) (value float64, err error) {
	for _, mapItem := range mapItem {
		if mapItem.Key == name {
			val, err := strconv.ParseFloat(mapItem.Value, 64)
			if err != nil {
				return val, errors.New("Could not parse response value to float")
			}
			return val, nil
		}
	}
	return 0, errors.New("Could not find response value")
}

func findResponseItem(response []types.SurveyItemResponse, itemKey string) (item []types.ResponseItem, err error) {
	for _, resp := range response {

		if strings.Contains(resp.Key, itemKey) {

			for _, item := range resp.Response.Items {

				if item.Key == "map" {
					return item.Items, nil
				}

			}
		}
	}
	return []types.ResponseItem{}, errors.New("Could not find response item")
}
