package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "io/ioutil"
    "os"
)

type TreeNode struct {
    Key      string      
    Value    interface{} 
    IsLeaf   bool        
    Children []*TreeNode
}

func buildTree(key string, value interface{}) *TreeNode {
    node := &TreeNode{
        Key: key,
    }

    switch v := value.(type) {
    case map[string]interface{}:
        for k, val := range v {
            child := buildTree(k, val)
            node.Children = append(node.Children, child)
        }
    case []interface{}:
        for i, val := range v {
            childKey := fmt.Sprintf("[%d]", i)
            child := buildTree(childKey, val)
            node.Children = append(node.Children, child)
        }
    default:
        node.Value = v
        node.IsLeaf = true
        return node
    }
    return node
}

func printTree(node *TreeNode, prefix string, isLast bool) {
    connector := ""
    if len(prefix) > 0 {
        if isLast {
            connector = "└── "
        } else {
            connector = "├── "
        }
        prefix += connector
    }

    if node.Key != "" {
        fmt.Print(prefix + node.Key)
        if !node.IsLeaf && len(node.Children) > 0 {
            fmt.Println("")
        } else if node.IsLeaf {
            fmt.Printf(": %v\n", node.Value)
        }
    }

    lastIndex := len(node.Children) - 1
    for idx, child := range node.Children {
        isChildLast := idx == lastIndex
        newPrefix := prefix[:len(prefix)-len(connector)]
        if isLast || idx < lastIndex {
            newPrefix += "│ "
        } else {
            newPrefix += "│ "
        }
        printTree(child, newPrefix, isChildLast)
    }
}


var (
    inputFile string
    treeView  bool
)

func init() {
    flag.StringVar(&inputFile, "i", "", "Input JSON file path")
    flag.BoolVar(&treeView, "t", false, "Enable tree view of JSON structure")
}

func main() {
    flag.Parse()

    if inputFile == "" {
        fmt.Println("Please provide an input file using the '-i' flag.")
        os.Exit(1)
    }

    data, err := ioutil.ReadFile(inputFile)
    if err != nil {
        fmt.Printf("Error reading file %s: %v\n", inputFile, err)
        os.Exit(1)
    }

    var parsedData interface{}
    err = json.Unmarshal(data, &parsedData)
    if err != nil {
        fmt.Printf("Error parsing JSON from %s: %v\n", inputFile, err)
        os.Exit(1)
    }

    root := buildTree("", parsedData)

    if treeView {
        fmt.Printf("%s\n", inputFile)
        printTree(root, "", true)
    } else {
        output, _ := json.MarshalIndent(parsedData, "", "\t")
        fmt.Println(string(output))
    }
}
