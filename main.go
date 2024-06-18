package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	log.SetFlags(log.Ltime | log.Lshortfile) // ログの出力書式を設定する
	if len(os.Args) < 2 {
		fmt.Println("Usage: tochhc2html <path_to_toc.hhc>")
		return
	}

	filePath := os.Args[1]
	nodes := parseHTML(filePath)

	if len(nodes) == 0 {
		fmt.Println("No nodes found")
	}

	summary := generateSummary(nodes)

	//os.WriteFile("tmp.txt", []byte(fmt.Sprintf("%v", hhc)), 0644)

	summary = "# Summary\n\n" + summary

	err := os.WriteFile("SUMMARY.md", []byte(summary), 0644)
	if err != nil {
		panic(fmt.Sprintf("Error writing SUMMARY.md file: %v", err))
	}

	fmt.Println("SUMMARY.md has been generated successfully.")
}

// Node represents a structure for holding the data
type Node struct {
	Name  string
	Local string
	Depth int
}

// parseHTML reads the HTML file and returns a slice of Nodes
func parseHTML(filePath string) []Node {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	htmlContent := strings.ToLower(string(file))
	htmlContent = strings.ReplaceAll(htmlContent, "\n", "")
	htmlContent = strings.ReplaceAll(htmlContent, "\r", "")
	htmlContent = strings.ReplaceAll(htmlContent, "\t", "")
	htmlContent = strings.ReplaceAll(htmlContent, " ", "")

	var nodes []Node
	depth := 0
	pos := 0

	for {
		startUL := strings.Index(htmlContent[pos:], "<ul>")
		endUL := strings.Index(htmlContent[pos:], "</ul>")
		startObject := strings.Index(htmlContent[pos:], "<object")

		if startUL == -1 && endUL == -1 && startObject == -1 {
			break
		}

		if startUL != -1 && (startUL < endUL || endUL == -1) && (startUL < startObject || startObject == -1) {
			depth++
			pos += startUL + len("<ul>")
		} else if endUL != -1 && (endUL < startUL || startUL == -1) {
			depth--
			pos += endUL + len("</ul>")
		} else if startObject != -1 {
			endObject := strings.Index(htmlContent[pos:], "</object>") + pos + len("</object>")
			if endObject == -1 {
				break
			}

			objectContent := htmlContent[pos+startObject : endObject]
			name := extractParam(objectContent, "name")
			local := extractParam(objectContent, "local")

			if name != "" && local != "" {
				node := Node{Name: name, Local: local, Depth: depth}
				nodes = append(nodes, node)
				//log.Printf("Depth: %d, Name: %s, Local: %s\n", depth, node.Name, node.Local)
			}
			pos = endObject
		}
	}

	return nodes
}

// extractParam extracts the value of a param from the object tag content
func extractParam(content, param string) string {
	start := strings.Index(content, fmt.Sprintf(`name="%s"value="`, param))
	if start == -1 {
		return ""
	}
	start += len(fmt.Sprintf(`name="%s"value="`, param))
	end := strings.Index(content[start:], `"`)
	if end == -1 {
		return ""
	}
	return content[start : start+end]
}

// generateSummary generates the SUMMARY.md content from the Node slice
func generateSummary(nodes []Node) string {
	var sb strings.Builder
	for _, node := range nodes {
		sb.WriteString(strings.Repeat("  ", node.Depth))
		sb.WriteString(fmt.Sprintf("- [%s](%s)\n", node.Name, node.Local))
	}
	return sb.String()
}
