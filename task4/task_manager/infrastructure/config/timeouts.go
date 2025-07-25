package config

import "time"

const (
	UserCreationTimeout time.Duration = 5 * time.Second
	UserLoginTimeout    time.Duration = 5 * time.Second
	DBConnectTimeOut    time.Duration = 10 * time.Second
	TokenTTL            time.Duration = 24 * time.Hour

	TaskCreationTimeout  time.Duration = 5 * time.Second
	TaskUpdateTimeout    time.Duration = 5 * time.Second
	TaskDeleteTimeout    time.Duration = 5 * time.Second
	GetTaskByIdTimeout   time.Duration = 5 * time.Second
	GetTaskByUserTimeout time.Duration = 8 * time.Second
)