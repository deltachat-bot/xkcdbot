package deltachat

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"

	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/channel"
)

// Delta Chat core Event
type Event struct {
	Type               string
	Msg                string
	File               string
	ChatId             uint64
	MsgId              uint64
	ContactId          uint64
	MsgIds             []uint64
	Timer              int
	Progress           uint
	Comment            string
	Path               string
	StatusUpdateSerial uint
}

type _Params struct {
	ContextId uint64
	Event     *Event
}

// Delta Chat core RPC
type Rpc interface {
	Start() error
	Stop()
	GetEventChannel(accountId uint64) <-chan *Event
	Call(method string, params ...any) error
	CallResult(result any, method string, params ...any) error
	String() string
}

// Delta Chat core RPC working over IO
type RpcIO struct {
	Stderr      *os.File
	AccountsDir string
	Cmd         string
	cmd         *exec.Cmd
	stdin       io.WriteCloser
	client      *jrpc2.Client
	ctx         context.Context
	events      map[uint64]chan *Event
	eventsMutex sync.Mutex
	closed      bool
}

func NewRpcIO() *RpcIO {
	return &RpcIO{Cmd: "deltachat-rpc-server", Stderr: os.Stderr}
}

// Implement Stringer.
func (self *RpcIO) String() string {
	return fmt.Sprintf("Rpc(AccountsDir=%#v)", self.AccountsDir)
}

func (self *RpcIO) Start() error {
	self.closed = false
	self.cmd = exec.Command(self.Cmd)
	if self.AccountsDir != "" {
		self.cmd.Env = append(os.Environ(), "DC_ACCOUNTS_PATH="+self.AccountsDir)
	}
	self.cmd.Stderr = self.Stderr
	self.stdin, _ = self.cmd.StdinPipe()
	stdout, _ := self.cmd.StdoutPipe()
	if err := self.cmd.Start(); err != nil {
		self.closed = true
		return err
	}

	self.ctx = context.Background()
	self.events = make(map[uint64]chan *Event)
	options := jrpc2.ClientOptions{OnNotify: self._onNotify}
	self.client = jrpc2.NewClient(channel.Line(stdout, self.stdin), &options)
	return nil
}

func (self *RpcIO) Stop() {
	self.eventsMutex.Lock()
	if !self.closed {
		self.stdin.Close()
		self.cmd.Process.Wait()
		for _, value := range self.events {
			close(value)
		}
		self.closed = true
	}
	self.eventsMutex.Unlock()
}

func (self *RpcIO) GetEventChannel(accountId uint64) <-chan *Event {
	self._initEventChannel(accountId)
	return self.events[accountId]
}

func (self *RpcIO) Call(method string, params ...any) error {
	_, err := self.client.Call(self.ctx, method, params)
	return err
}

func (self *RpcIO) CallResult(result any, method string, params ...any) error {
	return self.client.CallResult(self.ctx, method, params, &result)
}

func (self *RpcIO) _initEventChannel(accountId uint64) {
	self.eventsMutex.Lock()
	if _, ok := self.events[accountId]; !ok {
		self.events[accountId] = make(chan *Event, 10)
	}
	self.eventsMutex.Unlock()
}

func (self *RpcIO) _onNotify(req *jrpc2.Request) {
	if req.Method() == "event" {
		var params _Params
		req.UnmarshalParams(&params)
		self._initEventChannel(params.ContextId)
		if !self.closed {
			go func() { self.events[params.ContextId] <- params.Event }()
		}
	}
}
