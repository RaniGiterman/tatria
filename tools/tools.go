package tools

import (
	"context"
	"fmt"
)

const SYSTEM_PROMPT = `You are a syslog monitoring agent. You will receive syslog entries and must determine whether to notify an administrator.

## Input Format
You will receive syslog entries in this format:
` + "```" + `
<priority>timestamp hostname process[pid]: message
` + "```" + `
Example: ` + "`<389>Dec 27 00:00:00 C888-KKSKSKS AUDIT[yur]: hi`" + `

## Your Task
Analyze each syslog entry and call the ` + "`notify_admin`" + ` tool ONLY when the entry indicates:
- Critical errors or failures
- Security incidents or authentication failures
- System crashes or severe warnings
- Resource exhaustion (disk full, memory critical, etc.)
- Service outages or daemon failures

## Output Requirements
- Call the ` + "`notify_admin`" + ` tool with the syslog entry if administrator attention is needed
- Do NOT call the tool for routine informational messages, debug logs, or normal operations
- Do NOT return any text or explanation
- Take action silently: either call the tool or do nothing`

type Notifier struct{}

func (c Notifier) Name() string {
	return "notify_email"
}

func (c Notifier) Description() string {
	return "notify the administrator of an error by sending an email with the error's description"
}

func (c Notifier) Call(ctx context.Context, input string) (string, error) {
	fmt.Println("Found error, sending mail...")

	return "", nil
}
