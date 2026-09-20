### 1. N/A

1. route definition

- Url: /api/sdk/convert/png_to_elnk
- Method: POST
- Request: `PngToElnkReq`
- Response: `PngToElnkResp`

2. request definition



```golang
type PngToElnkReq struct {
	Url string `json:"url"`
}
```


3. response definition



```golang
type PngToElnkResp struct {
	Url string `json:"url"`
}
```

### 2. N/A

1. route definition

- Url: /api/sdk/ip/to_address
- Method: POST
- Request: `IpToAddressReq`
- Response: `IpToAddressResp`

2. request definition



```golang
type IpToAddressReq struct {
	Ip string `json:"ip"`
}
```


3. response definition



```golang
type IpToAddressResp struct {
	Country string `json:"country"`
	Province string `json:"province"`
	City string `json:"city"`
	Area string `json:"area"`
	Isp string `json:"isp"`
	Lng string `json:"lng"`
	Lat string `json:"lat"`
}
```

### 3. N/A

1. route definition

- Url: /api/sdk/message/check_code
- Method: POST
- Request: `CheckCodeReq`
- Response: `EmptyResp`

2. request definition



```golang
type CheckCodeReq struct {
	Platform string `json:"platform,optional,default=ali_sms"`
	SendType string `json:"send_type"`
	Mobile string `json:"mobile"`
	Code string `json:"code"`
}
```


3. response definition



```golang
type EmptyResp struct {
}
```

### 4. N/A

1. route definition

- Url: /api/sdk/message/send_code
- Method: POST
- Request: `SendCodeReq`
- Response: `EmptyResp`

2. request definition



```golang
type SendCodeReq struct {
	Platform string `json:"platform,optional,default=ali_sms"`
	SendType string `json:"send_type"`
	Mobile string `json:"mobile"`
}
```


3. response definition



```golang
type EmptyResp struct {
}
```

### 5. N/A

1. route definition

- Url: /api/sdk/oss/upload_file
- Method: POST
- Request: `UploadFileReq`
- Response: `UploadFileResp`

2. request definition



```golang
type UploadFileReq struct {
	FileName string `json:"file_name,optional" form:"file_name,optional"`
	File []byte `json:"file,optional" form:"file,optional"`
}
```


3. response definition



```golang
type UploadFileResp struct {
	Url string `json:"url"`
}
```

### 6. N/A

1. route definition

- Url: /api/sdk/oss/upload_token
- Method: POST
- Request: `UploadTokenReq`
- Response: `UploadTokenResp`

2. request definition



```golang
type UploadTokenReq struct {
	FileName string `json:"file_name"`
}
```


3. response definition



```golang
type UploadTokenResp struct {
	Token string `json:"token"`
	Expire int64 `json:"expire"`
}
```

### 7. N/A

1. route definition

- Url: /api/sdk/oss/upload_url
- Method: POST
- Request: `UploadUrlReq`
- Response: `UploadUrlResp`

2. request definition



```golang
type UploadUrlReq struct {
	FileName string `json:"file_name"`
	Url string `json:"url"`
}
```


3. response definition



```golang
type UploadUrlResp struct {
	Url string `json:"url"`
}
```

### 8. N/A

1. route definition

- Url: /api/sdk/weather/weather
- Method: POST
- Request: `WeatherReq`
- Response: `WeatherResp`

2. request definition



```golang
type WeatherReq struct {
	Code string `json:"code,optional,default="`
	Name string `json:"name,optional,default="`
	Type string `json:"type,optional,default=1"`
}
```


3. response definition



```golang
type WeatherResp struct {
	Date string `json:"date"`
	Week string `json:"week"`
	Weather string `json:"weather"`
	Temp string `json:"temp"`
	Wind string `json:"wind"`
	WeatherIcon string `json:"weather_icon"`
	Humidity string `json:"humidity"`
	Items []WeatherItem `json:"items"`
}
```


#### WeatherItem（嵌套类型，预报列表元素）

```golang
type WeatherItem struct {
    Date           string `json:"date"`
    Week           string `json:"week"`
    Weather        string `json:"weather"`
    Temp           string `json:"temp"`
    Wind           string `json:"wind"`
    WeatherIcon    string `json:"weather_icon"`
    NightWeather   string `json:"night_weather"`
    NightTemp      string `json:"night_temp"`
    NightWind      string `json:"night_wind"`
    DayWindPower   string `json:"day_wind_power"`
    NightWindPower string `json:"night_wind_power"`
}
```


#### 使用说明

- `code` 与 `name` 二选一（推荐 name，自动查 adcode，区级缺失时自动 fallback 到市级）。
- `type`：`1` = 实况天气，`2` = 未来 3 天预报，默认 `1`。
- 实况返回字段：`date / week / weather / temp / wind / humidity`，`items` 为空数组。
- 预报返回字段：`date / week / weather / temp / wind / humidity` 取自当前实况，`items` 填充 3 天预报。
- Redis 缓存 key：`cache:amigo:sdk:weather:{adcode}:{extensions}`，TTL 到当日 23:59:59。
- 依赖：
  - MySQL `base_code_item` 表（`sort_key=sdk, key=gaode.cityCode`）需有对应中文名记录
  - baseCode 字典 `gaode.weather.key`（高德开放平台申请的 key）

