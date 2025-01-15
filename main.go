package main

// Go provides a `flag` package supporting basic
// command-line flag parsing. We'll use this package to
// implement our example command-line program.
import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/fatih/color"
)

var (
	red    = color.New(color.FgRed).Add(color.Underline)
	green  = color.New(color.FgGreen).Add(color.Bold)
	orange = color.New(color.FgHiWhite, color.BgBlack).Add(color.Bold)
	blue   = color.New(color.FgBlue)
	yellow = color.New(color.FgYellow).Add(color.Italic)
)

func isOnlyWhitespace(s string) bool {
	for _, r := range s {
		if !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)

	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}

	return false
}

func evaluateCode(name string, std string) {
	// Step 1: Create .temp directory
	tempDir := ".temp"
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		err = os.Mkdir(tempDir, 0755)
		if err != nil {
			red.Println("Error creating .temp directory:", err)
			return
		}
	}

	// Step 2: Compile the C++ code to the .temp directory
	outputBinary := filepath.Join(tempDir, name)
	compileCmd := exec.Command("g++", "-D", "YASH_DEBUG=fsaf", fmt.Sprintf("%s.cpp", name), "-o", outputBinary, "-std="+std)
	var compileErr bytes.Buffer
	compileCmd.Stderr = &compileErr
	blue.Println("Using", std, "standard to compile.")
	startCompileTime := time.Now()
	err := compileCmd.Run()
	compileDuration := time.Since(startCompileTime)
	blue.Println("Took", compileDuration, "to compile.")

	if err != nil {
		red.Println("Compilation Error:", compileErr.String())
		return
	}

	inputFile := fmt.Sprintf("%s-in.txt", name)
	outputFile := fmt.Sprintf("%s-out.txt", name)

	inputData, err := os.ReadFile(inputFile)
	if err != nil {
		red.Println("Error reading input file:", err)
		return
	}

	expectedOutput, err := os.ReadFile(outputFile)
	if err != nil {
		red.Println("Error reading output file:", err)
		return
	}

	if isOnlyWhitespace(string(inputData)) {
		yellow.Println("Please fill your input testcases, and rerun")
		return
	}

	if isOnlyWhitespace(string(expectedOutput)) {
		yellow.Println("Please fill your expected outputs, and rerun")
		return
	}

	// Step 4: Run the program with a 4s timeout
	cmd := exec.Command(outputBinary)
	cmd.Stdin = bytes.NewReader(inputData)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	start := time.Now()
	err = cmd.Start()
	if err != nil {
		red.Println("Runtime Error:", err)
		return
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()
	var runtimeDuration time.Duration
	select {
	case <-time.After(60 * time.Second):
		red.Println("WARNING: Your code took a long time to execute. ")
		_ = cmd.Process.Kill()
		return
	case err := <-done:
		runtimeDuration = time.Since(start)
		if err != nil {
			red.Println("Runtime Error:", stderr.String())
			yellow.Println("Still checking for correctness of the code")
		}
	}

	actualOutput := strings.Split(stdout.String(), "\n")
	expectedLines := strings.Split(string(expectedOutput), "\n")

	blue.Println("\nYour actual output: ")
	yellow.Println(actualOutput)

	mismatch := false
	for i := 0; i < len(expectedLines); i++ {
		if i >= len(actualOutput) || strings.TrimSpace(actualOutput[i]) != strings.TrimSpace(expectedLines[i]) {
			if !mismatch {
				red.Println("\nMismatches:")
			}
			mismatch = true
			if i < len(actualOutput) {
				yellow.Printf("Mismatch at line %d:\n", i+1)
				orange.Printf("Expected: %s\nActual: %s\n", expectedLines[i], actualOutput[i])
			} else {
				yellow.Printf("Mismatch at line %d:\n", i+1)
				orange.Printf("Expected: %s\nActual: <no output>\n", expectedLines[i])
			}
		}
	}

	blue.Println("\nOK")
	blue.Println("Took", runtimeDuration)

	if !mismatch {
		green.Println("Success: Output matches expected result!")
	}

	err = os.Remove(outputBinary)
	if err != nil {
		red.Println("Error removing binary:", err)
	}

}

func genCode(problemName *string, cppTemplate *string) {

	blue.Println("Generating files...")
	if *problemName == "" {
		yellow.Println("You must provide a problem name")
		return
	}
	file, err := os.Create(*problemName + ".cpp")
	if err != nil {
		red.Printf("Failed to create problem file: %s\n", err.Error())
	}

	file1, err := os.Create(*problemName + "-out.txt")
	if err != nil {
		red.Printf("Failed to create problem file: %s\n", err.Error())
	}

	file2, err := os.Create(*problemName + "-in.txt")
	if err != nil {
		red.Printf("Failed to create problem file: %s\n", err.Error())
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			yellow.Println("Failed to close problem file: " + err.Error())
		}
		err = file1.Close()
		if err != nil {
			yellow.Printf("Failed to close problem file: %s\n", err.Error())
		}
		err = file2.Close()
		if err != nil {
			yellow.Println("Failed to close problem file: " + err.Error())
		}
	}(file)
	_, err = file.WriteString(string(*cppTemplate))
	if err != nil {
		fmt.Println(err.Error())
	}
	green.Println("Done.")
}

func main() {
	problemName := flag.String("name", "", "a problem name")
	isGenStage := flag.Bool("new", false, "a bool")
	isTester := flag.Bool("test", false, "a bool")
	templateName := flag.String("template", "default", "a template name")
	std := flag.String("std", "c++20", "C++ standard, c++17, c++20, c++23")
	flag.Parse()

	if *isGenStage {
		execPath, err := os.Executable()
		if err != nil {
			red.Println("Error getting executable:", err)
		}
		execDirPath := filepath.Dir(execPath)
		filePath := filepath.Join(execDirPath, "templates", *templateName+".cpp")
		content, err := os.ReadFile(filePath)
		if err != nil {
			red.Println("Error readin the template", templateName)
			red.Println("Error:", err)
			return
		}
		cppTemplate := string(content)
		if fileExists(*problemName+".cpp") || fileExists(*problemName+"-out.txt") || fileExists(*problemName+"-in.txt") {
			yellow.Println("the files were already generated for", *problemName)
			return
		}

		genCode(problemName, &cppTemplate)
	} else if *isTester {
		if *problemName == "" {
			yellow.Println("You must provide a problem name")
			return
		}
		blue.Println("Testing ...")

		evaluateCode(*problemName, *std)

		blue.Println("Test complete")
	} else {
		yellow.Println("No commands specified")
	}
}
