#!/bin/bash

echo "🚀 Starting airport deployment..."
sleep 2

# Check if the Airport service exists
if ! systemctl list-units --type=service | grep -q 'Airport.service'; then
    echo "🛑 Airport service not found. Creating Airport.service..."
    sudo bash -c 'cat > /etc/systemd/system/Airport.service << EOF
[Unit]
Description=Airport Service

[Service]
ExecStart=/usr/local/bin/main
WorkingDirectory=/usr/local/bin
Restart=always
User=root
EnvironmentFile=/usr/local/bin/.env

[Install]
WantedBy=multi-user.target
EOF'
    if [ $? -ne 0 ]; then
        echo "❌ Failed to create the service file!"
        exit 1
    fi
    echo "✅ Airport service created successfully."
else
    echo "✅ Airport service already exists. Continuing..."
fi

if [ -f /usr/local/bin/main ]; then
    echo "🗑️ Removing the old main file..."
    sudo rm /usr/local/bin/main
    if [ $? -ne 0 ]; then
        echo "❌ Failed to remove the old file!"
        exit 1
    fi
else
    echo "🔔 No old main file found. Skipping removal..."
fi

if [ -f .env ]; then
    echo "📂 Copying .env file to /usr/local/bin..."
    sudo cp .env /usr/local/bin/.env
    if [ $? -ne 0 ]; then
        echo "❌ Failed to copy .env file!"
        exit 1
    fi
else
    echo "⚠️ No .env file found in the current directory!"
fi

echo "⚙️ Building the Go file for Linux..."
sleep 2

GOOS=linux GOARCH=amd64 go build main.go
if [ $? -ne 0 ]; then
    echo "❌ Failed to build the Go file for Linux!"
    exit 1
fi

echo "📂 Moving the new build file..."
sleep 2
sudo mv main /usr/local/bin/
if [ $? -ne 0 ]; then
    echo "❌ Failed to move the new file!"
    exit 1
fi

echo "🔌 Enabling the service..."
sudo systemctl enable Airport.service

echo "▶️ Starting the service..."
sleep 2
sudo systemctl start Airport.service
if [ $? -ne 0 ]; then
    echo "❌ Failed to start the service!"
    exit 1
fi

echo "📊 Checking the service status..."
sleep 2
sudo systemctl status Airport.service

echo "✅ Web deployment completed successfully! 🎉"
echo "👨‍💻 Created by tmsoft12. 🎉 Enjoy! 😎"
