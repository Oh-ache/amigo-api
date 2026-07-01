package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v3"
)

type yamlConfig struct {
	DB struct {
		DataSource string `yaml:"DataSource"`
	} `yaml:"DB"`
}

const defaultConfigPath = "app/device/rpc/etc/device.yaml"

func main() {
	var (
		dsnArg        string
		configArg     string
		createArg     bool
		seedArg       bool
		dropArg       bool
		seedCasbinArg bool
	)
	flag.StringVar(&dsnArg, "dsn", "", "MySQL DSN 连接串（优先级最高）")
	flag.StringVar(&configArg, "config", "", "yaml 配置文件路径")
	flag.BoolVar(&createArg, "create", false, "建表")
	flag.BoolVar(&seedArg, "seed", false, "写入种子数据")
	flag.BoolVar(&dropArg, "drop", false, "先删表再重建（危险，需二次确认）")
	flag.BoolVar(&seedCasbinArg, "seed-casbin", false, "写入 Casbin 初始角色和权限策略")
	flag.Parse()

	if !createArg && !seedArg && !dropArg && !seedCasbinArg {
		flag.Usage()
		return
	}

	dsn := resolveDSN(dsnArg, configArg)
	if dsn == "" {
		log.Fatal("未找到 DSN，请通过 --dsn 或 --config 指定")
	}
	// 不在日志中输出密码
	fmt.Printf("DSN: %s\n", maskPassword(dsn))

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	defer db.Close()
	db.SetMaxOpenConns(4)
	db.SetConnMaxLifetime(30 * time.Second)

	if err := db.Ping(); err != nil {
		log.Fatalf("数据库 Ping 失败: %v", err)
	}
	fmt.Println("✓ 数据库连接成功")

	if dropArg {
		if !confirm("确定要删除所有业务表（device/app/device_app/sensor）？此操作不可逆！") {
			fmt.Println("已取消")
			return
		}
		dropTables(db)
	}

	if createArg {
		createTables(db)
	}

	if seedArg {
		seedData(db)
	}

	if seedCasbinArg {
		seedCasbinData(db)
	}

	fmt.Println("\n✓ 全部完成")
}

// ============================ DSN ============================

func resolveDSN(dsnFlag, configFlag string) string {
	if dsnFlag != "" {
		return dsnFlag
	}
	if env := os.Getenv("MYSQL_DSN"); env != "" {
		return env
	}
	cfgPath := configFlag
	if cfgPath == "" {
		cfgPath = defaultConfigPath
		if _, err := os.Stat(cfgPath); err != nil {
			// 尝试 baseCode 配置
			cfgPath = "app/baseCode/rpc/etc/baseCode.yaml"
			if _, err := os.Stat(cfgPath); err != nil {
				return ""
			}
		}
	}
	fmt.Printf("读取配置: %s\n", cfgPath)
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}
	var cfg yamlConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}
	return cfg.DB.DataSource
}

func maskPassword(dsn string) string {
	if idx := strings.Index(dsn, "@"); idx > 0 {
		userPart := dsn[:idx]
		if colon := strings.LastIndex(userPart, ":"); colon > 0 {
			return dsn[:colon+1] + "***@" + dsn[idx+1:]
		}
	}
	return dsn
}

func confirm(msg string) bool {
	fmt.Printf("\n⚠ %s\n输入 yes 确认: ", msg)
	var input string
	fmt.Scanln(&input)
	return strings.ToLower(strings.TrimSpace(input)) == "yes"
}

// ============================ 建表 ============================

