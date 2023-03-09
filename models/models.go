package models 

type Location struct {
    Key                string             `json:"Key"`
    LocalizedName      string             `json:"LocalizedName"`
    AdministrativeArea AdministrativeArea `json:"AdministrativeArea"`
}

type AdministrativeArea struct {
    ID            string `json:"ID"`
    LocalizedName string `json:"LocalizedName"`
}

type CurrentConditions struct {
    Temperature struct {
        Metric struct {
            Value float32 `json:"Value"`
        } `json:"Metric"`
        Imperial struct {
            Value float32 `json:"Value"`
        } `json:"Imperial"`
    } `json:"Temperature"`
    WeatherText string `json:"WeatherText"`
}