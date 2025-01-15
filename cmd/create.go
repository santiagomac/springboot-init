/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// FLAGS
var (
	projectType string
	projectName string
	groupID     string
	artifactID  string
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a starter Spring Boot project",
	Long: `Create a starter Spring Boot project with the followings types:
  - Web (Spring MVC)
  - Webflux (Reactor)
  `,
	Run: func(cmd *cobra.Command, args []string) {
		createProject()
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// CREATE FLAGS
	createCmd.Flags().StringVarP(&projectType, "type", "t", "web", "Project Type")
	createCmd.Flags().StringVarP(&projectName, "name", "n", "MyProject", "Project name")
	createCmd.Flags().StringVarP(&groupID, "group", "g", "com.example", "Project Group id")
	createCmd.Flags().StringVarP(&artifactID, "artifact", "a", "demo", "Project Artifact Id")
}

func createProject() {
	fmt.Println("Creating project...")

	var dependencies string

	if projectType == "webflux" {
		dependencies = "web"
	} else {
		dependencies = "webflux"
	}

	// URL BUILDING TO CONSUME Spring Initializr
	url := fmt.Sprintf("https://start.spring.io/starter.zip?type=gradle-project&language=java&bootVersion=3.3.0&groupId=%s&artifactId=%s&name=%s&dependencies=%s",
		groupID,
		artifactID,
		projectName,
		dependencies,
	)

	// DOWNLOAD THE PROJECT
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error downloading the project: %v\n", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error in response code: %s\n", resp.Status)
		return
	}

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		fmt.Printf("Error reading the downloaded file: %v\n", err)
		return
	}

	// Unzip the file
	err = unzip(buf.Bytes(), projectName)
	if err != nil {
		fmt.Printf("Error unziping the file: %v\n", err)
		return
	}

	fmt.Printf("Project '%s' succesfully created in folder: '%s'\n", projectName, projectName)
}

func unzip(data []byte, dest string) error {
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return err
	}

	for _, file := range reader.File {
		path := filepath.Join(dest, file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(path, os.ModePerm)
			continue
		}

		os.MkdirAll(filepath.Dir(path), os.ModePerm)

		outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}

		defer outFile.Close()
		rc, err := file.Open()
		if err != nil {
			return err
		}

		defer rc.Close()

		_, err = io.Copy(outFile, rc)
		if err != nil {
			return err
		}
	}

	return nil
}
