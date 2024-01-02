package vienv

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

type Tool struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	RepoLink string `json:"repoLink"`
	Path     string `json:"path"`
}

// typically our config.json
type Config struct {
	EnvName      string `json:"envName"`
	Architecture string `json:"arch"`
	Os           string `json:"os"`
	WorkingDir   string `json:"workingDirectory"`
	Tools        []Tool `json:"tools"`
}

type VirtualEnvironment struct {
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
func (v *VirtualEnvironment) createConfig() error {
	updatedContent, err := json.MarshalIndent(&v.config, "", "    ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(v.configPath, updatedContent, os.ModePerm); err != nil {
		return err
	}
	return nil
}

// this utility fetches tool from a repository...
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
	workDirAbsolutePath, err := filepath.Abs(dir)
	fmt.Printf("Working directory abs path: %v\n", workDirAbsolutePath)
	if err != nil {
		fmt.Printf("Error getting work directory %v Error: %v\n", dir, err)
		return nil, err
	}

	//this is actually not really necessary I think
	if _, err := os.Stat(workDirAbsolutePath); os.IsNotExist(err) {
		fmt.Printf("Directory %v does not exist\n", workDirAbsolutePath)
		return nil, err
	}

	//create configFile
	config := Config{
		WorkingDir:   workDirAbsolutePath,
		EnvName:      name,
		Architecture: runtime.GOARCH,
		Os:           runtime.GOOS,
		Tools:        make([]Tool, 0),
	}
	fmt.Printf("Config object generated: %v\n", config)
	config_path := filepath.Join(workDirAbsolutePath, "bioenv.json")
	fmt.Printf("Config file path: %v\n", config_path)

	venv = VirtualEnvironment{

		config:     &config,
		configPath: config_path,
		isActive:   true,
		stopped:    false,
	}

	fmt.Printf("Bioenv virtual env: %v\n", venv)

	fmt.Println("Creating virtual environment config...")
	if err = venv.createConfig(); err != nil {
		fmt.Println("Error could not create config file")
		return nil, err
	}

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
			fmt.Printf("Could not create config file %v\n", config_path)
			return err
		}
	} else {
		file, err := os.Create(config_path)
		if err != nil {
			fmt.Printf("Could not create config file %v\n", config_path)
			return err
		}
		defer file.Close()
		if _, err := loadConfigFromPath(config_path, &config); err != nil {
			fmt.Printf("Could not create config file %v\n", config_path)
			return err
		}
	}

	return nil

}

func loadConfigFromPath(path string, config *Config) (*Config, error) {
	fileContent, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Could not read config file %v\n", path)
		return nil, err
	}
	if err := json.Unmarshal([]byte(string(fileContent)), &config); err != nil {
		fmt.Printf("Error deserializing config file\n")
		return nil, err
	}

	return config, nil

}
