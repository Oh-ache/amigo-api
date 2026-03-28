package main

import (
	"fmt"
	"log"

	"amigo-api/common/utils/plug/miniapp"
)

func main() {
	// 示例：初始化小程序客户端
	testAppID := "test_app_id"
	testSecret := "test_secret"
	testToken := "test_token"
	testMchID := "test_mch_id"
	testKeyPath := "test_key_path"
	testCertPath := "test_cert_path"
	testKey := "test_key"

	// 创建配置
	cfg := &miniapp.MiniAppConfig{
		AppID:    testAppID,
		Secret:   testSecret,
		Token:    testToken,
		AESKey:   "abcdefghijklmnopqrstuvwxyz0123456789", // 有效的 base64 字符串
		MchID:    testMchID,
		KeyPath:  testKeyPath,
		CertPath: testCertPath,
		Key:      testKey,
	}

	// 初始化客户端
	client, err := miniapp.InitClient("test-miniapp", cfg)
	if err != nil {
		log.Fatalf("Failed to initialize miniapp client: %v", err)
	}

	fmt.Printf("Successfully initialized miniapp client: %v\n", client)

	// 测试获取客户端
	retrievedClient, err := miniapp.GetClient("test-miniapp")
	if err != nil {
		log.Fatalf("Failed to get miniapp client: %v", err)
	}

	fmt.Printf("Successfully retrieved miniapp client: %v\n", retrievedClient)

	// 测试列出所有客户端
	clients := miniapp.ListClients()
	fmt.Printf("Current miniapp clients: %v\n", clients)

	// 测试删除客户端
	miniapp.RemoveClient("test-miniapp")

	// 验证客户端是否已删除
	afterRemoveClient, err := miniapp.GetClient("test-miniapp")
	if err != nil {
		log.Fatalf("Failed to get miniapp client after removal: %v", err)
	}

	if afterRemoveClient == nil {
		fmt.Println("Successfully removed miniapp client")
	}

	// 再次列出所有客户端
	clientsAfterRemove := miniapp.ListClients()
	fmt.Printf("Miniapp clients after removal: %v\n", clientsAfterRemove)
}

