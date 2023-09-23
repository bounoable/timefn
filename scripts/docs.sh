#!/bin/sh
set -e

ROOT=$(git rev-parse --show-toplevel)

if ! command -v jotbot &> /dev/null; then
	echo "jotbot not found in PATH"
	echo
	echo "To install JotBot, run"
	echo
	echo "  go install github.com/modernice/jotbot/cmd/jotbot"
	echo
	echo "Then run this script again."
	exit 1
fi

if [ ! -f "$ROOT/jotbot.env" ]; then
	echo "jotbot.env not found"
	echo
	echo "To create jotbot.env, run"
	echo
	echo "  cp jotbot.env.example jotbot.env"
	echo
	echo "Then fill in the OPENAI_API_KEY and run this script again."
	exit 1
fi

set -a
. "$ROOT/jotbot.env"
set +a

if [ -z "$OPENAI_API_KEY" ]; then
	echo "OPENAI_API_KEY not set"
	echo
	echo "To set OPENAI_API_KEY, edit jotbot.env and run this script again."
	exit 1
fi

jotbot generate "$ROOT"
