package config

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"fmt"
)

type Database struct {
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
}

func (d Database) String(sslmode string) string {
	return fmt.Sprintf(""+
		"host=%s "+
		"port=%s "+
		"user=%s "+
		"dbname=%s "+
		"password=%s "+
		"sslmode=%v"+
		"", d.Host, d.Port, d.User, d.DBName, d.Password, sslmode)
}

type Config struct {
	DB Database

	Secret string
}

func GetConfig(cfg YamlConfig) Config {
	return Config{
		DB: Database{
			User:     cfg.Database.User,
			Password: cfg.Database.Password,
			Host:     cfg.Database.Host,
			Port:     cfg.Database.Port,
			DBName:   cfg.Database.DBName,
		},

		Secret: cfg.AppSecret,
	}
}
