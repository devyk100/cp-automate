# Mochi: C++ Problem Generator and Tester

Mochi is a command-line tool written in Go for generating and testing C++ problem files. It streamlines the process of creating problem files from templates and automates the testing of C++ code against predefined inputs and outputs.

## Requirements

- **Go**: Ensure you have Go installed to run this program.
- **GCC**: The program requires `g++` to compile C++ code.
- **Go Packages**: The program uses the `github.com/fatih/color` package for colored terminal output. You can install it using:
  ```bash
  go get github.com/fatih/color
  ```

## Installation

To install Mochi on a Linux system, you can use the provided `install.sh` script. This script will compile the Go project, move the executable and templates to the installation directory, and add the directory to your PATH.

### Steps to Install

1. Run the `install.sh` script:
   ```bash
   ./install.sh
   ```

2. Follow the on-screen instructions to complete the installation.

3. Reload your terminal or run the following command to update your PATH:
   ```bash
   export PATH="$PATH:$HOME/cp-automate"
   ```

4. Verify the installation by running:
   ```bash
   mochi -h
   ```

If the installation is successful, you should see the help message for Mochi.


First, add the executable to your PATH to use it conveniently as `mochi`. You can then run the program with the following command-line flags:

- `-name`: Specify the problem name.
- `-new`: Generate new problem files using a template.
- `-test`: Test the C++ code against input and output files.
- `-template`: Specify the template name (default is `default`).
- `-std`: Specify the C++ standard (default is `c++20`).

### Generating Problem Files

To generate new problem files, use the `-new` flag along with the `-name` and `-template` flags:

```bash
mochi -new -name=ProblemName -template=TemplateName
```

This will create the following files:
- `ProblemName.cpp`
- `ProblemName-in.txt`
- `ProblemName-out.txt`

### Testing C++ Code

To test the C++ code, use the `-test` flag along with the `-name` flag:

```bash
mochi -test -name=ProblemName
```

This will compile the C++ code and compare its output against the expected output.

## Templates

C++ templates are stored in the `templates` directory. You can add your own templates or modify existing ones to suit your needs.
