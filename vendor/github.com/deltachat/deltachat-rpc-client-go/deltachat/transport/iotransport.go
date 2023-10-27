package transport

import (
	"context"
	"io"
	"os"
	"os/exec"
	"sync"

	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/channel"
)

// Delta Chat RPC transport using external deltachat-rpc-server program
type IOTransport struct {
	Stderr      io.Writer
	AccountsDir string
	Cmd         string
	cmd         *exec.Cmd
	stdin       io.WriteCloser
	client      *jrpc2.Client
	ctx         context.Context
	cancel      context.CancelFunc
	mu          sync.Mutex
}

func NewIOTransport() *IOTransport {
	return &IOTransport{Cmd: deltachatRpcServerBin, Stderr: os.Stderr}
}

func (self *IOTransport) Open() error {
	self.mu.Lock()
	defer self.mu.Unlock()

	if self.ctx != nil && self.ctx.Err() == nil {
		return &TransportStartedErr{}
	}

	self.ctx, self.cancel = context.WithCancel(context.Background())
	self.cmd = exec.CommandContext(self.ctx, self.Cmd)
	if self.AccountsDir != "" {
		self.cmd.Env = append(os.Environ(), "DC_ACCOUNTS_PATH="+self.AccountsDir)
	}
	self.cmd.Stderr = self.Stderr
	self.stdin, _ = self.cmd.StdinPipe()
	stdout, _ := self.cmd.StdoutPipe()
	if err := self.cmd.Start(); err != nil {
		self.cancel()
		return err
	}

	self.client = jrpc2.NewClient(channel.Line(stdout, self.stdin), nil)
	return nil
}

func (self *IOTransport) Close() {
	self.mu.Lock()
	defer self.mu.Unlock()

	if self.ctx == nil || self.ctx.Err() != nil {
		return
	}

	self.stdin.Close()
	self.cancel()
	self.cmd.Process.Wait() //nolint:errcheck
}

func (self *IOTransport) Call(ctx context.Context, method string, params ...any) error {
	_, err := self.client.Call(ctx, method, params)
	return err
}

func (self *IOTransport) CallResult(ctx context.Context, result any, method string, params ...any) error {
	return self.client.CallResult(ctx, method, params, &result)
}

// TransportStartedErr is returned by IOTransport.Open() if the Transport is already started
type TransportStartedErr struct{}

func (self *TransportStartedErr) Error() string {
	return "Transport is already started"
}