func createTables(db *sql.DB) {
	fmt.Println("\n--- 建表 ---")

	// 1. device 新增字段（幂等）
	execIgnoreDup(db, "ALTER TABLE `device` ADD COLUMN `firmware` varchar(100) NOT NULL DEFAULT '' COMMENT '固件版本' AFTER `internal_ip`")
	execIgnoreDup(db, "ALTER TABLE `device` ADD COLUMN `device_group` varchar(100) NOT NULL DEFAULT '' COMMENT '设备分组' AFTER `firmware`")

	// 2. app 表
	exec(db, `CREATE TABLE IF NOT EXISTS app (
  app_id bigint unsigned NOT NULL AUTO_INCREMENT,
  name varchar(255) NOT NULL COMMENT '应用名称',
  app_key varchar(100) NOT NULL COMMENT '应用唯一标识',
  category varchar(50) NOT NULL DEFAULT '' COMMENT '分类(sensor/automation/integration/analytics)',
  version varchar(50) NOT NULL DEFAULT '' COMMENT '版本号',
  vendor varchar(100) NOT NULL DEFAULT '' COMMENT '厂商',
  icon_bg varchar(255) NOT NULL DEFAULT '' COMMENT '图标背景色',
  icon_color varchar(50) NOT NULL DEFAULT '' COMMENT '图标颜色',
  icon_label varchar(10) NOT NULL DEFAULT '' COMMENT '图标文字',
  description text COMMENT '应用描述',
  sensor_type varchar(50) NOT NULL DEFAULT '' COMMENT '传感器类型',
  extra_data text COMMENT 'JSON扩展数据(如data_points等)',
  is_delete tinyint NOT NULL DEFAULT '2' COMMENT '1删除2未删除',
  create_time bigint unsigned NOT NULL DEFAULT '0',
  update_time bigint unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (app_id),
  UNIQUE KEY idx_app_key (app_key)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='应用表'`)

	// 3. device_app 关联表
	exec(db, `CREATE TABLE IF NOT EXISTS device_app (
  device_app_id bigint unsigned NOT NULL AUTO_INCREMENT,
  device_id bigint unsigned NOT NULL COMMENT '设备id',
  app_id bigint unsigned NOT NULL COMMENT '应用id',
  config_data text COMMENT '设备上该应用的配置JSON',
  is_installed tinyint NOT NULL DEFAULT '2' COMMENT '1已安装2未安装',
  install_time bigint unsigned NOT NULL DEFAULT '0' COMMENT '安装时间',
  is_delete tinyint NOT NULL DEFAULT '2' COMMENT '1删除2未删除',
  create_time bigint unsigned NOT NULL DEFAULT '0',
  update_time bigint unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (device_app_id),
  UNIQUE KEY idx_device_app (device_id, app_id),
  KEY idx_device_id (device_id),
  KEY idx_app_id (app_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='设备-应用关联表'`)

	// 4. sensor 表
	exec(db, `CREATE TABLE IF NOT EXISTS sensor (
  sensor_id bigint unsigned NOT NULL AUTO_INCREMENT,
  device_id bigint unsigned NOT NULL COMMENT '所属设备id',
  sensor_name varchar(100) NOT NULL COMMENT '传感器名称',
  sensor_key varchar(100) NOT NULL COMMENT '传感器唯一标识',
  sensor_type varchar(50) NOT NULL COMMENT '传感器类型(temperature/humidity/pressure/motion等)',
  unit varchar(20) NOT NULL DEFAULT '' COMMENT '单位(°C/%/hPa等)',
  params_json text NOT NULL COMMENT '自定义JSON参数(计算公式/校准参数/阈值等)',
  is_running tinyint NOT NULL DEFAULT '2' COMMENT '1运行中2未运行',
  extra_data text COMMENT '额外JSON扩展数据',
  is_delete tinyint NOT NULL DEFAULT '2' COMMENT '1删除2未删除',
  create_time bigint unsigned NOT NULL DEFAULT '0',
  update_time bigint unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (sensor_id),
  UNIQUE KEY idx_device_sensor (device_id, sensor_key),
  KEY idx_device_id (device_id),
  KEY idx_sensor_type (sensor_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='传感器配置表'`)

	fmt.Println("✓ 建表完成")
}

// ============================ 删表 ============================

func dropTables(db *sql.DB) {
	fmt.Println("\n--- 删表 ---")
	execIgnoreErr(db, "DROP TABLE IF EXISTS sensor")
	execIgnoreErr(db, "DROP TABLE IF EXISTS device_app")
	execIgnoreErr(db, "DROP TABLE IF EXISTS app")
	// 不删 device 表本身，只回退新增字段
	execIgnoreErr(db, "ALTER TABLE `device` DROP COLUMN `device_group`")
	execIgnoreErr(db, "ALTER TABLE `device` DROP COLUMN `firmware`")
	fmt.Println("✓ 已删除")
}

// ============================ 种子数据 ============================

type appSeed struct {
	appKey      string
	name        string
	category    string
	version     string
	vendor      string
	iconBg      string
	iconColor   string
	iconLabel   string
	description string
	sensorType  string
	dataPoints  []string
}

