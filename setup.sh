#!/bin/bash

# Place all your setup commands here

# Clone the repository
git clone https://github.com/seal/auto-build
cd auto-build

# Ask for config details from the user
echo "Enter GitHub Token:"
read githubToken
echo "Enter Repo Owner:"
read repoOwner
echo "Enter Repo Name:"
read repoName
echo "Enter Code Directory:"
read codeDirectory
echo "Enter Code Run Command:"
read codeRun

# Create the config.json file
echo '{
  "githubToken": "'$githubToken'",
  "repoOwner": "'$repoOwner'",
  "repoName": "'$repoName'",
  "codeDirectory": "'$codeDirectory'",
  "codeRun": "'$codeRun'"
}' > config.json

# Build and run the Go program
go build -o myprogram
./myprogram &

# Install the Go program as a systemctl service
# (Note: You might need appropriate permissions to install as a systemctl service)
#sudo cp myprogram /usr/local/bin/
#sudo cp myprogram.service /etc/systemd/system/
#sudo systemctl enable myprogram.service
#sudo systemctl start myprogram.service

echo "Setup complete. Start program in new tmux session."



