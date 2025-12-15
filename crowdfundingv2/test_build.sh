#!/bin/bash
cd ~/crowdfunding/crowdfundingv2/contracts
echo "Testing Go build..."
go build -v -o crowdfunding 2>&1
if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
    ls -lh crowdfunding
    rm crowdfunding
else
    echo "❌ Build failed!"
fi
