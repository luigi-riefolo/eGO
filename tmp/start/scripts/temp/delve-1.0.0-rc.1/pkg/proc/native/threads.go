package native

import (
	"fmt"

	"github.com/derekparker/delve/pkg/proc"
)

// Thread represents a single thread in the traced process
// ID represents the thread id or port, Process holds a reference to the
// Process struct that contains info on the process as
// a whole, and Status represents the last result of a `wait` call
// on this thread.
type Thread struct {
	ID                       int              // Thread ID or mach port
	Status                   *WaitStatus      // Status returned from last wait call
	CurrentBreakpoint        *proc.Breakpoint // Breakpoint thread is currently stopped at
	BreakpointConditionMet   bool             // Output of evaluating the breakpoint's condition
	BreakpointConditionError error            // Error evaluating the breakpoint's condition

	dbp            *Process
	singleStepping bool
	running        bool
	os             *OSSpecificDetails
}

// Continue the execution of this thread.
//
// If we are currently at a breakpoint, we'll clear it
// first and then resume execution. Thread will continue until
// it hits a breakpoint or is signaled.
func (thread *Thread) Continue() error {
	pc, err := thread.PC()
	if err != nil {
		return err
	}
	// Check whether we are stopped at a breakpoint, and
	// if so, single step over it before continuing.
	if _, ok := thread.dbp.FindBreakpoint(pc); ok {
		if err := thread.StepInstruction(); err != nil {
			return err
		}
	}
	return thread.resume()
}

// StepInstruction steps a single instruction.
//
// Executes exactly one instruction and then returns.
// If the thread is at a breakpoint, we first clear it,
// execute the instruction, and then replace the breakpoint.
// Otherwise we simply execute the next instruction.
func (thread *Thread) StepInstruction() (err error) {
	thread.running = true
	thread.singleStepping = true
	defer func() {
		thread.singleStepping = false
		thread.running = false
	}()
	pc, err := thread.PC()
	if err != nil {
		return err
	}

	bp, ok := thread.dbp.FindBreakpoint(pc)
	if ok {
		// Clear the breakpoint so that we can continue execution.
		_, err = thread.ClearBreakpoint(bp)
		if err != nil {
			return err
		}

		// Restore breakpoint now that we have passed it.
		defer func() {
			err = thread.dbp.writeSoftwareBreakpoint(thread, bp.Addr)
		}()
	}

	err = thread.singleStep()
	if err != nil {
		if _, exited := err.(proc.ProcessExitedError); exited {
			return err
		}
		return fmt.Errorf("step failed: %s", err.Error())
	}
	return nil
}

// Location returns the threads location, including the file:line
// of the corresponding source code, the function we're in
// and the current instruction address.
func (thread *Thread) Location() (*proc.Location, error) {
	pc, err := thread.PC()
	if err != nil {
		return nil, err
	}
	f, l, fn := thread.dbp.bi.PCToLine(pc)
	return &proc.Location{PC: pc, File: f, Line: l, Fn: fn}, nil
}

func (thread *Thread) Arch() proc.Arch {
	return thread.dbp.bi.Arch
}

func (thread *Thread) BinInfo() *proc.BinaryInfo {
	return &thread.dbp.bi
}

// SetPC sets the PC for this thread.
func (thread *Thread) SetPC(pc uint64) error {
	regs, err := thread.Registers(false)
	if err != nil {
		return err
	}
	return regs.SetPC(thread, pc)
}

// Stopped returns whether the thread is stopped at
// the operating system level. Actual implementation
// is OS dependant, look in OS thread file.
func (thread *Thread) Stopped() bool {
	return thread.stopped()
}

// Halt stops this thread from executing. Actual
// implementation is OS dependant. Look in OS
// thread file.
func (thread *Thread) Halt() (err error) {
	defer func() {
		if err == nil {
			thread.running = false
		}
	}()
	if thread.Stopped() {
		return
	}
	err = thread.halt()
	return
}

// SetCurrentBreakpoint sets the current breakpoint that this
// thread is stopped at as CurrentBreakpoint on the thread struct.
func (thread *Thread) SetCurrentBreakpoint() error {
	thread.CurrentBreakpoint = nil
	pc, err := thread.PC()
	if err != nil {
		return err
	}
	if bp, ok := thread.dbp.FindBreakpoint(pc); ok {
		thread.CurrentBreakpoint = bp
		if err = thread.SetPC(bp.Addr); err != nil {
			return err
		}
		thread.BreakpointConditionMet, thread.BreakpointConditionError = bp.CheckCondition(thread)
		if thread.CurrentBreakpoint != nil && thread.BreakpointConditionMet {
			if g, err := proc.GetG(thread); err == nil {
				thread.CurrentBreakpoint.HitCount[g.ID]++
			}
			thread.CurrentBreakpoint.TotalHitCount++
		}
	}
	return nil
}

func (thread *Thread) clearBreakpointState() {
	thread.CurrentBreakpoint = nil
	thread.BreakpointConditionMet = false
	thread.BreakpointConditionError = nil
}

func (th *Thread) Breakpoint() (*proc.Breakpoint, bool, error) {
	return th.CurrentBreakpoint, (th.CurrentBreakpoint != nil && th.BreakpointConditionMet), th.BreakpointConditionError
}

func (th *Thread) ThreadID() int {
	return th.ID
}

// ClearBreakpoint clears the specified breakpoint.
func (thread *Thread) ClearBreakpoint(bp *proc.Breakpoint) (*proc.Breakpoint, error) {
	if _, err := thread.WriteMemory(uintptr(bp.Addr), bp.OriginalData); err != nil {
		return nil, fmt.Errorf("could not clear breakpoint %s", err)
	}
	return bp, nil
}

// Registers obtains register values from the debugged process.
func (t *Thread) Registers(floatingPoint bool) (proc.Registers, error) {
	return registers(t, floatingPoint)
}

func (t *Thread) PC() (uint64, error) {
	regs, err := t.Registers(false)
	if err != nil {
		return 0, err
	}
	return regs.PC(), nil
}
