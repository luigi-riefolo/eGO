package terminal

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/derekparker/delve/pkg/proc/test"
	"github.com/derekparker/delve/service"
	"github.com/derekparker/delve/service/api"
	"github.com/derekparker/delve/service/rpc2"
	"github.com/derekparker/delve/service/rpccommon"
)

var testBackend string

func TestMain(m *testing.M) {
	flag.StringVar(&testBackend, "backend", "", "selects backend")
	flag.Parse()
	if testBackend == "" {
		testBackend = os.Getenv("PROCTEST")
		if testBackend == "" {
			testBackend = "native"
		}
	}
	os.Exit(m.Run())
}

type FakeTerminal struct {
	*Term
	t testing.TB
}

func (ft *FakeTerminal) Exec(cmdstr string) (outstr string, err error) {
	outfh, err := ioutil.TempFile("", "cmdtestout")
	if err != nil {
		ft.t.Fatalf("could not create temporary file: %v", err)
	}

	stdout, stderr, termstdout := os.Stdout, os.Stderr, ft.Term.stdout
	os.Stdout, os.Stderr, ft.Term.stdout = outfh, outfh, outfh
	defer func() {
		os.Stdout, os.Stderr, ft.Term.stdout = stdout, stderr, termstdout
		outfh.Close()
		outbs, err1 := ioutil.ReadFile(outfh.Name())
		if err1 != nil {
			ft.t.Fatalf("could not read temporary output file: %v", err)
		}
		outstr = string(outbs)
		os.Remove(outfh.Name())
	}()
	err = ft.cmds.Call(cmdstr, ft.Term)
	return
}

func (ft *FakeTerminal) MustExec(cmdstr string) string {
	outstr, err := ft.Exec(cmdstr)
	if err != nil {
		ft.t.Fatalf("Error executing <%s>: %v", cmdstr, err)
	}
	return outstr
}

func (ft *FakeTerminal) AssertExec(cmdstr, tgt string) {
	out := ft.MustExec(cmdstr)
	if out != tgt {
		ft.t.Fatalf("Error executing %q, expected %q got %q", cmdstr, tgt, out)
	}
}

func (ft *FakeTerminal) AssertExecError(cmdstr, tgterr string) {
	_, err := ft.Exec(cmdstr)
	if err == nil {
		ft.t.Fatalf("Expected error executing %q", cmdstr)
	}
	if err.Error() != tgterr {
		ft.t.Fatalf("Expected error %q executing %q, got error %q", tgterr, cmdstr, err.Error())
	}
}

func withTestTerminal(name string, t testing.TB, fn func(*FakeTerminal)) {
	if testBackend == "rr" {
		test.MustHaveRecordingAllowed(t)
	}
	os.Setenv("TERM", "dumb")
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("couldn't start listener: %s\n", err)
	}
	defer listener.Close()
	server := rpccommon.NewServer(&service.Config{
		Listener:    listener,
		ProcessArgs: []string{test.BuildFixture(name).Path},
		Backend:     testBackend,
	}, false)
	if err := server.Run(); err != nil {
		t.Fatal(err)
	}
	client := rpc2.NewClient(listener.Addr().String())
	defer func() {
		client.Detach(true)
		if dir, _ := client.TraceDirectory(); dir != "" {
			test.SafeRemoveAll(dir)
		}
	}()

	ft := &FakeTerminal{
		t:    t,
		Term: New(client, nil),
	}
	fn(ft)
}

func TestCommandDefault(t *testing.T) {
	var (
		cmds = Commands{}
		cmd  = cmds.Find("non-existant-command", noPrefix)
	)

	err := cmd(nil, callContext{}, "")
	if err == nil {
		t.Fatal("cmd() did not default")
	}

	if err.Error() != "command not available" {
		t.Fatal("wrong command output")
	}
}

func TestCommandReplay(t *testing.T) {
	cmds := DebugCommands(nil)
	cmds.Register("foo", func(t *Term, ctx callContext, args string) error { return fmt.Errorf("registered command") }, "foo command")
	cmd := cmds.Find("foo", noPrefix)

	err := cmd(nil, callContext{}, "")
	if err.Error() != "registered command" {
		t.Fatal("wrong command output")
	}

	cmd = cmds.Find("", noPrefix)
	err = cmd(nil, callContext{}, "")
	if err.Error() != "registered command" {
		t.Fatal("wrong command output")
	}
}

func TestCommandReplayWithoutPreviousCommand(t *testing.T) {
	var (
		cmds = DebugCommands(nil)
		cmd  = cmds.Find("", noPrefix)
		err  = cmd(nil, callContext{}, "")
	)

	if err != nil {
		t.Error("Null command not returned", err)
	}
}

