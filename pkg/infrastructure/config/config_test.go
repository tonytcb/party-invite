package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	t.Parallel()

	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("error to load current directory: %v", err)
	}

	cfg, err := Load(currentDir + "/../../../")
	if err != nil {
		t.Fatalf("error to load config: %v", err)
	}

	assert.NotNil(t, cfg)
}

func TestConfig_IsValid(t *testing.T) {
	t.Parallel()

	var validConfig = &Config{
		AppName:        "test",
		HTTPPort:       1000,
		BaseLocation:   "dublin",
		LocationNearTo: 100,
	}

	type fields struct {
		Config *Config
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "should return no errors on a valid config",
			fields: fields{
				Config: validConfig,
			},
			wantErr: assert.NoError,
		},
		{
			name: "should error on missing APP_NAME env var",
			fields: fields{
				Config: func() *Config {
					c := *validConfig
					c.AppName = ""
					return &c
				}(),
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "undefined APP_NAME env var")
			},
		},
		{
			name: "should error on missing HTTP_PORT env var",
			fields: fields{
				Config: func() *Config {
					c := *validConfig
					c.HTTPPort = 0
					return &c
				}(),
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "undefined or invalid HTTP_PORT env var")
			},
		},

		{
			name: "should error on missing BASE_LOCATION env var",
			fields: fields{
				Config: func() *Config {
					c := *validConfig
					c.BaseLocation = "atlanta"
					return &c
				}(),
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "invalid BASE_LOCATION env var")
			},
		},
		{
			name: "should error on missing LOCATION_NEAR_TO env var",
			fields: fields{
				Config: func() *Config {
					c := *validConfig
					c.LocationNearTo = 0
					return &c
				}(),
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "undefined or invalid LOCATION_NEAR_TO env var")
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			got := tt.fields.Config.IsValid()

			tt.wantErr(t, got, "IsValid()")
		})
	}
}
