package erlang

import (
	"testing"
)

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: Config{
				Cookie:         "super_secret",
				MyNodeName:     "testnode",
				RemoteNodeName: "remote@localhost",
				Creation:       1,
				Port:           9999,
			},
			wantErr: false,
		},
		{
			name: "empty cookie",
			config: Config{
				Cookie:         "",
				MyNodeName:     "testnode",
				RemoteNodeName: "remote@localhost",
				Creation:       1,
				Port:           9999,
			},
			wantErr: true,
		},
		{
			name: "empty node name",
			config: Config{
				Cookie:         "super_secret",
				MyNodeName:     "",
				RemoteNodeName: "remote@localhost",
				Creation:       1,
				Port:           9999,
			},
			wantErr: true,
		},
		{
			name: "empty remote node name",
			config: Config{
				Cookie:         "super_secret",
				MyNodeName:     "testnode",
				RemoteNodeName: "",
				Creation:       1,
				Port:           9999,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewConnection(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConnection() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigDefaultValues(t *testing.T) {
	config := Config{
		Cookie:         "super_secret",
		MyNodeName:     "testnode",
		RemoteNodeName: "remote@localhost",
		Creation:       1,
		Port:           9999,
	}

	// Test that the config struct can be created and accessed
	if config.Cookie != "super_secret" {
		t.Errorf("Expected cookie 'super_secret', got '%s'", config.Cookie)
	}

	if config.MyNodeName != "testnode" {
		t.Errorf("Expected node name 'testnode', got '%s'", config.MyNodeName)
	}

	if config.RemoteNodeName != "remote@localhost" {
		t.Errorf("Expected remote node name 'remote@localhost', got '%s'", config.RemoteNodeName)
	}

	if config.Creation != 1 {
		t.Errorf("Expected creation 1, got %d", config.Creation)
	}

	if config.Port != 9999 {
		t.Errorf("Expected port 9999, got %d", config.Port)
	}
}

// TestConnectionStruct tests that the Connection struct can be created
// This test doesn't actually connect to Erlang, just tests the struct
func TestConnectionStruct(t *testing.T) {
	// This test is mainly to ensure the struct can be created
	// without actually trying to connect to Erlang
	conn := &Connection{}
	if conn == nil {
		t.Error("Failed to create Connection struct")
	}
}