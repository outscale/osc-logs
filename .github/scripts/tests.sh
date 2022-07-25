#!/bin/env bash
set -eu

mkdir ~/.osc
cat > ~/.osc/config.json << EOF 
{
    "default": {
        "access_key": "${OSC_ACCESS_KEY}",
        "secret_key": "${OSC_SECRET_KEY}",
        "host": "outscale.com",
        "https": true,
        "method": "POST",
        "region_name": "${OSC_REGION}"
    }
}
EOF

./osc-logs --write="logs.json" --interval=2 --count=2

if [ -s "logs.json" ];then
  echo "log has some entries"
else
  echo "log has no entries"
  exit 1
fi

jq -s -r '(.[] ) ' logs.json