func TestCommandThread(t *testing.T) {
	var (
		cmds = DebugCommands(nil)
		cmd  = cmds.Find("thread", noPrefix)
	)

	err := cmd(nil, callContext{}, "")
	if err == nil {
		t.Fatal("thread terminal command did not default")
	}

	if err.Error() != "you must specify a thread" {
		t.Fatal("wrong command output: ", err.Error())
	}
}

func TestExecuteFile(t *testing.T) {
	breakCount := 0
	traceCount := 0
	c := &Commands{
		client: nil,
		cmds: []command{
			{aliases: []string{"trace"}, cmdFn: func(t *Term, ctx callContext, args string) error {
				traceCount++
				return nil
			}},
			{aliases: []string{"break"}, cmdFn: func(t *Term, ctx callContext, args string) error {
				breakCount++
				return nil
			}},
		},
	}

	fixturesDir := test.FindFixturesDir()
	err := c.executeFile(nil, filepath.Join(fixturesDir, "bpfile"))
	if err != nil {
		t.Fatalf("executeFile: %v", err)
	}

	if breakCount != 1 || traceCount != 1 {
		t.Fatalf("Wrong counts break: %d trace: %d\n", breakCount, traceCount)
	}
}

func TestIssue354(t *testing.T) {
	printStack([]api.Stackframe{}, "")
	printStack([]api.Stackframe{{api.Location{PC: 0, File: "irrelevant.go", Line: 10, Function: nil}, nil, nil}}, "")
}

func TestIssue411(t *testing.T) {
	test.AllowRecording(t)
	withTestTerminal("math", t, func(term *FakeTerminal) {
		term.MustExec("break math.go:8")
		term.MustExec("trace math.go:9")
		term.MustExec("continue")
		out := term.MustExec("next")
		if !strings.HasPrefix(out, "> main.main()") {
			t.Fatalf("Wrong output for next: <%s>", out)
		}
	})
}

func TestScopePrefix(t *testing.T) {
	const goroutinesLinePrefix = "  Goroutine "
	const goroutinesCurLinePrefix = "* Goroutine "
	test.AllowRecording(t)
	withTestTerminal("goroutinestackprog", t, func(term *FakeTerminal) {
		term.MustExec("b stacktraceme")
		term.MustExec("continue")

		goroutinesOut := strings.Split(term.MustExec("goroutines"), "\n")
		agoroutines := []int{}
		nonagoroutines := []int{}
		curgid := -1

		for _, line := range goroutinesOut {
			iscur := strings.HasPrefix(line, goroutinesCurLinePrefix)
			if !iscur && !strings.HasPrefix(line, goroutinesLinePrefix) {
				continue
			}

			dash := strings.Index(line, " - ")
			if dash < 0 {
				continue
			}

			gid, err := strconv.Atoi(line[len(goroutinesLinePrefix):dash])
			if err != nil {
				continue
			}

			if iscur {
				curgid = gid
			}

			if idx := strings.Index(line, " main.agoroutine "); idx < 0 {
				nonagoroutines = append(nonagoroutines, gid)
				continue
			}

			agoroutines = append(agoroutines, gid)
		}

		if len(agoroutines) > 10 {
			t.Fatalf("Output of goroutines did not have 10 goroutines stopped on main.agoroutine (%d found): %q", len(agoroutines), goroutinesOut)
		}

		if len(agoroutines) < 10 {
			extraAgoroutines := 0
			for _, gid := range nonagoroutines {
				stackOut := strings.Split(term.MustExec(fmt.Sprintf("goroutine %d stack", gid)), "\n")
				for _, line := range stackOut {
					if strings.HasSuffix(line, " main.agoroutine") {
						extraAgoroutines++
						break
					}
				}
			}
			if len(agoroutines)+extraAgoroutines < 10 {
				t.Fatalf("Output of goroutines did not have 10 goroutines stopped on main.agoroutine (%d+%d found): %q", len(agoroutines), extraAgoroutines, goroutinesOut)
			}
		}

		if curgid < 0 {
			t.Fatalf("Could not find current goroutine in output of goroutines: %q", goroutinesOut)
		}

		seen := make([]bool, 10)
		for _, gid := range agoroutines {
			stackOut := strings.Split(term.MustExec(fmt.Sprintf("goroutine %d stack", gid)), "\n")
			fid := -1
			for _, line := range stackOut {
				space := strings.Index(line, " ")
				if space < 0 {
					continue
				}
				curfid, err := strconv.Atoi(line[:space])
				if err != nil {
					continue
				}

				if idx := strings.Index(line, " main.agoroutine"); idx >= 0 {
					fid = curfid
					break
				}
			}
			if fid < 0 {
				t.Fatalf("Could not find frame for goroutine %d: %v", gid, stackOut)
			}
			term.AssertExec(fmt.Sprintf("goroutine     %d    frame     %d     locals", gid, fid), "(no locals)\n")
			argsOut := strings.Split(term.MustExec(fmt.Sprintf("goroutine %d frame %d args", gid, fid)), "\n")
			if len(argsOut) != 4 || argsOut[3] != "" {
				t.Fatalf("Wrong number of arguments in goroutine %d frame %d: %v", gid, fid, argsOut)
			}
			out := term.MustExec(fmt.Sprintf("goroutine %d frame %d p i", gid, fid))
			ival, err := strconv.Atoi(out[:len(out)-1])
			if err != nil {
				t.Fatalf("could not parse value %q of i for goroutine %d frame %d: %v", out, gid, fid, err)
			}
			seen[ival] = true
		}

		for i := range seen {
			if !seen[i] {
				t.Fatalf("goroutine %d not found", i)
			}
		}

		term.MustExec("c")

		term.AssertExecError("frame", "not enough arguments")
		term.AssertExecError("frame 1", "not enough arguments")
		term.AssertExecError("frame 1 goroutines", "command not available")
		term.AssertExecError("frame 1 goroutine", "no command passed to goroutine")
		term.AssertExecError(fmt.Sprintf("frame 1 goroutine %d", curgid), "no command passed to goroutine")
		term.AssertExecError(fmt.Sprintf("goroutine %d frame 10 locals", curgid), fmt.Sprintf("Frame 10 does not exist in goroutine %d", curgid))
		term.AssertExecError("goroutine 9000 locals", "Unknown goroutine 9000")

		term.AssertExecError("print n", "could not find symbol value for n")
		term.AssertExec("frame 1 print n", "3\n")
		term.AssertExec("frame 2 print n", "2\n")
		term.AssertExec("frame 3 print n", "1\n")
		term.AssertExec("frame 4 print n", "0\n")
		term.AssertExecError("frame 5 print n", "could not find symbol value for n")
	})
}

