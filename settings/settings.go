/*
	Copyright (C) 2018 Nirmal Almara

    This file is part of Joyread.

    Joyread is free software: you can redistribute it and/or modify
    it under the terms of the GNU Affero General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    Joyread is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU Affero General Public License for more details.

    You should have received a copy of the GNU Affero General Public License
	along with Joyread.  If not, see <https://www.gnu.org/licenses/>.
*/

package settings

import (
	"fmt"
	"io/ioutil"

	cError "gitlab.com/joyread/ultimate/error"
	"gopkg.in/yaml.v2"
)

var conf BaseStruct

// BaseStruct struct
type BaseStruct struct {
	BaseValues BaseValuesStruct `yaml:"server" binding:"required"`
}

// BaseValuesStruct struct
type BaseValuesStruct struct {
	ServerPort string         `yaml:"port" binding:"required"`
	AssetPath  string         `yaml:"asset_path" binding:"required"`
	DataPath   string         `yaml:"data_path" binding:"required"`
	DBValues   DBValuesStruct `yaml:"database" binding:"required"`
}

// DBValuesStruct struct
type DBValuesStruct struct {
	DBHostname string `yaml:"hostname" binding:"required"`
	DBPort     string `yaml:"port" binding:"required"`
	DBName     string `yaml:"name" binding:"required"`
	DBUsername string `yaml:"username" binding:"required"`
	DBPassword string `yaml:"password" binding:"required"`
	DBSSLMode  string `yaml:"sslmode" binding:"required"`
}

func init() {
	fmt.Println("Running init ...")
	yamlFile, err := ioutil.ReadFile("config/app.yaml")
	cError.CheckError(err)

	err = yaml.Unmarshal(yamlFile, &conf)
	cError.CheckError(err)
}

// GetConf ...
func GetConf() *BaseStruct {
	return &conf
}
