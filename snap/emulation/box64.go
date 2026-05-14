// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2025 Canonical Ltd
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

package emulation

import (
	"os"
	"path/filepath"
	"runtime"
)

// box64BinaryNames are the possible names for the box64 binary
var box64BinaryNames = []string{"box64", "box64-stable", "box64-develop"}

// box64SourceArchs are the architectures box64 can emulate
var box64SourceArchs = []string{"amd64"}

// box64TargetArchs are the architectures box64 can run on
var box64TargetArchs = []string{"arm64", "riscv64"}

// Box64SnapName is the name of the box64 snap package
const Box64SnapName = "box64-with-gl4es" // "box64"

// Box64SnapPath is the path to box64 when installed as a snap
// This path is accessible from inside any snap's namespace because
// /snap is bind-mounted into every snap's mount namespace
const Box64SnapPath = "/snap/box64-with-gl4es/current/usr/bin/box64" //"/snap/box64/current/usr/bin/box64"

// detectBox64 checks if box64 is available on the system
func detectBox64() (*EmulatorInfo, error) {
	// First check if we're on a supported target architecture
	currentArch := runtime.GOARCH
	targetArch := goArchToDpkg(currentArch)
	if !contains(box64TargetArchs, targetArch) {
		return nil, ErrEmulatorNotSupported
	}

	// Look for box64 binary
	path, err := findBox64()
	if err != nil {
		return nil, err
	}

	return &EmulatorInfo{
		Type:        EmulatorBox64,
		Path:        path,
		SourceArchs: box64SourceArchs,
		TargetArchs: box64TargetArchs,
		Flags:       Box64Flags(),
		Env:         Box64Env(),
	}, nil
}

// findBox64 searches for the box64 binary in various locations
func findBox64() (string, error) {
	// First try PATH lookup
	for _, name := range box64BinaryNames {
		if path, err := findExecutable(name); err == nil {
			return path, nil
		}
	}

	// Try common installation locations
	commonPaths := []string{
		"/usr/bin/box64",
		"/usr/local/bin/box64",
		"/opt/box64/box64",
		"/snap/box64/current/usr/bin/box64",
	}

	for _, path := range commonPaths {
		if _, err := os.Stat(path); err == nil {
			// Check if executable
			if isExecutable(path) {
				return path, nil
			}
		}
	}

	// Check if installed as a snap
	if snapPath := findBox64Snap(); snapPath != "" {
		return snapPath, nil
	}

	return "", ErrEmulatorNotFound
}

// findBox64Snap checks if box64 is installed as a snap
func findBox64Snap() string {
	// Check common snap installation locations
	snapDirs := []string{
		"/snap/box64/current",
		"/var/lib/snapd/snap/box64/current",
	}

	for _, snapDir := range snapDirs {
		// Look for box64 binary in typical snap locations
		binPaths := []string{
			filepath.Join(snapDir, "usr/bin/box64"),
			filepath.Join(snapDir, "bin/box64"),
		}

		for _, binPath := range binPaths {
			if _, err := os.Stat(binPath); err == nil && isExecutable(binPath) {
				return binPath
			}
		}
	}

	return ""
}

// isExecutable checks if a file is executable
func isExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode()&0111 != 0
}

// goArchToDpkg converts Go architecture to dpkg architecture
func goArchToDpkg(goarch string) string {
	mapping := map[string]string{
		"amd64":   "amd64",
		"arm64":   "arm64",
		"arm":     "armhf",
		"riscv64": "riscv64",
		"ppc64le": "ppc64el",
		"s390x":   "s390x",
		"386":     "i386",
	}
	if dpkg, ok := mapping[goarch]; ok {
		return dpkg
	}
	return goarch
}

// Box64Flags returns additional flags for box64 based on the application
func Box64Flags() []string {
	// Default box64 flags that work well with most applications
	// These can be extended based on specific snap requirements
	return []string{}
}

// Box64Env returns environment variables for box64
func Box64Env() map[string]string {
	return map[string]string{
		// Enable dynamic recompilation optimizations
		"BOX64_DYNAREC": "1",
		// Allow using system libraries when available
		"BOX64_ALLOW_SYSTEM_LIBS": "1",
	}
}
