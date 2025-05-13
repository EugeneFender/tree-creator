```markdown
# Project Structure Generator

A simple Go-based tool to generate directory structures from configuration files in JSON, YAML, or XML format.

## Features
- Supports multiple configuration formats: JSON, YAML, XML
- Recursively creates directories and files
- Automatic detection of file format
- Human-readable CLI output

## Requirements
- Go 1.16+ (for `os.MkdirAll` and `io.ReadFile`)

## Installation
```bash
go get github.com/EugeneFender/tree-creator
```

## Usage
1. Create a configuration file (see examples below)
2. Run the tool:
```bash
# Use specific config file
go run main.go structure.yaml

# Or let the tool auto-detect default files
go run main.go
```

## Configuration File Examples

### JSON
```json
{
  "my-project": {
    "src": {
      "main.go": "package main\n\nfunc main() {}",
      "utils": {
        "@helper.go": "package utils"
      }
    },
    "README.md": "# Project Documentation"
  }
}
```

### YAML
```yaml
my-project:
  src:
    main.go: |
      package main
      
      func main() {}
    utils:
      @helper.go: package utils
  README.md: "# Project Documentation"
```

### XML
```xml
<structure>
  <my-project>
    <src>
      <main.go>package main\n\nfunc main() {}</main.go>
      <utils>
        <helper.go>package utils</helper.go>
      </utils>
    </src>
    <README.md># Project Documentation</README.md>
  </my-project>
</structure>
```

## Behavior Notes
1. **Directory Detection**: Any map without `@` prefixed keys will be treated as a directory
2. **File Creation**: Nodes containing `@` prefixed keys will be created as files
3. **Content Handling**: 
   - Empty values will create empty files
   - All values are treated as string content
4. **Auto-Directory Creation**: Parent directories for files are automatically created

## Output Example
```bash
$ go run main.go
📂 Format: YAML
✅ Structure 'my-project' created successfully!
```

## Error Handling
- Will display error messages in case of:
  - Missing configuration files
  - Invalid file formats
  - Permission issues
  - Invalid structure definitions

## License
MIT License (see LICENSE file for details)
```
