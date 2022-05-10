package types

type TickBiteMapData struct {

	Time int64   `bson:"time,omitempty" json:"time,omitempty"`
	Lng  float64 `bson:"lng" json:"lng"` 			//longitude
	Lat  float64 `bson:"lat" json:"lat"`            //latitude
	Type string  `bson:"type,omitempty" json:"type,omitempty"` //example: "TB"
	
}

type TimeSlider struct{

	MinLabel string `json:"minLabel"`
	MaxLabel string `json:"maxLabel"`
	Labels[] string `json:"labels"`

}


type ReportMapData struct{

	Slider	 TimeSlider 	`json:"slider"`
	Series[][] TickBiteMapData `json:"series"`
}
