// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

// +build mage

package main

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/pkg/errors"

	"github.com/borischernov/beats/dev-tools/mage"
)

func init() {
	mage.BeatDescription = "Filebeat sends log files to Logstash or directly to Elasticsearch."
}

// Build builds the Beat binary.
func Build() error {
	return mage.Build(mage.DefaultBuildArgs())
}

// GolangCrossBuild build the Beat binary inside of the golang-builder.
// Do not use directly, use crossBuild instead.
func GolangCrossBuild() error {
	return mage.GolangCrossBuild(mage.DefaultGolangCrossBuildArgs())
}

// BuildGoDaemon builds the go-daemon binary (use crossBuildGoDaemon).
func BuildGoDaemon() error {
	return mage.BuildGoDaemon()
}

// CrossBuild cross-builds the beat for all target platforms.
func CrossBuild() error {
	return mage.CrossBuild()
}

// CrossBuildXPack cross-builds the beat with XPack for all target platforms.
func CrossBuildXPack() error {
	return mage.CrossBuildXPack()
}

// CrossBuildGoDaemon cross-builds the go-daemon binary using Docker.
func CrossBuildGoDaemon() error {
	return mage.CrossBuildGoDaemon()
}

// Clean cleans all generated files and build artifacts.
func Clean() error {
	return mage.Clean()
}

// Package packages the Beat for distribution.
// Use SNAPSHOT=true to build snapshots.
// Use PLATFORMS to control the target platforms.
func Package() {
	start := time.Now()
	defer func() { fmt.Println("package ran for", time.Since(start)) }()

	mage.UseElasticBeatPackaging()
	customizePackaging()

	mg.Deps(Update, prepareModulePackaging)
	mg.Deps(CrossBuild, CrossBuildXPack, CrossBuildGoDaemon)
	mg.SerialDeps(mage.Package, TestPackages)
}

// TestPackages tests the generated packages (i.e. file modes, owners, groups).
func TestPackages() error {
	return mage.TestPackages(mage.WithModules(), mage.WithModulesD())
}

// Update updates the generated files (aka make update).
func Update() error {
	return sh.Run("make", "update")
}

// Fields generates a fields.yml for the Beat.
func Fields() error {
	return mage.GenerateFieldsYAML("module")
}

// GoTestUnit executes the Go unit tests.
// Use TEST_COVERAGE=true to enable code coverage profiling.
// Use RACE_DETECTOR=true to enable the race detector.
func GoTestUnit(ctx context.Context) error {
	return mage.GoTest(ctx, mage.DefaultGoTestUnitArgs())
}

// GoTestIntegration executes the Go integration tests.
// Use TEST_COVERAGE=true to enable code coverage profiling.
// Use RACE_DETECTOR=true to enable the race detector.
func GoTestIntegration(ctx context.Context) error {
	return mage.GoTest(ctx, mage.DefaultGoTestIntegrationArgs())
}

// -----------------------------------------------------------------------------
// Customizations specific to Filebeat.
// - Include modules directory in packages (minus _meta and test files).
// - Include modules.d directory in packages.

var modulesDirGenerated = filepath.Clean("build/packaging/modules")

// customizePackaging modifies the package specs to add the modules and
// modules.d directory.
func customizePackaging() {
	var (
		moduleTarget = "module"
		module       = mage.PackageFile{
			Mode:   0644,
			Source: modulesDirGenerated,
		}

		modulesDTarget = "modules.d"
		modulesD       = mage.PackageFile{
			Mode:   0644,
			Source: "modules.d",
			Config: true,
		}
	)

	for _, args := range mage.Packages {
		pkgType := args.Types[0]
		switch pkgType {
		case mage.TarGz, mage.Zip:
			args.Spec.Files[moduleTarget] = module
			args.Spec.Files[modulesDTarget] = modulesD
		case mage.Deb, mage.RPM:
			args.Spec.Files["/usr/share/{{.BeatName}}/"+moduleTarget] = module
			args.Spec.Files["/etc/{{.BeatName}}/"+modulesDTarget] = modulesD
		case mage.DMG:
			args.Spec.Files["/Library/Application Support/{{.BeatVendor}}/{{.BeatName}}"+moduleTarget] = module
			args.Spec.Files["/etc/{{.BeatName}}/"+modulesDTarget] = modulesD
		default:
			panic(errors.Errorf("unhandled package type: %v", pkgType))
		}
	}
}

// prepareModulePackaging copies the module dir to the build dir and excludes
// _meta and test files so that they are not included in packages.
func prepareModulePackaging() error {
	if err := sh.Rm(modulesDirGenerated); err != nil {
		return err
	}

	copy := &mage.CopyTask{
		Source:  "module",
		Dest:    modulesDirGenerated,
		Mode:    0644,
		DirMode: 0755,
		Exclude: []string{
			"/_meta",
			"/test",
		},
	}
	return copy.Execute()
}
