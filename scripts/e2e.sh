#!/bin/bash

# Shell script for testing the entire process from link to deploy, against the real production environment

# save current directory
current_directory=$(pwd)
# Set the path to the main.go file
main_go_path="cmd/azion/main.go"
# use bin_path instead of main_go_path, if you want to run the script using a built binary file
bin_path="bin/azion" 
# Combine the current_directory and main_go_path
full_main_go_path="${current_directory}/${main_go_path}"
full_bin_path="${current_directory}/${bin_path}"
# Set the path to the expected folder after the link command
expected_folder="azion"
# Set the path to the expected folder after the build command
expected_edge_folder=".edge"
# token from secrets
token=$TOKEN

# Function to list and delete edge storage buckets
delete_edge_storage_buckets() {
  # List the buckets with page size 1000, skip the header line, and delete each bucket
  "$full_bin_path" list edge-storage bucket --page-size 1000 | awk 'NR>1 {print $1}' | while read -r name; do
    "$full_bin_path" delete edge-storage bucket --name "$name" --force
  done
}

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

check_azion_json() {
    json_file="${current_directory}/react-static/azion/azion.json"

    # Extract the application_id from the JSON file
    local application_id
    application_id=$(jq -r '.application.id' "$json_file")

    # Check if the application_id is a valid number
    if ! [[ "$application_id" =~ ^[0-9]+$ ]]; then
        echo "Error: Invalid application_id value: $application_id"
        exit 1
    fi

    # Check if application_id is greater than zero
    if (( application_id > 0 )); then
        echo "Application ID retrieved, is valid and greater than zero: $application_id"
        return 0
    else
        echo "Error: azion.json was not retrieved successfully"
        exit 1
    fi
}

make build
"$full_bin_path" -t "$token"
# go run "$full_main_go_path" -t "$token"
delete_edge_storage_buckets


# Check if the main.go file exists
if [ -f "$full_main_go_path" ]; then
    echo "Current Directory: $current_directory"

    rm -rf vulcan
    rm -rf react-static

    git clone https://github.com/aziontech/vulcan-examples.git
    cp -r vulcan-examples/examples/react-static ./ 
    cd react-static

    npm install
    # Run the link command with the specified options
    echo "Running cmd/azion/main.go link --preset react --auto --debug"
    # go run "$full_main_go_path" link --preset react --auto --debug
    # uncomment the line below and comment the line above if using a built binary file
    "$full_bin_path" link --preset react --auto --debug

    # Check the exit status of the last command
    if [ $? -eq 0 ]; then
        echo "Link command executed successfully."
        check_folder_exists "$expected_folder"

        # Run the build command
        echo "Running cmd/azion/main.go build --debug"
        # go run "$full_main_go_path" build --debug
        # uncomment the line below and comment the line above if using a built binary file
        "$full_bin_path" build --debug

        # Check the exit status of the build command
        if [ $? -eq 0 ]; then
            echo "Build command executed successfully."

            # Check if the expected edge folder exists
            check_folder_exists "$expected_edge_folder"

             # Run the build command
            echo "Running cmd/azion/main.go deploy --debug"
            # go run "$full_main_go_path" deploy --debug
            # uncomment the line below and comment the line above if using a built binary file
            "$full_bin_path" deploy --debug

            # Check the exit status of the deploy command
            if [ $? -eq 0 ]; then
                echo "Deploy command executed successfully."

                check_azion_json

            else
                echo "Error: An error occurred while running the build command."
                exit 1
            fi

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