func seedData(db *sql.DB) {
	fmt.Println("\n--- 种子数据 ---")

	now := time.Now().Unix()
	apps := []appSeed{
		{
			appKey: "sensor-temperature", name: "温湿度传感器", category: "sensor",
			version: "v1.2.4", vendor: "Amigo Labs",
			iconBg: "linear-gradient(135deg, #fb923c 0%, #ef4444 100%)", iconColor: "#fb923c", iconLabel: "温",
			description: "采集环境温度与相对湿度，每 30 秒上报一次，支持阈值告警。",
			sensorType:  "temperature", dataPoints: []string{"temperature", "humidity"},
		},
		{
			appKey: "sensor-air", name: "空气质量传感器", category: "sensor",
			version: "v1.5.0", vendor: "Amigo Labs",
			iconBg: "linear-gradient(135deg, #10b981 0%, #059669 100%)", iconColor: "#10b981", iconLabel: "气",
			description: "实时监测 PM2.5、CO₂、TVOC 等空气指标，守护健康呼吸。",
			sensorType:  "air", dataPoints: []string{"pm25", "co2", "tvoc"},
		},
		{
			appKey: "sensor-motion", name: "人体红外感应", category: "sensor",
			version: "v1.0.3", vendor: "Amigo Labs",
			iconBg: "linear-gradient(135deg, #6366f1 0%, #4f46e5 100%)", iconColor: "#6366f1", iconLabel: "感",
			description: "基于 PIR + 微波双鉴，精准检测人体活动并触发联动。",
			sensorType:  "motion", dataPoints: []string{"triggered"},
		},
		{
			appKey: "sensor-attitude", name: "六轴姿态传感", category: "sensor",
			version: "v1.1.2", vendor: "Amigo Labs",
			iconBg: "linear-gradient(135deg, #0ea5e9 0%, #2563eb 100%)", iconColor: "#0ea5e9", iconLabel: "姿",
			description: "MPU6050 三轴加速度 + 三轴陀螺仪，可用于跌倒检测与姿态识别。",
			sensorType:  "attitude", dataPoints: []string{"pitch", "roll", "yaw"},
		},
		{
			appKey: "automation-rule", name: "自动化规则引擎", category: "automation",
			version: "v1.3.2", vendor: "Amigo Labs",
			iconBg: "linear-gradient(135deg, #f59e0b 0%, #d97706 100%)", iconColor: "#f59e0b", iconLabel: "规",
			description: "可视化 IF-THEN 规则编排，支持时间、设备状态、传感器多条件组合。",
		},
		{
			appKey: "automation-schedule", name: "定时任务调度", category: "automation",
			version: "v1.3.1", vendor: "Amigo Labs",
			iconBg: "linear-gradient(135deg, #a855f7 0%, #7c3aed 100%)", iconColor: "#a855f7", iconLabel: "时",
			description: "基于 cron 表达式的定时任务，支持设备批量联动与节假日例外。",
		},
		{
			appKey: "integration-webhook", name: "Webhook 推送", category: "integration",
			version: "v1.2.0", vendor: "Amigo Connect",
			iconBg: "linear-gradient(135deg, #2563eb 0%, #1d4ed8 100%)", iconColor: "#2563eb", iconLabel: "钩",
			description: "将设备事件实时推送到任意 HTTP(S) 端点，支持签名校验与重试。",
		},
		{
			appKey: "integration-mqtt", name: "MQTT 桥接", category: "integration",
			version: "v1.4.2", vendor: "Amigo Connect",
			iconBg: "linear-gradient(135deg, #14b8a6 0%, #0d9488 100%)", iconColor: "#14b8a6", iconLabel: "MQ",
			description: "桥接外部 MQTT Broker，可双向同步设备状态与下发指令。",
		},
	}

	for _, a := range apps {
		extra := map[string]interface{}{}
		if len(a.dataPoints) > 0 {
			extra["data_points"] = a.dataPoints
		}
		extraJSON := "{}"
		if b, err := json.Marshal(extra); err == nil && string(b) != "{}" {
			extraJSON = string(b)
		}

		_, err := db.Exec(
			`INSERT INTO app (name, app_key, category, version, vendor, icon_bg, icon_color, icon_label, description, sensor_type, extra_data, create_time, update_time)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			 ON DUPLICATE KEY UPDATE name=VALUES(name), category=VALUES(category), version=VALUES(version),
			 vendor=VALUES(vendor), icon_bg=VALUES(icon_bg), icon_color=VALUES(icon_color), icon_label=VALUES(icon_label),
			 description=VALUES(description), sensor_type=VALUES(sensor_type), extra_data=VALUES(extra_data), update_time=VALUES(update_time)`,
			a.name, a.appKey, a.category, a.version, a.vendor, a.iconBg, a.iconColor, a.iconLabel, a.description, a.sensorType, extraJSON, now, now,
		)
		if err != nil {
			log.Printf("插入种子数据失败 [%s]: %v", a.appKey, err)
		} else {
			fmt.Printf("  ✓ %s\n", a.appKey)
		}
	}

	fmt.Printf("✓ 种子数据完成 (%d 条)\n", len(apps))
}

