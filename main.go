package main

// Go provides a `flag` package supporting basic
// command-line flag parsing. We'll use this package to
// implement our example command-line program.
import (
	"bytes"
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)
import "fmt"

func isOnlyWhitespace(s string) bool {
	// Iterate over each rune in the string
	for _, r := range s {
		if !unicode.IsSpace(r) {
			return false // Found a non-whitespace character
		}
	}
	return true // All characters are whitespace
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
			fmt.Println("\033[31mError creating .temp directory:\033[0m", err)
			return
		}
	}

	// Step 2: Compile the C++ code to the .temp directory
	outputBinary := filepath.Join(tempDir, name)
	compileCmd := exec.Command("clang++", fmt.Sprintf("%s.cpp", name), "-o", outputBinary, "-std="+std)
	var compileErr bytes.Buffer
	compileCmd.Stderr = &compileErr
	startCompileTime := time.Now()
	err := compileCmd.Run()
	compileDuration := time.Since(startCompileTime)
	
	if err != nil {
		fmt.Println("\033[31mCompilation Error:\033[0m", compileErr.String())
		return
	}


	inputFile := fmt.Sprintf("%s-in.txt", name)
	outputFile := fmt.Sprintf("%s-out.txt", name)

	inputData, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Println("\033[31mError reading input file:\033[0m", err)
		return
	}

	expectedOutput, err := os.ReadFile(outputFile)
	if err != nil {
		fmt.Println("\033[31mError reading output file:\033[0m", err)
		return
	}

	if isOnlyWhitespace(inputData) {
		fmt.Println("Please fill your input testcases, and rerun")
		return
	}

	if isOnlyWhitespace(expectedOutput) {
		fmt.Println("Please fill your expected outputs, and rerun")
		return
	}

	// Step 4: Run the program with a 4s timeout
	cmd := exec.Command(outputBinary)
	cmd.Stdin = bytes.NewReader(inputData)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Start timing the execution
	start := time.Now()
	err = cmd.Start()
	if err != nil {
		fmt.Println("\033[31mRuntime Error:\033[0m", err)
		return
	}

	// Run with timeout
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(45 * time.Second):
		// Time limit exceeded, but we'll let the program continue
		fmt.Println("\033[33mWARNING: Your code took a long time to execute. \033[0m")
		_ = cmd.Process.Kill()
		return
	case err := <-done:
		runtimeDuration := time.Since(start)
		// If there's a runtime error
		if err != nil {
			fmt.Println("\033[31mRuntime Error:\033[0m", stderr.String())
			fmt.Println("Still checking for correctness of the code")
		}
	}

	// Step 5: Compare output line by line
	actualOutput := strings.Split(stdout.String(), "\n")
	expectedLines := strings.Split(string(expectedOutput), "\n")

	fmt.Println("Your actual output: ")
	fmt.Println(actualOutput)

	mismatch := false
	for i := 0; i < len(expectedLines); i++ {
		if i >= len(actualOutput) || strings.TrimSpace(actualOutput[i]) != strings.TrimSpace(expectedLines[i]) {
			if i < len(actualOutput) {
				fmt.Printf("\033[31mMismatch at line %d:\033[0m\nExpected: %s\nActual: %s\n",
					i+1, expectedLines[i], actualOutput[i])
			} else {
				fmt.Printf("\033[31mMismatch at line %d:\033[0m\nExpected: %s\nActual: <no output>\n",
					i+1, expectedLines[i])
			}
			mismatch = true
		}
	}

	fmt.Println("\033[33mOK\033[0m")
	fmt.Println("Took", runtimeDuration)
	
	if !mismatch {
		fmt.Println("\033[32mSuccess: Output matches expected result!\033[0m")
	}

	err = os.Remove(outputBinary)
	if err != nil {
		fmt.Println("\033[31mError removing binary:\033[0m", err)
	}

}

func genCode(problemName *string, cppTemplate *string) {

	fmt.Println("Generating files...")
	if *problemName == "" {
		fmt.Println("You must provide a problem name")
		return
	}
	file, err := os.Create(*problemName + ".cpp")
	file1, err := os.Create(*problemName + "-out.txt")
	file2, err := os.Create(*problemName + "-in.txt")
	if err != nil {
		fmt.Printf("Failed to create problem file: %s\n", err.Error())
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Failed to close problem file: " + err.Error())
		}
		err = file1.Close()
		if err != nil {
			fmt.Printf("Failed to close problem file: %s\n", err.Error())
		}
		err = file2.Close()
		if err != nil {
			fmt.Println("Failed to close problem file: " + err.Error())
		}
	}(file)
	_, err = file.WriteString(string(*cppTemplate))
	fmt.Println("Done.")
}

func main() {
	//wordPtr := flag.String("word", "foo", "a string")
	problemName := flag.String("name", "", "a problem name")
	isGenStage := flag.Bool("new", false, "a bool")
	isTester := flag.Bool("test", false, "a bool")
	templateName := flag.String("t", "default", "a template name")
	std := flag.String("std", "c++20", "C++ standard, c++17, c++20, c++23")
	flag.Parse()

	if *isGenStage == true {
		execPath, err := os.Executable()
		if err != nil {
			fmt.Println("\033[31mError getting executable:\033[0m", err)

		}
		execDirPath := filepath.Dir(execPath)
		filePath := filepath.Join(execDirPath, "templates", *templateName+".cpp")
		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Println("Error readin the template", templateName)
			println("\033[31mError reading template file:\033[0m", err)
			return
		}
		cppTemplate := string(content)
		if fileExists(*problemName+".cpp") || fileExists(*problemName+"-out.txt") || fileExists(*problemName+"-in.txt") {
			fmt.Println("the files were already generated for", *problemName)
			return
		}

		genCode(problemName, &cppTemplate)
	} else if *isTester == true {
		if *problemName == "" {
			fmt.Println("You must provide a problem name")
			return
		}
		fmt.Println("Testing ...")

		evaluateCode(*problemName, std)

		fmt.Println("Test complete")
	}
}
