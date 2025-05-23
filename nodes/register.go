package nodes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

type NodeMetadata struct {
	ID       uuid.UUID `yaml:"id"`
	OS       string    `yaml:"os"`
	Hostname string    `yaml:"hostname"`
}

type NodeRegistrationRequest struct {
	DeploymentURL string `json:"deployment_url"`
	MetricsURL    string `json:"metrics_url"`
}

type NodeRegistrationResponse struct {
	NodeUUID string `json:"id"`
}

func CheckNodeFileExists(filePath string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	defer file.Close()

	if file != nil {
		return true, nil
	}

	return false, nil
}

func CreateNodeFile(filePath string, nodeID uuid.UUID) error {
	nodeMetadata := NodeMetadata{
		ID:       nodeID,
		OS:       "linux",
		Hostname: "localhost",
	}

	dirPath := strings.Split(filePath, "/")[:len(strings.Split(filePath, "/"))-2]
	dirPathStr := strings.Join(dirPath, "/")
	err := os.MkdirAll(dirPathStr, os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := yaml.Marshal(nodeMetadata)
	if err != nil {
		return err
	}

	if _, err := file.Write(data); err != nil {
		return err
	}

	return nil
}

func ReadNodeFile(filePath string) (NodeMetadata, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return NodeMetadata{}, err
	}
	defer file.Close()

	var nodeMetadata NodeMetadata
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&nodeMetadata); err != nil {
		return NodeMetadata{}, err
	}

	return nodeMetadata, nil
}

func GetPublicIP() (string, error) {
	cli := http.Client{}
	resp, err := cli.Get("https://api.ipify.org?format=text")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(ip), nil
}

func SendRegisterRequest(orchestratorURL string, nodePort string) (NodeRegistrationResponse, error) {
	ip, err := GetPublicIP()
	if err != nil {
		return NodeRegistrationResponse{}, err
	}

	deployURL := fmt.Sprintf("http://%s:%s/deploy", ip, nodePort)
	metricsURL := fmt.Sprintf("http://%s:%s/metrics", ip, nodePort)
	healthURL := fmt.Sprintf("http://%s:%s/health", ip, nodePort)
	connectionURL := fmt.Sprintf("http://%s:%s/connect", ip, nodePort)

	requestBody, err := json.Marshal(map[string]string{
		"deployment_url": deployURL,
		"metrics_url":    metricsURL,
		"health_url":     healthURL,
		"connection_url": connectionURL,
	})
	if err != nil {
		return NodeRegistrationResponse{}, err
	}

	req, err := http.NewRequest("POST", orchestratorURL+"/nodes/register", bytes.NewBuffer(requestBody))
	if err != nil {
		return NodeRegistrationResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	cli := http.Client{}
	raw_response, err := cli.Do(req)
	if err != nil {
		return NodeRegistrationResponse{}, err
	}
	defer raw_response.Body.Close()

	var response NodeRegistrationResponse
	if err := json.NewDecoder(raw_response.Body).Decode(&response); err != nil {
		return NodeRegistrationResponse{}, err
	}
	if raw_response.StatusCode != http.StatusCreated {
		return NodeRegistrationResponse{}, fmt.Errorf("failed to register node: %s", raw_response.Status)
	}

	return response, nil
}

func RegisterNode(orchestratorURL string) (NodeMetadata, error) {
	filePath := "./.metadata/node.yaml"
	logger := log.New(os.Stdout, "INFO: ", log.LstdFlags)

	ok, err := CheckNodeFileExists(filePath)
	if err != nil {
		return NodeMetadata{}, err
	}

	if !ok {
		logger.Println("Node file does not exist, registering node...")
		nodeRegData, err := SendRegisterRequest(orchestratorURL, "8081")
		if err != nil {
			return NodeMetadata{}, err
		}

		nodeUUID, err := uuid.Parse(nodeRegData.NodeUUID)
		if err != nil {
			return NodeMetadata{}, err
		}

		err = CreateNodeFile(filePath, nodeUUID)
		if err != nil {
			return NodeMetadata{}, err
		}
		logger.Println("Node registered successfully, node file created.")
	}

	nodeMetadata, err := ReadNodeFile(filePath)
	if err != nil {
		return NodeMetadata{}, err
	}

	return nodeMetadata, nil
}
