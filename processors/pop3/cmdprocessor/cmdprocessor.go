//Package cmdprocessor provides a set of structs, variables and methods to process POP3 clients' commands.
package cmdprocessor

import (
	"github.com/trapped/gomaild/parsers/textual"
	. "github.com/trapped/gomaild/processors/pop3/session"
	"strings"
	//POP3 commands
	"github.com/trapped/gomaild/processors/pop3/cmdprocessor/apop"
	"github.com/trapped/gomaild/processors/pop3/cmdprocessor/auth"
	"github.com/trapped/gomaild/processors/pop3/cmdprocessor/capa"
	"github.com/trapped/gomaild/processors/pop3/cmdprocessor/dele"
	"github.com/trapped/gomaild/processors/pop3/cmdprocessor/list"
	"github.com/trapped/gomaild/processors/pop3/cmdprocessor/pass"
	"github.com/trapped/gomaild/processors/pop3/cmdprocessor/quit"
	"github.com/trapped/gomaild/processors/pop3/cmdprocessor/retr"
	"github.com/trapped/gomaild/processors/pop3/cmdprocessor/rset"
	"github.com/trapped/gomaild/processors/pop3/cmdprocessor/stat"
	"github.com/trapped/gomaild/processors/pop3/cmdprocessor/stls"
	"github.com/trapped/gomaild/processors/pop3/cmdprocessor/top"
	"github.com/trapped/gomaild/processors/pop3/cmdprocessor/uidl"
	"github.com/trapped/gomaild/processors/pop3/cmdprocessor/user"
)

var (
	//Contains commands and their relative functions (to be executed when a command is issued)
	Commands map[string]func(*Session, textual.Statement) (string, error) = map[string]func(*Session, textual.Statement) (string, error){
		"apop": apop.Process,
		"capa": capa.Process,
		"user": user.Process,
		"pass": pass.Process,
		"stat": stat.Process,
		"list": list.Process,
		"uidl": uidl.Process,
		"top":  top.Process,
		"retr": retr.Process,
		"dele": dele.Process,
		"rset": rset.Process,
		"quit": quit.Process,
		"auth": auth.Process,
		"stls": stls.Process,
	}
)

//Struct to provide a throw-away command processor and session for POP3.
type Processor struct {
	Session     *Session //POP3 session, accessible by both the commands and the client handler
	LastCommand string   //The last successfully issued command
}

//Processes a POP3 command and returns a result.
func (p *Processor) Process(s string) string {
	//Prepare a textual parser.
	parser := textual.Parser{
		Prefix:             "",
		Suffix:             "",
		OpenBrackets:       false,
		Brackets:           []byte{},
		Trim:               true,
		ArgumentSeparators: []byte{' '},
	}

	//Parse the command with the parser.
	z := parser.Parse(s)

	if p.Session.AuthState == AUTHWUSER || p.Session.AuthState == AUTHWPASS {
		z.Name = p.LastCommand
	}

	//Run the processor for the command issued by the client (if it exists) and return the result with the correct POP3 result prefix.
	if _, exists := Commands[strings.ToLower(z.Name)]; exists {
		p.LastCommand = z.Name
		res, err := Commands[strings.ToLower(z.Name)](p.Session, z)
		if err != nil {
			return "-ERR " + err.Error()
		}
		return "+OK " + res
	}
	return "-ERR command not found"
}
