// Copyright (C) 2020-2022 Red Hat, Inc.
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along
// with this program; if not, write to the Free Software Foundation, Inc.,
// 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.
package offlinecheck

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

type OfflineChecker struct{}

func LoadCatalogs() error {
	offlineDBPath := os.Getenv("TNF_OFFLINE_DB")
	if offlineDBPath == "" {
		return fmt.Errorf("no offline DB provided")
	}

	log.Infof("Offline DB location: %s", offlineDBPath)

	if err := loadContainersCatalog(offlineDBPath); err != nil {
		return fmt.Errorf("cannot load containers catalog, err: %v", err)
	}
	if err := loadHelmCatalog(offlineDBPath); err != nil {
		return fmt.Errorf("cannot load helm charts catalog, err: %v", err)
	}
	if err := loadOperatorsCatalog(offlineDBPath); err != nil {
		return fmt.Errorf("cannot load operators catalog, err: %v", err)
	}

	return nil
}

func (checker OfflineChecker) IsServiceReachable() bool {
	return true
}
