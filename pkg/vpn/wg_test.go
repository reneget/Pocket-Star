package vpn

import (
	"strings"
	"testing"
)

func TestGenerateKeyPair(t *testing.T) {
	kp, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair failed: %v", err)
	}
	if kp.Private == "" {
		t.Error("private key is empty")
	}
	if kp.Public == "" {
		t.Error("public key is empty")
	}
	if kp.Private == kp.Public {
		t.Error("private and public keys should differ")
	}

	// Keys should be valid base64
	if len(kp.Private) != 44 { // 32 bytes -> base64 = 44 chars
		t.Errorf("unexpected private key length: %d", len(kp.Private))
	}
	if len(kp.Public) != 44 {
		t.Errorf("unexpected public key length: %d", len(kp.Public))
	}
}

func TestGenerateHubConfig(t *testing.T) {
	hubCfg, nodeCfgs, err := GenerateHubConfig("10.0.0.1/24", 51820)
	if err != nil {
		t.Fatalf("GenerateHubConfig failed: %v", err)
	}

	if !strings.Contains(hubCfg, "[Interface]") {
		t.Error("hub config missing [Interface]")
	}
	if !strings.Contains(hubCfg, "ListenPort = 51820") {
		t.Error("hub config missing ListenPort")
	}

	expectedNodes := []string{"server", "work-laptop"}
	for _, name := range expectedNodes {
		cfg, ok := nodeCfgs[name]
		if !ok {
			t.Errorf("missing node config: %s", name)
			continue
		}
		if !strings.Contains(cfg, "[Interface]") {
			t.Errorf("node %s config missing [Interface]", name)
		}
		if !strings.Contains(cfg, "CHANGE_TO_HUB_IP") {
			t.Errorf("node %s config missing hub IP placeholder", name)
		}
	}
}

func TestGenerateNodeConfig(t *testing.T) {
	cfg := GenerateNodeConfig("test-node", "1.2.3.4:51820", "10.0.0.0/24")
	if !strings.Contains(cfg, "1.2.3.4:51820") {
		t.Error("node config missing hub endpoint")
	}
	if !strings.Contains(cfg, "<HUB_PUBLIC_KEY>") {
		t.Error("node config missing hub public key placeholder")
	}
}
