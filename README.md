# Mochi: C++ Problem Generator and Tester

Mochi is a Go program that facilitates the generation and testing of C++ problem files. It allows users to create new problem files from templates and test C++ code against specified input and output files.

## Requirements

- **Go**: Ensure you have Go installed to run this program.
- **GCC**: The program requires `g++` to compile C++ code.
- **Go Packages**: The program uses the `github.com/fatih/color` package for colored terminal output. You can install it using:
  ```bash
  go get github.com/fatih/color
  ```

## Usage

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
