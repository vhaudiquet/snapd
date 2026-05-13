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
	"errors"
	"fmt"
	"strings"
)

// Error types for emulation
var (
	// ErrEmulatorNotFound indicates the emulator binary was not found
	ErrEmulatorNotFound = errors.New("emulator not found")
	// ErrEmulatorNotSupported indicates emulation is not supported for this architecture
	ErrEmulatorNotSupported = errors.New("emulation not supported for this architecture")
	// ErrEmulationDisabled indicates emulation is disabled
	ErrEmulationDisabled = errors.New("emulation is disabled")
)

// EmulatorNotFoundError indicates a specific emulator was not found
type EmulatorNotFoundError struct {
	Emulator EmulatorType
	Reason   string
}

func (e *EmulatorNotFoundError) Error() string {
	if e.Reason != "" {
		return fmt.Sprintf("emulator %s not found: %s", e.Emulator, e.Reason)
	}
	return fmt.Sprintf("emulator %s not found", e.Emulator)
}

func (e *EmulatorNotFoundError) Unwrap() error {
	return ErrEmulatorNotFound
}

// EmulationNotSupportedError indicates the architecture combination is not supported
type EmulationNotSupportedError struct {
	SourceArch string
	TargetArch string
	Available  []string
}

func (e *EmulationNotSupportedError) Error() string {
	if len(e.Available) > 0 {
		return fmt.Sprintf("emulation from %s to %s is not supported (available: %s)",
			e.SourceArch, e.TargetArch, strings.Join(e.Available, ", "))
	}
	return fmt.Sprintf("emulation from %s to %s is not supported", e.SourceArch, e.TargetArch)
}

func (e *EmulationNotSupportedError) Unwrap() error {
	return ErrEmulatorNotSupported
}

// SnapNeedsEmulationError indicates a snap requires emulation to run
type SnapNeedsEmulationError struct {
	Snap        string
	SourceArchs []string
	TargetArch  string
}

func (e *SnapNeedsEmulationError) Error() string {
	return fmt.Sprintf("snap %q is built for %s and requires emulation to run on %s (use --emulate flag)",
		e.Snap, strings.Join(e.SourceArchs, "/"), e.TargetArch)
}

// IsEmulationError checks if an error is related to emulation
func IsEmulationError(err error) bool {
	var emulatorNotFound *EmulatorNotFoundError
	var emulationNotSupported *EmulationNotSupportedError
	var snapNeedsEmulation *SnapNeedsEmulationError
	return errors.As(err, &emulatorNotFound) ||
		errors.As(err, &emulationNotSupported) ||
		errors.As(err, &snapNeedsEmulation) ||
		errors.Is(err, ErrEmulatorNotFound) ||
		errors.Is(err, ErrEmulatorNotSupported)
}
