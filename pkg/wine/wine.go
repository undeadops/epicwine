package wine

// Request objects
type WineRequestInput struct {
	ID                  int
	Country             string `json:"country" csv:"country"`
	Description         string `json:"description" csv:"description"`
	Designation         string `json:"designation" csv:"designation"`
	Points              string `json:"points" csv:"points"`
	Price               string `json:"price" csv:"price"`
	Province            string `json:"province" csv:"province"`
	Region1             string `json:"region_1" csv:"region_1"`
	Region2             string `json:"region_2" csv:"region_2"`
	TasterName          string `json:"taster_name" csv:"taster_name"`
	TasterTwitterHandle string `json:"taster_twitter_handle" csv:"taster_twitter_handle"`
	Title               string `json:"title" csv:"title"`
	Variety             string `json:"variety" csv:"variety"`
	Winery              string `json:"winery" csv:"winery"`
}
