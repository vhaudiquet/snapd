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
	"fmt"
	"os/exec"
	"sync"
)

// EmulatorType identifies the emulator to use
type EmulatorType string

const (
	// EmulatorBox64 is the box64 x86_64 emulator for arm64/riscv64
	EmulatorBox64 EmulatorType = "box64"
)

// Config holds emulation configuration for a snap
type Config struct {
	// Enabled indicates if emulation is active for this snap
	Enabled bool `json:"enabled"`
	// Emulator is the type of emulator to use
	Emulator EmulatorType `json:"emulator"`
	// SourceArch is the architecture of the snap binaries (e.g., "amd64")
	SourceArch string `json:"source_arch"`
	// TargetArch is the system architecture (e.g., "arm64")
	TargetArch string `json:"target_arch"`
	// EmulatorPath is the path to the emulator binary
	EmulatorPath string `json:"emulator_path"`
}

// EmulatorInfo contains information about an available emulator
type EmulatorInfo struct {
	// Type is the emulator type
	Type EmulatorType
	// Path is the path to the emulator binary
	Path string
	// SourceArchs lists architectures this emulator can run
	SourceArchs []string
	// TargetArchs lists architectures this emulator can run on
	TargetArchs []string
	// Flags are command-line flags to pass to the emulator
	Flags []string
	// Env contains environment variables to set for emulated execution
	Env map[string]string
}

// Registry maintains available emulators
type Registry struct {
	mu        sync.RWMutex
	available map[EmulatorType]*EmulatorInfo
}

var (
	globalRegistry     *Registry
	globalRegistryOnce sync.Once
)

// GetRegistry returns the global emulator registry
func GetRegistry() *Registry {
	globalRegistryOnce.Do(func() {
		globalRegistry = &Registry{
			available: make(map[EmulatorType]*EmulatorInfo),
		}
		globalRegistry.detectEmulators()
	})
	return globalRegistry
}

// detectEmulators scans the system for available emulators
func (r *Registry) detectEmulators() {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Detect box64
	if info, err := detectBox64(); err == nil {
		r.available[EmulatorBox64] = info
	}
}

// Available returns the list of available emulators
func (r *Registry) Available() []*EmulatorInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*EmulatorInfo, 0, len(r.available))
	for _, info := range r.available {
		result = append(result, info)
	}
	return result
}

// Get returns information about a specific emulator
func (r *Registry) Get(t EmulatorType) (*EmulatorInfo, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	info, ok := r.available[t]
	return info, ok
}

// CanEmulate checks if emulation is possible between architectures
func (r *Registry) CanEmulate(sourceArch, targetArch string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, info := range r.available {
		if contains(info.SourceArchs, sourceArch) && contains(info.TargetArchs, targetArch) {
			return true
		}
	}
	return false
}

// GetEmulatorFor returns an emulator that can handle the given architecture translation
func (r *Registry) GetEmulatorFor(sourceArch, targetArch string) (*EmulatorInfo, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, info := range r.available {
		if contains(info.SourceArchs, sourceArch) && contains(info.TargetArchs, targetArch) {
			return info, true
		}
	}
	return nil, false
}

// IsEmulationSupported checks if emulation is possible between architectures
// using the global registry
func IsEmulationSupported(sourceArch, targetArch string) bool {
	return GetRegistry().CanEmulate(sourceArch, targetArch)
}

// GetEmulatorConfig returns a Config for emulating sourceArch on targetArch
func GetEmulatorConfig(sourceArch, targetArch string) (*Config, error) {
	registry := GetRegistry()
	info, ok := registry.GetEmulatorFor(sourceArch, targetArch)
	if !ok {
		return nil, fmt.Errorf("no emulator available for %s → %s", sourceArch, targetArch)
	}

	return &Config{
		Enabled:      true,
		Emulator:     info.Type,
		SourceArch:   sourceArch,
		TargetArch:   targetArch,
		EmulatorPath: info.Path,
	}, nil
}

// EmulatedArchitectures returns the list of architectures that can be emulated
// on the given target architecture
func EmulatedArchitectures(targetArch string) []string {
	registry := GetRegistry()
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	archSet := make(map[string]bool)
	for _, info := range registry.available {
		if contains(info.TargetArchs, targetArch) {
			for _, arch := range info.SourceArchs {
				archSet[arch] = true
			}
		}
	}

	result := make([]string, 0, len(archSet))
	for arch := range archSet {
		result = append(result, arch)
	}
	return result
}

// findExecutable searches for an executable in PATH
func findExecutable(name string) (string, error) {
	path, err := exec.LookPath(name)
	if err != nil {
		return "", fmt.Errorf("cannot find %s: %v", name, err)
	}
	return path, nil
}

// contains checks if a string is in a slice
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
