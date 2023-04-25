// Package config
// Title       : config.go
// Author      : Tuffy  2023/4/4 15:58
// Description :
package config

type ApiConfig struct {
	RootURL       string `json:"root_url"`
	ListeningHost string `json:"listen_host"`
	ListeningPort int    `json:"listen_port"`
	Version       string `json:"version"`
}