// ============================ 辅助函数 ============================

func exec(db *sql.DB, sql string) {
	_, err := db.Exec(sql)
	if err != nil {
		log.Fatalf("SQL 执行失败: %v\n%s", err, sql)
	}
	// 提取第一行作为标识
	firstLine := strings.SplitN(strings.TrimSpace(sql), "\n", 2)[0]
	if len(firstLine) > 60 {
		firstLine = firstLine[:60] + "..."
	}
	fmt.Printf("  ✓ %s\n", firstLine)
}

func execIgnoreErr(db *sql.DB, sql string) {
	_, err := db.Exec(sql)
	if err != nil {
		fmt.Printf("  ⚠ 跳过（%v）: %s\n", err, sql[:min(len(sql), 60)])
	} else {
		fmt.Printf("  ✓ %s\n", sql[:min(len(sql), 60)])
	}
}

func execIgnoreDup(db *sql.DB, sql string) {
	_, err := db.Exec(sql)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate column") || strings.Contains(err.Error(), "Duplicate key") {
			fmt.Printf("  - 已存在: %s\n", sql[:min(len(sql), 60)])
			return
		}
		log.Fatalf("SQL 执行失败: %v\n%s", err, sql)
	}
	fmt.Printf("  ✓ %s\n", sql[:min(len(sql), 60)])
}

// ============================ Casbin 种子数据 ============================

func seedCasbinData(db *sql.DB) {
	fmt.Println("\n--- Casbin 种子数据 ---")

	// 先清空旧的种子数据，避免重复（幂等）
	execIgnoreErr(db, "DELETE FROM casbin_rule WHERE ptype='g' AND v0='1' AND v1='super_admin' AND v2='amigo-admin'")
	execIgnoreErr(db, "DELETE FROM casbin_rule WHERE ptype='p' AND v1='amigo-admin'")

	casbinSQLs := []string{
		// 角色绑定：admin_id=1 → super_admin
		"INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('g', '1', 'super_admin', 'amigo-admin')",
		// 超级管理员权限策略
		"INSERT INTO casbin_rule (ptype, v0, v1, v2, v3) VALUES ('p', 'super_admin', 'amigo-admin', '/api/admin/*', 'GET')",
		"INSERT INTO casbin_rule (ptype, v0, v1, v2, v3) VALUES ('p', 'super_admin', 'amigo-admin', '/api/admin/*', 'POST')",
		"INSERT INTO casbin_rule (ptype, v0, v1, v2, v3) VALUES ('p', 'super_admin', 'amigo-admin', '/api/admin/*', 'DELETE')",
		"INSERT INTO casbin_rule (ptype, v0, v1, v2, v3) VALUES ('p', 'super_admin', 'amigo-admin', '/api/device/*', 'GET')",
		"INSERT INTO casbin_rule (ptype, v0, v1, v2, v3) VALUES ('p', 'super_admin', 'amigo-admin', '/api/device/*', 'POST')",
		"INSERT INTO casbin_rule (ptype, v0, v1, v2, v3) VALUES ('p', 'super_admin', 'amigo-admin', '/api/base_code_*/*', 'GET')",
		"INSERT INTO casbin_rule (ptype, v0, v1, v2, v3) VALUES ('p', 'super_admin', 'amigo-admin', '/api/base_code_*/*', 'POST')",
	}

	for _, s := range casbinSQLs {
		exec(db, s)
	}

	fmt.Println("✓ Casbin 种子数据完成")
}