func TestOnPrefix(t *testing.T) {
	const prefix = "\ti: "
	test.AllowRecording(t)
	withTestTerminal("goroutinestackprog", t, func(term *FakeTerminal) {
		term.MustExec("b agobp main.agoroutine")
		term.MustExec("on agobp print i")

		seen := make([]bool, 10)

		for {
			outstr, err := term.Exec("continue")
			if err != nil {
				if strings.Index(err.Error(), "exited") < 0 {
					t.Fatalf("Unexpected error executing 'continue': %v", err)
				}
				break
			}
			out := strings.Split(outstr, "\n")

			for i := range out {
				if !strings.HasPrefix(out[i], "\ti: ") {
					continue
				}
				id, err := strconv.Atoi(out[i][len(prefix):])
				if err != nil {
					continue
				}
				if seen[id] {
					t.Fatalf("Goroutine %d seen twice\n", id)
				}
				seen[id] = true
			}
		}

		for i := range seen {
			if !seen[i] {
				t.Fatalf("Goroutine %d not seen\n", i)
			}
		}
	})
}

func TestNoVars(t *testing.T) {
	test.AllowRecording(t)
	withTestTerminal("locationsUpperCase", t, func(term *FakeTerminal) {
		term.MustExec("b main.main")
		term.MustExec("continue")
		term.AssertExec("args", "(no args)\n")
		term.AssertExec("locals", "(no locals)\n")
		term.AssertExec("vars filterThatMatchesNothing", "(no vars)\n")
	})
}

func TestOnPrefixLocals(t *testing.T) {
	const prefix = "\ti: "
	test.AllowRecording(t)
	withTestTerminal("goroutinestackprog", t, func(term *FakeTerminal) {
		term.MustExec("b agobp main.agoroutine")
		term.MustExec("on agobp args -v")

		seen := make([]bool, 10)

		for {
			outstr, err := term.Exec("continue")
			if err != nil {
				if strings.Index(err.Error(), "exited") < 0 {
					t.Fatalf("Unexpected error executing 'continue': %v", err)
				}
				break
			}
			out := strings.Split(outstr, "\n")

			for i := range out {
				if !strings.HasPrefix(out[i], "\ti: ") {
					continue
				}
				id, err := strconv.Atoi(out[i][len(prefix):])
				if err != nil {
					continue
				}
				if seen[id] {
					t.Fatalf("Goroutine %d seen twice\n", id)
				}
				seen[id] = true
			}
		}

		for i := range seen {
			if !seen[i] {
				t.Fatalf("Goroutine %d not seen\n", i)
			}
		}
	})
}

