package main

// Go provides a `flag` package supporting basic
// command-line flag parsing. We'll use this package to
// implement our example command-line program.
import (
	"bytes"
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)
import "fmt"

func evaluateCode(name string) {
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
	compileCmd := exec.Command("g++", fmt.Sprintf("%s.cpp", name), "-o", outputBinary)
	var compileErr bytes.Buffer
	compileCmd.Stderr = &compileErr
	err := compileCmd.Run()

	// If compilation fails
	if err != nil {
		fmt.Println("\033[31mCompilation Error:\033[0m", compileErr.String())
		return
	}

	// Step 3: Execute the compiled program with input and compare output
	inputFile := fmt.Sprintf("%s-in.txt", name)
	outputFile := fmt.Sprintf("%s-out.txt", name)

	// Read input from the input file
	inputData, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Println("\033[31mError reading input file:\033[0m", err)
		return
	}

	// Read expected output from the output file
	expectedOutput, err := os.ReadFile(outputFile)
	if err != nil {
		fmt.Println("\033[31mError reading output file:\033[0m", err)
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
	case <-time.After(4 * time.Second):
		// Time limit exceeded, but we'll let the program continue
		fmt.Println("\033[33mWARNING: Time Limit Exceeded (TLE), still checking for correctness...\033[0m")
		_ = cmd.Process.Kill()
		return
	case err := <-done:
		// If there's a runtime error
		if err != nil {
			fmt.Println("\033[31mRuntime Error:\033[0m", stderr.String())
			return
		}
	}

	// Step 5: Compare output line by line
	actualOutput := strings.Split(stdout.String(), "\n")
	expectedLines := strings.Split(string(expectedOutput), "\n")

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

	fmt.Println("\033[33mOK\033[0m", time.Since(start))

	// If output doesn't match, report a mismatch
	if !mismatch {
		fmt.Println("\033[32mSuccess: Output matches expected result!\033[0m")
	}

	// Step 6: Cleanup - remove the compiled binary
	err = os.Remove(outputBinary)
	if err != nil {
		fmt.Println("\033[31mError removing binary:\033[0m", err)
	}

}

func main() {
	//wordPtr := flag.String("word", "foo", "a string")
	problemName := flag.String("name", "", "a problem name")
	isGenStage := flag.Bool("gen-problem", false, "a bool")
	isTester := flag.Bool("test", false, "a bool")
	flag.Parse()
	cppTemplate := `#include <bits/stdc++.h>
using namespace std;

template<typename A, typename B> ostream& operator<<(ostream &os, const pair<A, B> &p) { return os << '(' << p.first << ", " << p.second << ')'; }
template<typename T_container, typename T = typename enable_if<!is_same<T_container, string>::value, typename T_container::value_type>::type> ostream& operator<<(ostream &os, const T_container &v) { os << '{'; string sep; for (const T &x : v) os << sep << x, sep = ", "; return os << '}'; }
void dbg_out() { cerr << endl; }
template<typename Head, typename... Tail> void dbg_out(Head H, Tail... T) { cerr << ' ' << H; dbg_out(T...); }
#ifdef LOCAL
#define dbg(...) cerr << "(" << #__VA_ARGS__ << "):", dbg_out(__VA_ARGS__)
#else
#define dbg(...)
#endif

#define ar array
#define ll long long
#define ld long double
#define sza(x) ((int)x.size())
#define all(a) (a).begin(), (a).end()

const int MAX_N = 1e5 + 5;
const ll MOD = 1e9 + 7;
const ll INF = 1e9;
const ld EPS = 1e-9;



void solve() {
    
}

int main() {
    ios_base::sync_with_stdio(0);
    cin.tie(0); cout.tie(0);
    int tc = 1;
    // cin >> tc;
    for (int t = 1; t <= tc; t++) {
        // cout << "Case #" << t << ": ";
        solve();
    }
}
`
	if *isGenStage == true {
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
		_, err = file.WriteString(cppTemplate)
		fmt.Println("Done.")
	} else if *isTester == true {
		if *problemName == "" {
			fmt.Println("You must provide a problem name")
			return
		}
		fmt.Println("Testing ...")

		evaluateCode(*problemName)

		fmt.Println("Test complete")
	}
}
