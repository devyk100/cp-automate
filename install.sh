#!/bin/bash

# Set variables
EXECUTABLE_NAME="mochi"       # The name of the executable
INSTALL_DIR="$HOME/cp-automate"  # Directory to place the executable and templates
EXECUTABLE_PATH="$INSTALL_DIR/$EXECUTABLE_NAME"
TEMPLATES_DIR="$INSTALL_DIR/templates"
ZSHRC="$HOME/.zshrc"          # Assuming the user is using Zsh (change to .bashrc if using bash)

# Step 1: Compile the Go project
echo "Compiling the Go project..."
go build -o "$EXECUTABLE_NAME" .

# Check if the build was successful
if [ ! -f "$EXECUTABLE_NAME" ]; then
    echo "Error: Failed to compile the Go project."
    exit 1
fi

# Step 2: Create the installation directory if it doesn't exist
if [ ! -d "$INSTALL_DIR" ]; then
    echo "Creating directory: $INSTALL_DIR"
    mkdir -p "$INSTALL_DIR"
fi

# Step 3: Move the executable to the installation directory
echo "Moving $EXECUTABLE_NAME to $INSTALL_DIR"
mv "$EXECUTABLE_NAME" "$INSTALL_DIR/"

# Step 4: Move the templates directory to the installation directory
if [ -d "./templates" ]; then
    echo "Moving templates/ directory to $INSTALL_DIR"
    cp -r ./templates "$INSTALL_DIR/"
else
    echo "Warning: No templates directory found."
fi

# Step 5: Add the installation directory to the PATH in .zshrc (or .bashrc)
if ! grep -q "$INSTALL_DIR" "$ZSHRC"; then
    echo "Adding $INSTALL_DIR to the PATH in $ZSHRC"
    echo "export PATH=\"\$PATH:$INSTALL_DIR\"" >> "$ZSHRC"
else
    echo "$INSTALL_DIR is already in the PATH."
fi

# Step 6: Reload .zshrc (or .bashrc) to apply changes
echo "Reload the terminal"
echo "Or do `export PATH=\"\$PATH:$INSTALL_DIR\"`"

# Step 7: Verify the installation
echo "Verifying the installation..."
if command -v "$EXECUTABLE_NAME" &> /dev/null; then
    echo "$EXECUTABLE_NAME is successfully installed and accessible from the PATH!"
else
    echo "Error: $EXECUTABLE_NAME is not found in the PATH."
    exit 1
fi
