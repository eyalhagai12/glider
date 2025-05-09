package nodes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

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

func CheckNodeFileExists(filePath string, nodeID uuid.UUID) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	defer file.Close()

	return false, nil
}

func CreateNodeFile(filePath string, nodeID uuid.UUID) error {
	nodeMetadata := NodeMetadata{
		ID:       nodeID,
		OS:       "linux",
		Hostname: "localhost",
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

func SendRegisterRequest(orchestratorURL string, nodeID uuid.UUID) error {
	ip, err := GetPublicIP()
	if err != nil {
		return err
	}

	deployURL := fmt.Sprintf("http://%s/deploy", ip)
	metricsURL := fmt.Sprintf("http://%s/metrics", ip)

	requestBody, err := json.Marshal(map[string]string{
		"deployment_url": deployURL,
		"metrics_url":    metricsURL,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", orchestratorURL+"/node/register", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	cli := http.Client{}
	_, err = cli.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func RegisterNode(orchestratorURL string) error {
	filePath := "./.metadata/node.yaml"
	ok, err := CheckNodeFileExists(filePath, uuid.Nil)
	if err != nil {
		return err
	}

	if !ok {
		nodeID := uuid.New()
		err := CreateNodeFile(filePath, nodeID)
		if err != nil {
			return err
		}
	}

	nodeMetadata, err := ReadNodeFile(filePath)
	if err != nil {
		return err
	}

	err = SendRegisterRequest(orchestratorURL, nodeMetadata.ID)
	if err != nil {
		return err
	}

	return nil
}
