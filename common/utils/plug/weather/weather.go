package weather

import (
	"amigo-api/common/utils"

	"github.com/valyala/fasthttp"
)

type GaodeData struct {
	Status   string `json:"status"`
	Count    string `json:"count"`
	Info     string `json:"info"`
	Infocode string `json:"infocode"`
	Lives    []struct {
		Province         string `json:"province"`
		City             string `json:"city"`
		Adcode           string `json:"adcode"`
		Weather          string `json:"weather"`
		Temperature      string `json:"temperature"`
		Winddirection    string `json:"winddirection"`
		Windpower        string `json:"windpower"`
		Humidity         string `json:"humidity"`
		Reporttime       string `json:"reporttime"`
		TemperatureFloat string `json:"temperature_float"`
		HumidityFloat    string `json:"humidity_float"`
	} `json:"lives"`
	Forecasts []struct {
		City       string `json:"city"`
		Adcode     string `json:"adcode"`
		Province   string `json:"province"`
		Reporttime string `json:"reporttime"`
		Casts      []struct {
			Date           string `json:"date"`
			Week           string `json:"week"`
			Dayweather     string `json:"dayweather"`
			Nightweather   string `json:"nightweather"`
			Daytemp        string `json:"daytemp"`
			Nighttemp      string `json:"nighttemp"`
			Daywind        string `json:"daywind"`
			Nightwind      string `json:"nightwind"`
			Daypower       string `json:"daypower"`
			Nightpower     string `json:"nightpower"`
			DaytempFloat   string `json:"daytemp_float"`
			NighttempFloat string `json:"nighttemp_float"`
		} `json:"casts"`
	} `json:"forecasts"`
}

func GetWeather(key, code, extensions string) (resp *GaodeData, err error) {
	url := "https://restapi.amap.com/v3/weather/weatherInfo"
	params := map[string]string{
		"key":        key,
		"city":       code,
		"extensions": extensions,
	}

	resp = &GaodeData{}
	err = utils.FastWithDo(resp, fasthttp.MethodGet, url, params, nil, nil)
	return
}
