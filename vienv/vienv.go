package vienv

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

type Tool struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	RepoLink string `json:"repoLink"`
	Path     string `json:"path"`
}

//typically our config.json
type Config struct {
	EnvName      string `json:"envName"`
	Architecture string `json:"arch"`
	WorkingDir   string `json:"workingDirectory"`
	Tools        []Tool `json:"tools"`
}

type VirtualEnvironment struct {
	//tempName string
	//envName    string
	configPath string
	config     *Config
	isActive   bool
	stopped    bool
}

var venv VirtualEnvironment

func (v *VirtualEnvironment) InstallTool(toolName string) error {
	//fetch tool into the bioenv/tools dir
	tool, err := v.getToolFromRepo(toolName)
	if err != nil {
		fmt.Printf("Error installing tool: %s\n", toolName)
		return err
	}

	//add tool to env and update config...
	if err := v.updateConfig(*tool); err != nil {
		fmt.Printf("Error updating configuration: %v\n", err)
		return err
	}

	//add tool config to v.config json... toolName, version,

	return nil
}

func (v *VirtualEnvironment) updateConfig(newTool Tool) error {
	fileContent, err := os.ReadFile("config.json")
	if err != nil {
		return err
	}
	if err := json.Unmarshal(fileContent, &v.config); err != nil {
		return err
	}
	v.config.Tools = append(v.config.Tools, newTool)

	updatedContent, err := json.MarshalIndent(&v.config, "", "    ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(v.configPath, updatedContent, os.ModePerm); err != nil {
		return err
	}
	return nil
}

//this utility fetches tool from a repository...
func (v *VirtualEnvironment) getToolFromRepo(toolName string) (*Tool, error) {
	fmt.Printf("Installing new bio tool: %s", toolName)
	repoURL := "https://github.com/propenster/bioenv"
	toolsDirFromRepo := fmt.Sprintf("tools/%s", toolName) //this is the subdirectory on the remote repo for this particular tool

	//this is on our pc where we want the contents of gatk to be dumped
	toolsLocalTargetDir := fmt.Sprintf("%v\\%v", v.config.WorkingDir, "tools")

	cmd := exec.Command("git", "clone", "--filter=blob:none", "--sparse", fmt.Sprintf("%s.git", repoURL))
	cmd.Dir = toolsLocalTargetDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		//log.Fatal(err)
		return nil, err
	}

	cmd = exec.Command("git", "sparse-checkout", "init", "--cone")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	cmd = exec.Command("git", "sparse-checkout", "set", toolsDirFromRepo)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	//create a tool object
	toolRepoUrl := fmt.Sprintf("%s/%s/%s", repoURL, "tree/master/bioenv/tools", toolName)
	tool := &Tool{
		Name:     toolName,
		Version:  "0",
		RepoLink: toolRepoUrl,
		Path:     toolsLocalTargetDir,
	}
	return tool, nil
}

func Init(dir, name string) (*VirtualEnvironment, error) {
	//make sure dir is valid dir...
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Printf("Directory %v does not exist", dir)
		return nil, err
	}

	//create configFile
	var config Config
	config_path := fmt.Sprintf("%s\\%s", dir, "bioenv.json")
	if err := loadConfig(config_path, config); err != nil {
		fmt.Printf("Error creating config file %v", config_path)
		return nil, err
	}

	//done
	config.EnvName = name

	venv = VirtualEnvironment{

		config:     &config,
		configPath: config_path,
		isActive:   true,
		stopped:    false,
	}

	fmt.Printf("Bioenv virtual env: %v", venv)

	//addendum - setPrompt
	if runtime.GOOS == "windows" {
		os.Setenv("PROMPT", config.EnvName)
	} else {
		os.Setenv("PS1", config.EnvName)
	}

	return &venv, nil
}

func loadConfig(config_path string, config Config) error {
	if _, err := os.Stat(config_path); os.IsExist(err) {
		//config file already exists.
		if _, err := loadConfigFromPath(config_path, &config); err != nil {
			fmt.Printf("Could not create config file %v", config_path)
			return err
		}
	} else {
		file, err := os.Create(config_path)
		if err != nil {
			fmt.Printf("Could not create config file %v", config_path)
			return err
		}
		defer file.Close()
		if _, err := loadConfigFromPath(config_path, &config); err != nil {
			fmt.Printf("Could not create config file %v", config_path)
			return err
		}
	}

	return nil

}

func loadConfigFromPath(path string, config *Config) (*Config, error) {
	fileContent, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Could not read config file %v", path)
		return nil, err
	}
	if err := json.Unmarshal([]byte(string(fileContent)), &config); err != nil {
		fmt.Printf("Error deserializing config file")
		return nil, err
	}

	return config, nil

}
