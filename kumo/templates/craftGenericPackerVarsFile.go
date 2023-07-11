package templates

import (
	"os"
	"path/filepath"
	"text/template"

	"github.com/ed3899/kumo/utils"
	"github.com/pkg/errors"
)

type PackerEnvironmentI interface {
	*AWS_PackerEnvironment | *GeneralPackerEnvironment
}

func CraftGenericPackerVarsFile[E PackerEnvironmentI](templateName, packerVarsFileName string, env E) (resultingPackerVarsPath string, err error) {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		err = errors.Wrap(err, "Error occurred while getting current working directory")
		return "", err
	}

	// Parse template file
	generalTemplatePath := filepath.Join(cwd, "templates", templateName)
	tmpl, err := template.ParseFiles(generalTemplatePath)
	if err != nil {
		err = errors.Wrap(err, "Error occurred while parsing Packer AWS Vars template file")
		return "", err
	}

	// Get Packer HCL directory path
	phcldir, err := utils.GetPackerHclDirPath()
	if err != nil {
		err = errors.Wrap(err, "Error occurred while getting Packer HCL directory path")
		return "", err
	}

	// Create Packer Vars file
	resultingPackerVarsPath = filepath.Join(phcldir, packerVarsFileName)
	file, err := os.Create(resultingPackerVarsPath)
	if err != nil {
		err = errors.Wrap(err, "Error occurred while creating Packer AWS Vars file")
		return "", err
	}
	defer file.Close()

	// Execute template file
	err = tmpl.Execute(file, env)
	if err != nil {
		err = errors.Wrap(err, "Error occurred while executing Packer AWS Vars template file")
		return "", err
	}

	// Return path to Packer Vars file
	return resultingPackerVarsPath, nil
}
