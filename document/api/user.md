### 1. N/A

1. route definition

- Url: /api/admin/login
- Method: POST
- Request: `AdminLoginReq`
- Response: `AdminLoginResp`

2. request definition



```golang
type AdminLoginReq struct {
	Mobile string `json:"mobile"`
	Password string `json:"password"`
}
```


3. response definition



```golang
type AdminLoginResp struct {
	AdminId uint64 `json:"admin_id"`
	Mobile string `json:"mobile"`
	Username string `json:"username"`
	Avatar string `json:"avatar"`
	CreateTime int64 `json:"create_time"`
	Token string `json:"token"`
}
```

### 2. N/A

1. route definition

- Url: /api/admin/add
- Method: POST
- Request: `AdminAddReq`
- Response: `GetAdminResp`

2. request definition



```golang
type AdminAddReq struct {
	Mobile string `json:"mobile"`
	Username string `json:"username"`
	Avatar string `json:"avatar"`
	Password string `json:"password"`
	RePassword string `json:"re_password"`
}
```


3. response definition



```golang
type GetAdminResp struct {
	AdminId uint64 `json:"admin_id"`
	Mobile string `json:"mobile"`
	Username string `json:"username"`
	Avatar string `json:"avatar"`
	CreateTime int64 `json:"create_time"`
}
```

### 3. N/A

1. route definition

- Url: /api/admin/delete
- Method: POST
- Request: `AdminDeleteReq`
- Response: `EmptyResp`

2. request definition



```golang
type AdminDeleteReq struct {
	AdminId uint64 `json:"admin_id"`
}
```


3. response definition



```golang
type EmptyResp struct {
}
```

### 4. N/A

1. route definition

- Url: /api/admin/get
- Method: GET
- Request: `-`
- Response: `GetAdminResp`

2. request definition



3. response definition



```golang
type GetAdminResp struct {
	AdminId uint64 `json:"admin_id"`
	Mobile string `json:"mobile"`
	Username string `json:"username"`
	Avatar string `json:"avatar"`
	CreateTime int64 `json:"create_time"`
}
```

### 5. N/A

1. route definition

- Url: /api/admin/update
- Method: POST
- Request: `AdminUpdateReq`
- Response: `EmptyResp`

2. request definition



```golang
type AdminUpdateReq struct {
	AdminId uint64 `json:"admin_id"`
	Username string `json:"username"`
	Avatar string `json:"avatar"`
}
```


3. response definition



```golang
type EmptyResp struct {
}
```

### 6. N/A

1. route definition

- Url: /api/admin/update_mobile
- Method: POST
- Request: `AdminUpdateMobileReq`
- Response: `EmptyResp`

2. request definition



```golang
type AdminUpdateMobileReq struct {
	AdminId uint64 `json:"admin_id"`
	Mobile string `json:"mobile"`
}
```


3. response definition



```golang
type EmptyResp struct {
}
```

### 7. N/A

1. route definition

- Url: /api/admin/update_password
- Method: POST
- Request: `AdminUpdatePasswordReq`
- Response: `EmptyResp`

2. request definition



```golang
type AdminUpdatePasswordReq struct {
	AdminId uint64 `json:"admin_id"`
	Password string `json:"password"`
	RePassword string `json:"re_password"`
}
```


3. response definition



```golang
type EmptyResp struct {
}
```

### 8. N/A

1. route definition

- Url: /api/user/login
- Method: POST
- Request: `UserLoginReq`
- Response: `UserLoginResp`

2. request definition



```golang
type UserLoginReq struct {
	Mobile string `json:"mobile,optional"`
	Username string `json:"username,optional"`
	Password string `json:"password"`
}
```


3. response definition



```golang
type UserLoginResp struct {
	UserId uint64 `json:"user_id"`
	Mobile string `json:"mobile"`
	Username string `json:"username"`
	Avatar string `json:"avatar"`
	CreateTime int64 `json:"create_time"`
	Token string `json:"token"`
}
```

### 9. N/A

1. route definition

- Url: /api/user/third_login
- Method: POST
- Request: `UserThirdLoginReq`
- Response: `UserLoginResp`

2. request definition



```golang
type UserThirdLoginReq struct {
	Code string `json:"code"`
	AppType string `json:"app_type"`
}
```


3. response definition



```golang
type UserLoginResp struct {
	UserId uint64 `json:"user_id"`
	Mobile string `json:"mobile"`
	Username string `json:"username"`
	Avatar string `json:"avatar"`
	CreateTime int64 `json:"create_time"`
	Token string `json:"token"`
}
```

### 10. N/A

1. route definition

- Url: /api/user/add
- Method: POST
- Request: `UserAddReq`
- Response: `GetUserResp`

2. request definition



```golang
type UserAddReq struct {
	Mobile string `json:"mobile"`
	Username string `json:"username"`
	Avatar string `json:"avatar"`
	Password string `json:"password"`
	RePassword string `json:"re_password"`
}
```


3. response definition



```golang
type GetUserResp struct {
	UserId uint64 `json:"user_id"`
	Mobile string `json:"mobile"`
	Username string `json:"username"`
	Avatar string `json:"avatar"`
	CreateTime int64 `json:"create_time"`
}
```

### 11. N/A

1. route definition

- Url: /api/user/delete
- Method: POST
- Request: `UserDeleteReq`
- Response: `EmptyResp`

2. request definition



```golang
type UserDeleteReq struct {
	UserId uint64 `json:"user_id"`
}
```


3. response definition



```golang
type EmptyResp struct {
}
```

### 12. N/A

1. route definition

- Url: /api/user/get
- Method: GET
- Request: `-`
- Response: `GetUserResp`

2. request definition



3. response definition



```golang
type GetUserResp struct {
	UserId uint64 `json:"user_id"`
	Mobile string `json:"mobile"`
	Username string `json:"username"`
	Avatar string `json:"avatar"`
	CreateTime int64 `json:"create_time"`
}
```

### 13. N/A

1. route definition

- Url: /api/user/update
- Method: POST
- Request: `UserUpdateReq`
- Response: `EmptyResp`

2. request definition



```golang
type UserUpdateReq struct {
	UserId uint64 `json:"user_id"`
	Username string `json:"username"`
	Avatar string `json:"avatar"`
}
```


3. response definition



```golang
type EmptyResp struct {
}
```

### 14. N/A

1. route definition

- Url: /api/user/update_mobile
- Method: POST
- Request: `UserUpdateMobileReq`
- Response: `EmptyResp`

2. request definition



```golang
type UserUpdateMobileReq struct {
	UserId uint64 `json:"user_id"`
	Mobile string `json:"mobile"`
}
```


3. response definition



```golang
type EmptyResp struct {
}
```

### 15. N/A

1. route definition

- Url: /api/user/update_password
- Method: POST
- Request: `UserUpdatePasswordReq`
- Response: `EmptyResp`

2. request definition



```golang
type UserUpdatePasswordReq struct {
	UserId uint64 `json:"user_id"`
	Password string `json:"password"`
	RePassword string `json:"re_password"`
}
```


3. response definition



```golang
type EmptyResp struct {
}
```

