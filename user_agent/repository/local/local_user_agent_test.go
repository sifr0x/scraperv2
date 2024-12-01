package local_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/datahuys/scraperv2/user_agent/repository/local"
	"github.com/stretchr/testify/assert"
)

func Test_GetRandomUserAgent(t *testing.T) {
	t.Run("simple-random", func(t *testing.T) {
		agents := []string{"abc", "dfg"}
		agentsJson, err := json.Marshal(agents)
		assert.NoError(t, err)

		file, err := os.CreateTemp("", "agent_test")
		assert.NoError(t, err)
		defer os.Remove(file.Name())

		_, err = file.Write(agentsJson)
		assert.NoError(t, err)
		err = file.Close()
		assert.NoError(t, err)

		time.Sleep(2 * time.Second)

		repo, err := local.NewLocalUserAgentRepository(file.Name())
		assert.NoError(t, err)

		time.Sleep(1 * time.Second)

		random, err := repo.GetRandomUserAgent()
		assert.NoError(t, err)
		assert.Contains(t, agents, random)
	})

	t.Run("agents-replace", func(t *testing.T) {
		agents := []string{"abc", "dfg"}
		agentsJson, err := json.Marshal(agents)
		assert.NoError(t, err)

		file, err := os.CreateTemp("", "agent_test")
		assert.NoError(t, err)
		defer os.Remove(file.Name())

		_, err = file.Write(agentsJson)
		assert.NoError(t, err)
		err = file.Close()
		assert.NoError(t, err)

		repo, err := local.NewLocalUserAgentRepository(file.Name())
		assert.NoError(t, err)
		random, err := repo.GetRandomUserAgent()
		assert.NoError(t, err)
		assert.Contains(t, agents, random)

		time.Sleep(1 * time.Second)

		// Create new agents file
		agentsNew := []string{"ggg", "hhh"}
		agentsNewJson, err := json.Marshal(agentsNew)
		assert.NoError(t, err)

		// Rewriting file
		file2, err := os.Create(file.Name())
		assert.NoError(t, err)
		_, err = file2.Write(agentsNewJson)
		assert.NoError(t, err)
		err = file2.Close()
		assert.NoError(t, err)

		time.Sleep(5 * time.Second)

		// Testing if new agents have been applied
		random2, err := repo.GetRandomUserAgent()
		assert.NoError(t, err)
		assert.Contains(t, agentsNew, random2)

		// Making sure random does not return old agent
		random3, err := repo.GetRandomUserAgent()
		assert.NoError(t, err)
		assert.NotContains(t, agents, random3)
	})

	t.Run("test-all-accessed", func(t *testing.T) {
		agents := []string{"abc", "dfg", "foo"}
		agentsJson, err := json.Marshal(agents)
		assert.NoError(t, err)

		file, err := os.CreateTemp("", "agent_test")
		assert.NoError(t, err)
		defer os.Remove(file.Name())

		_, err = file.Write(agentsJson)
		assert.NoError(t, err)
		err = file.Close()
		assert.NoError(t, err)

		time.Sleep(2 * time.Second)

		repo, err := local.NewLocalUserAgentRepository(file.Name())
		assert.NoError(t, err)

		time.Sleep(1 * time.Second)

		first := false
		second := false
		third := false
		i := 0

		for !(first && second && third) {
			i++
			random, err := repo.GetRandomUserAgent()
			if err != nil {
				continue
			}

			if random == agents[0] {
				first = true
			}
			if random == agents[1] {
				second = true
			}
			if random == agents[2] {
				third = true
			}
		}
	})
}