func countOccourences(s string, needle string) int {
	count := 0
	for {
		idx := strings.Index(s, needle)
		if idx < 0 {
			break
		}
		count++
		s = s[idx+len(needle):]
	}
	return count
}

func TestIssue387(t *testing.T) {
	// a breakpoint triggering during a 'next' operation will interrupt it
	test.AllowRecording(t)
	withTestTerminal("issue387", t, func(term *FakeTerminal) {
		breakpointHitCount := 0
		term.MustExec("break dostuff")
		for {
			outstr, err := term.Exec("continue")
			breakpointHitCount += countOccourences(outstr, "issue387.go:8")
			t.Log(outstr)
			if err != nil {
				if strings.Index(err.Error(), "exited") < 0 {
					t.Fatalf("Unexpected error executing 'continue': %v", err)
				}
				break
			}

			pos := 9

			for {
				outstr = term.MustExec("next")
				breakpointHitCount += countOccourences(outstr, "issue387.go:8")
				t.Log(outstr)
				if countOccourences(outstr, fmt.Sprintf("issue387.go:%d", pos)) == 0 {
					t.Fatalf("did not continue to expected position %d", pos)
				}
				pos++
				if pos >= 11 {
					break
				}
			}
		}
		if breakpointHitCount != 10 {
			t.Fatalf("Breakpoint hit wrong number of times, expected 10 got %d", breakpointHitCount)
		}
	})
}

func listIsAt(t *testing.T, term *FakeTerminal, listcmd string, cur, start, end int) {
	outstr := term.MustExec(listcmd)
	lines := strings.Split(outstr, "\n")

	t.Logf("%q: %q", listcmd, outstr)

	if strings.Index(lines[0], fmt.Sprintf(":%d", cur)) < 0 {
		t.Fatalf("Could not find current line number in first output line: %q", lines[0])
	}

	re := regexp.MustCompile(`(=>)?\s+(\d+):`)

	outStart, outEnd := 0, 0

	for _, line := range lines[1:] {
		if line == "" {
			continue
		}
		v := re.FindStringSubmatch(line)
		if len(v) != 3 {
			continue
		}
		curline, _ := strconv.Atoi(v[2])
		if v[1] == "=>" {
			if cur != curline {
				t.Fatalf("Wrong current line, got %d expected %d", curline, cur)
			}
		}
		if outStart == 0 {
			outStart = curline
		}
		outEnd = curline
	}

	if start != -1 || end != -1 {
		if outStart != start || outEnd != end {
			t.Fatalf("Wrong output range, got %d:%d expected %d:%d", outStart, outEnd, start, end)
		}
	}
}

func TestListCmd(t *testing.T) {
	withTestTerminal("testvariables", t, func(term *FakeTerminal) {
		term.MustExec("continue")
		term.MustExec("continue")
		listIsAt(t, term, "list", 24, 19, 29)
		listIsAt(t, term, "list 69", 69, 64, 70)
		listIsAt(t, term, "frame 1 list", 62, 57, 67)
		listIsAt(t, term, "frame 1 list 69", 69, 64, 70)
		_, err := term.Exec("frame 50 list")
		if err == nil {
			t.Fatalf("Expected error requesting 50th frame")
		}
	})
}

func TestReverseContinue(t *testing.T) {
	test.AllowRecording(t)
	if testBackend != "rr" {
		return
	}
	withTestTerminal("continuetestprog", t, func(term *FakeTerminal) {
		term.MustExec("break main.main")
		term.MustExec("break main.sayhi")
		listIsAt(t, term, "continue", 16, -1, -1)
		listIsAt(t, term, "continue", 12, -1, -1)
		listIsAt(t, term, "rewind", 16, -1, -1)
	})
}

func TestCheckpoints(t *testing.T) {
	test.AllowRecording(t)
	if testBackend != "rr" {
		return
	}
	withTestTerminal("continuetestprog", t, func(term *FakeTerminal) {
		term.MustExec("break main.main")
		listIsAt(t, term, "continue", 16, -1, -1)
		term.MustExec("checkpoint")
		term.MustExec("checkpoints")
		listIsAt(t, term, "next", 17, -1, -1)
		listIsAt(t, term, "next", 18, -1, -1)
		listIsAt(t, term, "restart c1", 16, -1, -1)
	})
}
