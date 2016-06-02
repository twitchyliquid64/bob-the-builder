// +build windows

package rpio

import (
  "errors"
)

type Direction uint8
type Pin uint8
type State uint8
type Pull uint8





// Set pin as Input
func (pin Pin) Input() {
}

// Set pin as Output
func (pin Pin) Output() {
}

// Set pin High
func (pin Pin) High() {
}

// Set pin Low
func (pin Pin) Low() {
}

// Toggle pin state
func (pin Pin) Toggle() {
}

// Set pin Direction
func (pin Pin) Mode(dir Direction) {
}

// Set pin state (high/low)
func (pin Pin) Write(state State) {
}

// Read pin state (high/low)
func (pin Pin) Read() State {
  return State(0)
}

// Set a given pull up/down mode
func (pin Pin) Pull(pull Pull) {
}

// Pull up pin
func (pin Pin) PullUp() {
}

// Pull down pin
func (pin Pin) PullDown() {
}

// Disable pullup/down on pin
func (pin Pin) PullOff() {
}

func Open() (err error) {
  return errors.New("Not supported on this platform")
}

func Close() error {
  return nil
}
