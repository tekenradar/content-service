package types

type MapData struct {

	Time int64   `bson:"time,omitempty"`
	Lng  float64 `bson:"lng"`            //longitude
	Lat  float64 `bson:"lat"`            //latitude
	Type string  `bson:"type,omitempty"` //example: "TB"
	
}
