// This file is based on code from the original project:
// https://github.com/jamescherti/watch-xfce-xfconf
// Copyright (C) 2021-2025 James Cherti
// Licensed under the MIT License.

package main

import (
	"encoding/xml"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Xfconf struct {
	xfconfItems map[string]XfconfItem
}

type XfconfItem struct {
	Channel       string
	PropertyPath  string
	PropertyType  string
	PropertyValue interface{}
}

type Property struct {
	Name     string     `xml:"name,attr"`
	Type     string     `xml:"type,attr"`
	Value    string     `xml:"value,attr"`
	Property []Property `xml:"property"`
}

func NewXfconf() (*Xfconf, error) {
	xfconf := &Xfconf{
		xfconfItems: make(map[string]XfconfItem),
	}

	// TODO: check XDG_CONFIG_DIR here
	dirXfconf := filepath.Join(os.Getenv("HOME"), ".config", "xfce4", "xfconf", "xfce-perchannel-xml")
	files, err := filepath.Glob(filepath.Join(dirXfconf, "*.xml"))
	if err != nil {
		return nil, err
	}

	for _, xmlFile := range files {
		if err := xfconf.parseXfconfPerchannelXML(xmlFile); err != nil {
			return nil, err
		}
	}

	return xfconf, nil
}

// Diff returns the settings that have been changed.
func (xfconf *Xfconf) Diff() ([]string, error) {

	newXfconf, err := NewXfconf()
	if err != nil {
		return nil, err
	}

	before := xfconf.String()
	after := newXfconf.String()

	xfconf.xfconfItems = newXfconf.xfconfItems

	return diffStrings(before, after), nil
}

// String returns the string representation of the Xfconf settings.
func (xfconf *Xfconf) String() string {
	var commands []string
	for _, item := range xfconf.xfconfItems {
		cmd := fmt.Sprintf("xfconf-query --create -c %s -p %s",
			quoteCommand(item.Channel),
			quoteCommand(item.PropertyPath))

		if item.PropertyType == "array" {
			for _, arrayItem := range item.PropertyValue.([]interface{}) {
				arrayItemType := arrayItem.(map[string]interface{})["type"].(string)
				arrayItemValue := arrayItem.(map[string]interface{})["value"].(string)
				cmd += fmt.Sprintf(" --type %s --set %s",
					quoteCommand(arrayItemType),
					quoteCommand(arrayItemValue))
			}
		} else {
			cmd += fmt.Sprintf(" --type %s --set %s",
				quoteCommand(item.PropertyType),
				quoteCommand(fmt.Sprintf("%v", item.PropertyValue)))
		}

		commands = append(commands, cmd)
	}

	return strings.Join(commands, "\n")
}

// parseXfconfPerchannelXML parses the Xfconf XML file.
func (xfconf *Xfconf) parseXfconfPerchannelXML(xmlFile string) error {
	var channel struct {
		Name     string     `xml:"name,attr"`
		Property []Property `xml:"property"`
	}

	data, err := ioutil.ReadFile(xmlFile)
	if err != nil {
		return err
	}
	if err := xml.Unmarshal(data, &channel); err != nil {
		return err
	}

	/*if channel.Version != "1.0" || channel.XMLName.Local != "channel" {
	    return errors.New("invalid XML file: " + xmlFile)
	}*/

	for _, prop := range channel.Property {
		xfconf.parseProperty(prop, channel.Name, "")
	}

	return nil
}

// parseProperty parses a single property.
func (xfconf *Xfconf) parseProperty(prop Property, channelName, propertyPath string) {
	curPropertyPath := propertyPath + "/" + prop.Name
	if prop.Type == "empty" {
		for _, subProp := range prop.Property {
			xfconf.parseProperty(subProp, channelName, curPropertyPath)
		}
		return
	}

	var propertyValue interface{}
	if prop.Type == "array" {
		var arrayItems []interface{}
		for _, arrayItem := range prop.Property {
			arrayItems = append(arrayItems, map[string]interface{}{
				"type":  arrayItem.Type,
				"value": arrayItem.Value,
			})
		}
		propertyValue = arrayItems
	} else {
		propertyValue = prop.Value
	}

	xfconf.xfconfItems[curPropertyPath] = XfconfItem{
		Channel:       channelName,
		PropertyPath:  curPropertyPath,
		PropertyType:  prop.Type,
		PropertyValue: propertyValue,
	}
}

func quoteCommand(command string) string {
	return "'" + strings.ReplaceAll(command, "'", "'\\''") + "'"
}

func diffStrings(before, after string) []string {
	beforeLines := strings.Split(before, "\n")
	afterLines := strings.Split(after, "\n")

	var diff []string
	for _, line := range afterLines {
		if !contains(beforeLines, line) {
			diff = append(diff, line)
		}
	}

	return diff
}

// contains checks if a string is in a slice of strings.
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

func recordProfile() {
	fmt.Println("Recording changes to xfconf...")
	xfconf, err := NewXfconf()
	if err != nil {
		fmt.Println("Error initializing Xfconf:", err)
		return
	}

	blue := color.New(color.FgHiBlue).SprintFunc()

	for {
		changes, err := xfconf.Diff()
		if err != nil {
			fmt.Println("Error getting diff:", err)
			return
		}

		if len(changes) > 0 {
			for _, change := range changes {
				fmt.Printf("%s %s\n", blue("•"), change)
			}
		}

		time.Sleep(250 * time.Millisecond)
	}
}
