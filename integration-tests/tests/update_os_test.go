// -*- Mode: Go; indent-tabs-mode: t -*-
// +build !excludeintegration,!excludereboots

/*
 * Copyright (C) 2015, 2016 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package tests

import (
	//	"io/ioutil"
	//	"path"

	"github.com/ubuntu-core/snappy/snappy"

	"github.com/ubuntu-core/snappy/integration-tests/testutils/common"
	"github.com/ubuntu-core/snappy/integration-tests/testutils/partition"
	"github.com/ubuntu-core/snappy/integration-tests/testutils/updates"

	"gopkg.in/check.v1"
)

var _ = check.Suite(&updateOSSuite{})

type updateOSSuite struct {
	common.SnappySuite
}

func (s *updateOSSuite) assertBootDirContents(c *check.C) {
	system, err := partition.BootSystem()
	c.Assert(err, check.IsNil, check.Commentf("Error getting the boot system: %s", err))
	// TODO ask mvo about the new style boot dir.
	//	files, err := ioutil.ReadDir(
	//		path.Join(partition.BootDir(system), partition.OtherPartition(current)))
	//	c.Assert(err, check.IsNil, check.Commentf("Error reading the other partition boot dir: %s", err))

	// no filenames to check on amd64 the vmlinuz/initrd comes out of
	// the squashfs via grub loop mounts
	expectedFileNames := []string{}
	if system == "uboot" {
		expectedFileNames = []string{"dtbs", "initrd.img", "vmlinuz"}
	}

	fileNames := []string{}
	//	for _, f := range files {
	//	fileNames = append(fileNames, f.Name())
	//	}
	c.Assert(fileNames, check.DeepEquals, expectedFileNames,
		check.Commentf("Wrong files in the other partition boot dir"))
}

// Test that the ubuntu-core update to the same release and channel must install a newer
// version. If there is no update available, the channel version will be
// modified to fake an update. If there is a version available, the image will
// be up-to-date after running this test.
func (s *updateOSSuite) TestUpdateToSameReleaseAndChannel(c *check.C) {
	if common.BeforeReboot() {
		updateOutput := updates.CallFakeOSUpdate(c)
		expected := "(?ms)" +
			".*" +
			"^Reboot to use ubuntu-core version .*\\.\n"
		c.Assert(updateOutput, check.Matches, expected)
		s.assertBootDirContents(c)
		common.Reboot(c)
	} else if common.AfterReboot(c) {
		common.RemoveRebootMark(c)
		currentVersion := common.GetCurrentUbuntuCoreVersion(c)
		c.Assert(snappy.VersionCompare(currentVersion, common.GetSavedVersion(c)), check.Equals, 1,
			check.Commentf("Rebooted to the wrong version: %d", currentVersion))
	}
}
