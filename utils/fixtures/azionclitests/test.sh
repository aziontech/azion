#!/bin/bash

# Shell script for testing the md/azion/main.go link and build commands

# Set the path to the main.go file
main_go_path="cmd/azion/main.go"
# Set the path to the expected folder after the link command
expected_folder="azion"
# Set the path to the expected folder after the build command
expected_edge_folder=".edge"

# Function to check if a folder exists
check_folder_exists() {
    local folder_path="$1"
    if [ -d "$folder_path" ]; then
        echo "Folder $folder_path exists."
    else
        echo "Error: Folder $folder_path not found."
        exit 1
    fi
}

# Check if the main.go file exists
if [ -f "$main_go_path" ]; then
    # Run the link command with the specified options
    echo "Running cmd/azion/main.go link --preset astro --mode deliver --auto"
    go run "$main_go_path" link --preset astro --mode deliver --auto --debug

    # Check the exit status of the last command
    if [ $? -eq 0 ]; then
        echo "Link command executed successfully."
        check_folder_exists "$expected_folder"

        # Run the build command
        echo "Running cmd/azion/main.go build --debug"
        go run "$main_go_path" build

        # Check the exit status of the build command
        if [ $? -eq 0 ]; then
            echo "Build command executed successfully."

            # Check if the expected edge folder exists
            check_folder_exists "$expected_edge_folder"
        else
            echo "Error: An error occurred while running the build command."
            exit 1
        fi
    else
        echo "Error: An error occurred while running the link command."
        exit 1
    fi
else
    echo "Error: cmd/azion/main.go not found. Please make sure the file exists."
    exit 1
fi
