package agents

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// ListAgents retrieves all agents for the authenticated user.
func ListAgents() ([]Agent, error) {
	var response AgentListResponse
	err := doJSON(http.MethodGet, "", nil, &response)
	if err != nil {
		return nil, err
	}

	formattedJSON, err := json.MarshalIndent(response.Agents, "", "    ")
	if err != nil {
		return nil, errors.New("failed to format JSON")
	}

	fmt.Println(string(formattedJSON))

	return response.Agents, nil
}

// CreateAgent creates a new agent from the provided request body.
func CreateAgent(body CreateAgentBody) (*CreateAgentResponse, error) {
	var response CreateAgentResponse
	err := doJSON(http.MethodPost, "", body, &response)
	if err != nil {
		return nil, err
	}

	formattedJSON, err := json.MarshalIndent(response, "", "    ")
	if err != nil {
		return nil, errors.New("failed to format JSON")
	}

	fmt.Println(string(formattedJSON))

	return &response, nil
}

// ListEngines retrieves the agent engines enabled for this deployment.
func ListEngines() (*EnginesResponse, error) {
	var response EnginesResponse
	err := doJSON(http.MethodGet, "/engines", nil, &response)
	if err != nil {
		return nil, err
	}

	formattedJSON, err := json.MarshalIndent(response.Engines, "", "    ")
	if err != nil {
		return nil, errors.New("failed to format JSON")
	}

	fmt.Println(string(formattedJSON))

	return &response, nil
}

// GetAgent retrieves detailed information about a specific agent.
func GetAgent(agentID string) (*AgentDetailResponse, error) {
	var response AgentDetailResponse
	err := doJSON(http.MethodGet, "/"+agentID, nil, &response)
	if err != nil {
		return nil, err
	}

	formattedJSON, err := json.MarshalIndent(response, "", "    ")
	if err != nil {
		return nil, errors.New("failed to format JSON")
	}

	fmt.Println(string(formattedJSON))

	return &response, nil
}

// DeleteAgent deletes an agent by ID.
func DeleteAgent(agentID string) error {
	var response DeleteAgentResponse
	err := doJSON(http.MethodDelete, "/"+agentID, nil, &response)
	if err != nil {
		return err
	}

	fmt.Println("Agent deleted")

	return nil
}

// RestartAgent restarts an agent by ID.
func RestartAgent(agentID string) (*RestartResponse, error) {
	var response RestartResponse
	err := doJSON(http.MethodPost, "/"+agentID+"/restart", nil, &response)
	if err != nil {
		return nil, err
	}

	formattedJSON, err := json.MarshalIndent(response, "", "    ")
	if err != nil {
		return nil, errors.New("failed to format JSON")
	}

	fmt.Println(string(formattedJSON))

	return &response, nil
}

// GetAgentLogs retrieves logs for a specific agent.
func GetAgentLogs(agentID string) (string, error) {
	resp, err := doRequest(http.MethodGet, "/"+agentID+"/logs", nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("server returned error %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var response LogsResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", errors.Join(err, errors.New("failed to decode response"))
	}

	fmt.Println(response.Logs)

	return response.Logs, nil
}
