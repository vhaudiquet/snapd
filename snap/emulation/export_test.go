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
	"sync"
)

// MockBox64BinaryNames allows testing with custom binary names
func MockBox64BinaryNames(names []string) (restore func()) {
	old := box64BinaryNames
	box64BinaryNames = names
	return func() {
		box64BinaryNames = old
	}
}

// MockBox64SourceArchs allows testing with custom source architectures
func MockBox64SourceArchs(archs []string) (restore func()) {
	old := box64SourceArchs
	box64SourceArchs = archs
	return func() {
		box64SourceArchs = old
	}
}

// MockBox64TargetArchs allows testing with custom target architectures
func MockBox64TargetArchs(archs []string) (restore func()) {
	old := box64TargetArchs
	box64TargetArchs = archs
	return func() {
		box64TargetArchs = old
	}
}

// ResetRegistry resets the global registry to force re-detection
func ResetRegistry() {
	globalRegistry = nil
	globalRegistryOnce = sync.Once{}
}

// SetRegistry sets a mock registry for testing
func SetRegistry(r *Registry) (restore func()) {
	oldRegistry := globalRegistry
	globalRegistry = r
	return func() {
		globalRegistry = oldRegistry
	}
}

// NewTestRegistry creates a registry for testing with specified emulators
func NewTestRegistry(emulators map[EmulatorType]*EmulatorInfo) *Registry {
	return &Registry{
		available: emulators,
	}
}

// DetectBox64 exports detectBox64 for testing
var DetectBox64 = detectBox64

// FindBox64 exports findBox64 for testing
var FindBox64 = findBox64
