### 1. N/A

1. route definition

- Url: /api/device/get
- Method: GET
- Request: `GetDeviceReq`
- Response: `GetDeviceResp`

2. request definition



```golang
type GetDeviceReq struct {
	DeviceId uint64 `form:"device_id,default=0,optional"`
	MacAddress string `form:"mac_address,default=,optional"`
}
```


3. response definition



```golang
type GetDeviceResp struct {
	DeviceId uint64 `json:"device_id"`
	Name string `json:"name"`
	UserId uint64 `json:"user_id"`
	MacAddress string `json:"mac_address"`
	InternalIp string `json:"internal_ip"`
	BmpImage string `json:"bmp_image"`
	IsRunning int64 `json:"is_running"`
	ExtraData string `json:"extra_data"`
	IsDelete int64 `json:"is_delete"`
	CreateTime uint64 `json:"create_time"`
	UpdateTime uint64 `json:"update_time"`
}
```

### 2. N/A

1. route definition

- Url: /api/device/list
- Method: GET
- Request: `ListDeviceReq`
- Response: `ListDeviceResp`

2. request definition



```golang
type ListDeviceReq struct {
	Name string `form:"name,default=,optional"`
	UserId uint64 `form:"user_id,default=0,optional"`
	MacAddress string `form:"mac_address,default=,optional"`
	InternalIp string `form:"internal_ip,default=,optional"`
	IsRunning int64 `form:"is_running,default=0,optional"`
	IsDelete int64 `form:"is_delete,default=2,optional"`
	Page int64 `form:"page,default=1,optional"`
	PageSize int64 `form:"page_size,default=10,optional"`
}
```


3. response definition



```golang
type ListDeviceResp struct {
	List []GetDeviceResp `json:"list"`
	Total int64 `json:"total"`
}
```

### 3. N/A

1. route definition

- Url: /api/device/add
- Method: POST
- Request: `AddDeviceReq`
- Response: `EmptyResp`

2. request definition



```golang
type AddDeviceReq struct {
	Name string `json:"name"`
	UserId uint64 `json:"user_id"`
	MacAddress string `json:"mac_address"`
	InternalIp string `json:"internal_ip"`
	BmpImage string `json:"bmp_image,optional,default="`
	IsRunning int64 `json:"is_running,optional,default=2"`
	ExtraData string `json:"extra_data,optional,default="`
}
```


3. response definition



```golang
type EmptyResp struct {
}
```

### 4. N/A

1. route definition

- Url: /api/device/delete
- Method: POST
- Request: `DeleteDeviceReq`
- Response: `EmptyResp`

2. request definition



```golang
type DeleteDeviceReq struct {
	DeviceId uint64 `json:"device_id,default=0,optional"`
	MacAddress string `json:"mac_address,default=,optional"`
}
```


3. response definition



```golang
type EmptyResp struct {
}
```

### 5. N/A

1. route definition

- Url: /api/device/update
- Method: POST
- Request: `UpdateDeviceReq`
- Response: `EmptyResp`

2. request definition



```golang
type UpdateDeviceReq struct {
	DeviceId uint64 `json:"device_id"`
	Name string `json:"name"`
	UserId uint64 `json:"user_id"`
	MacAddress string `json:"mac_address"`
	InternalIp string `json:"internal_ip"`
	BmpImage string `json:"bmp_image,optional,default="`
	IsRunning int64 `json:"is_running,optional,default=2"`
	ExtraData string `json:"extra_data,optional,default="`
	IsDelete int64 `json:"is_delete"`
}
```


3. response definition



```golang
type EmptyResp struct {
}
```

