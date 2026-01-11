package langchain

import (
	llm_tools "tatria/tools"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/tools"
)

var Executor *agents.Executor

func Init() error {
	llm, err := openai.New(
		// openai.WithModel("gpt-4o-mini"),
		openai.WithModel("gpt-5-chat-latest"),
	)
	if err != nil {
		return err
	}
	agentTools := []tools.Tool{
		llm_tools.Notifier{},
	}

	agent := agents.NewOneShotAgent(llm,
		agentTools,
		agents.WithMaxIterations(3),
		agents.WithPromptPrefix(llm_tools.SYSTEM_PROMPT+"\n\n"),
	)

	Executor = agents.NewExecutor(agent)

	return nil
}
