package svc

import (
	"encoding/json"
	"fmt"
	"time"

	"amigo-api/app/device/model"
	"amigo-api/app/device/rpc/internal/config"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const (
	mqttMsgKey   = "mqtt:messages:latest"
	mqttMsgLimit = 100
	mqttTopic    = "devices/+/telemetry"
)

type ServiceContext struct {
	Config              config.Config
	RedisClient         *redis.Redis
	MQTTClient          mqtt.Client
	DeviceModel         model.DeviceModel
	AppModel            model.AppModel
	DeviceEventModel    model.DeviceEventModel
	FirmwareModel       model.FirmwareModel
	FirmwareTaskModel   model.FirmwareTaskModel
	WorkOrderModel      model.WorkOrderModel
	WorkOrderReplyModel model.WorkOrderReplyModel
}

type MQTTMessage struct {
	Topic     string `json:"topic"`
	Payload   string `json:"payload"`
	Timestamp int64  `json:"timestamp"`
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.DB.DataSource)
	rdb := redis.New(c.Redis.Host, func(r *redis.Redis) {
		r.Type = c.Redis.Type
		r.Pass = c.Redis.Pass
	})

	svcCtx := &ServiceContext{
		Config:              c,
		RedisClient:         rdb,
		DeviceModel:         model.NewDeviceModel(sqlConn, c.Cache),
		AppModel:            model.NewAppModel(sqlConn, c.Cache),
		DeviceEventModel:    model.NewDeviceEventModel(sqlConn, c.Cache),
		FirmwareModel:       model.NewFirmwareModel(sqlConn, c.Cache),
		FirmwareTaskModel:   model.NewFirmwareTaskModel(sqlConn, c.Cache),
		WorkOrderModel:      model.NewWorkOrderModel(sqlConn, c.Cache),
		WorkOrderReplyModel: model.NewWorkOrderReplyModel(sqlConn, c.Cache),
	}

	svcCtx.MQTTClient = newMQTTClient(c.EMQX, rdb)

	return svcCtx
}

func (m *MQTTMessage) String() string {
	return fmt.Sprintf("[%s] %s -> %s", time.UnixMilli(m.Timestamp).Format("15:04:05"), m.Topic, m.Payload)
}

func newMQTTClient(cfg config.EMQXConf, rdb *redis.Redis) mqtt.Client {
	opts := mqtt.NewClientOptions().
		AddBroker(cfg.Broker).
		SetClientID(cfg.ClientIdPrefix + "-rpc").
		SetKeepAlive(time.Duration(cfg.KeepAlive) * time.Second).
		SetConnectTimeout(time.Duration(cfg.ConnectTimeout) * time.Second).
		SetAutoReconnect(cfg.AutoReconnect).
		SetCleanSession(cfg.CleanSession).
		SetOnConnectHandler(func(client mqtt.Client) {
			logx.Infof("[MQTT] 已连接到 EMQX: %s", cfg.Broker)
			subscribeMQTT(client, rdb)
		}).
		SetConnectionLostHandler(func(client mqtt.Client, err error) {
			logx.Errorf("[MQTT] 连接断开: %v", err)
		})

	if cfg.Username != "" {
		opts.SetUsername(cfg.Username)
	}
	if cfg.Password != "" {
		opts.SetPassword(cfg.Password)
	}

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		logx.Errorf("[MQTT] 连接 EMQX 失败: %v", token.Error())
	}

	return client
}

func subscribeMQTT(client mqtt.Client, rdb *redis.Redis) {
	token := client.Subscribe(mqttTopic, 0, func(c mqtt.Client, msg mqtt.Message) {
		m := MQTTMessage{
			Topic:     msg.Topic(),
			Payload:   string(msg.Payload()),
			Timestamp: time.Now().UnixMilli(),
		}

		data, _ := json.Marshal(m)
		if _, err := rdb.Lpush(mqttMsgKey, string(data)); err != nil {
			logx.Errorf("[MQTT] 存储消息失败: %v", err)
			return
		}

		// 只保留最近 100 条
		if err := rdb.Ltrim(mqttMsgKey, 0, mqttMsgLimit-1); err != nil {
			logx.Errorf("[MQTT] 裁剪消息列表失败: %v", err)
		}

		logx.Infof("[MQTT] 收到: %s", m.String())
	})

	if token.Wait() && token.Error() != nil {
		logx.Errorf("[MQTT] 订阅 %s 失败: %v", mqttTopic, token.Error())
	} else {
		logx.Infof("[MQTT] 已订阅: %s", mqttTopic)
	}
}
