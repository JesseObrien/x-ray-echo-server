package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

type XraySegment xray.Segment

func (xs XraySegment) Display(isSubsegment bool, padding *string) {

	xs.DisplayInfo(isSubsegment, padding)

	p := *padding
	if isSubsegment {
		p = fmt.Sprintf("%s  ", *padding)
	}

	if len(xs.Subsegments) > 0 {
		xs.DisplaySubsegments(isSubsegment, &p)
	}
}

func (xs XraySegment) DisplayInfo(isSubsegment bool, padding *string) {
	title := "└─ "

	if isSubsegment {
		title = fmt.Sprintf("%s└─ subsegment: ", *padding)
	}

	subPadding := fmt.Sprintf("%s  ", *padding)

	elapsedTime := fmt.Sprintf("%fs", (xs.EndTime - xs.StartTime))

	metaData := log.Fields{
		"Elapsed Time": elapsedTime,
	}

	if len(xs.Metadata) > 0 {
		md := traverseMetadata(xs.Metadata, subPadding)
		for k, v := range md {
			metaData[k] = v
		}
	}

	log.WithFields(metaData).Infof("%s%s", title, xs.Name)

	if len(xs.AWS) > 0 {
		log.WithFields(getFields(xs.AWS)).Infof("%s└─ data", subPadding)
	}

	if xs.HTTP != nil {
		log.Infof("%s└─ http", subPadding)

		if request := xs.HTTP.GetRequest(); request != nil {
			log.WithFields(getFields(request)).Infof("%s  └─> Request", subPadding)
		}

		if response := xs.HTTP.GetResponse(); response != nil {
			log.WithFields(getFields(response)).Infof("%s  └─> Response", subPadding)
		}
	}
}

func traverseMetadata(metadata map[string]map[string]interface{}, padding string) logrus.Fields {
	md := logrus.Fields{}

	for _, v := range metadata {
		if len(v) > 0 {
			for key, val := range v {

				if m, ok := val.(map[string]interface{}); ok {
					for subKey, subVal := range getFields(m) {
						md[subKey] = subVal
					}
				} else {
					md[key] = val
				}
			}
		}
	}

	return md
}

func (xs XraySegment) DisplaySubsegments(isSubsegment bool, padding *string) {
	p := *padding

	if isSubsegment {
		p = fmt.Sprintf("%s", *padding)
	}

	for _, subSegment := range xs.Subsegments {
		subSeg := &XraySegment{}
		err := json.Unmarshal(subSegment, subSeg)

		if err != nil {
			log.WithError(err).Error("Could not unmarshal subsegment into a segment")
		}

		subSeg.Display(true, &p)
	}
}

func getFields(i interface{}) logrus.Fields {
	b, _ := json.Marshal(i)
	fields := logrus.Fields{}
	json.Unmarshal(b, &fields)
	return fields
}
