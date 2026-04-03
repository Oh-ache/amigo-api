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

