package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// Recursive XML element parser
func parseXMLElement(dec *xml.Decoder, se xml.StartElement) (interface{}, error) {
	children := make(map[string]interface{})
	hasChildren := false
	var textContent string

	for {
		token, err := dec.Token()
		if err != nil {
			return nil, err
		}

		switch t := token.(type) {
		case xml.StartElement:
			hasChildren = true
			childVal, err := parseXMLElement(dec, t)
			if err != nil {
				return nil, err
			}
			children[t.Name.Local] = childVal
		case xml.CharData:
			textContent += strings.TrimSpace(string(t))
		case xml.EndElement:
			if t.Name == se.Name {
				if hasChildren {
					return children, nil
				}
				return textContent, nil
			}
		}
	}
}

// Entry point for XML parsing
func parseXML(dec *xml.Decoder) (map[string]interface{}, error) {
	for {
		token, err := dec.Token()
		if err == io.EOF {
			return nil, fmt.Errorf("empty XML document")
		}
		if err != nil {
			return nil, err
		}
		if se, ok := token.(xml.StartElement); ok {
			content, err := parseXMLElement(dec, se)
			if err != nil {
				return nil, err
			}
			return map[string]interface{}{se.Name.Local: content}, nil
		}
	}
}

// Load the structure from a file based on its format
func loadStructureFromFile(path string) (map[string]interface{}, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buf := make([]byte, 1024)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return nil, err
	}
	head := string(buf[:n])

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	var structure interface{}

	switch {
	case strings.HasPrefix(head, "{"):
		fmt.Println("📂 Format: JSON")
		if err := json.NewDecoder(file).Decode(&structure); err != nil {
			return nil, err
		}
	case strings.Contains(head, "---") || strings.Contains(head, ":"):
		fmt.Println("📂 Format: YAML")
		data, _ := io.ReadAll(file)
		if err := yaml.Unmarshal(data, &structure); err != nil {
			return nil, err
		}
	case strings.HasPrefix(strings.TrimSpace(head), "<"):
		fmt.Println("📂 Format: XML")
		dec := xml.NewDecoder(file)
		structure, err = parseXML(dec)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unable to detect file format")
	}

	if m, ok := structure.(map[string]interface{}); ok {
		return m, nil
	}
	return nil, fmt.Errorf("structure is not a map")
}

// Recursively create directories and files from the structure
func createStructure(data map[string]interface{}, basePath string) error {
	for name, value := range data {
		currentPath := filepath.Join(basePath, name)

		if m, ok := value.(map[string]interface{}); ok {
			hasAtKey := false
			for k := range m {
				if strings.HasPrefix(k, "@") {
					hasAtKey = true
					break
				}
			}
			if !hasAtKey {
				if err := os.MkdirAll(currentPath, 0755); err != nil {
					return err
				}
				if err := createStructure(m, currentPath); err != nil {
					return err
				}
				continue
			}
		}

		dir := filepath.Dir(currentPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}

		content := fmt.Sprintf("%v", value)
		if content == "<nil>" {
			content = ""
		}
		if err := os.WriteFile(currentPath, []byte(content), 0644); err != nil {
			return err
		}
	}
	return nil
}

// Main function
func main() {
	// Check for command-line argument
	var inputFile string
	if len(os.Args) > 1 {
		inputFile = os.Args[1]
		if _, err := os.Stat(inputFile); os.IsNotExist(err) {
			fmt.Printf("❌ The file '%s' is not exist\n", inputFile)
			return
		}
		fmt.Printf("📂 Use: %s\n", inputFile)
	} else {
		// No argument provided, look for default files
		inputFiles := []string{"structure.json", "structure.yaml", "structure.yml", "structure.xml"}
		for _, file := range inputFiles {
			if _, err := os.Stat(file); err == nil {
				inputFile = file
				break
			}
		}

		if inputFile == "" {
			fmt.Println("❌ Not found structure.json / structure.yaml / structure.xml")
			return
		}
	}

	structure, err := loadStructureFromFile(inputFile)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	for projectName, projectData := range structure {
		projectMap, ok := projectData.(map[string]interface{})
		if !ok {
			fmt.Println("❌ Error: пproject data is not a valid dict")
			return
		}

		if err := os.MkdirAll(projectName, 0755); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}

		if err := createStructure(projectMap, projectName); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}

		fmt.Printf("✅ Structure '%s' created successfully!\n", projectName)
		break
	}
}
