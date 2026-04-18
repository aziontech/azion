#!/usr/bin/env python3
"""
Azion CLI Configuration Generator

Generates azion.json and args.json templates for Azion projects.
Run this script to create configuration files for your project.

Usage:
    python generate_config.py --type <preset> [--name <project-name>] [--output <dir>]

Examples:
    python generate_config.py --type nextjs --name my-app
    python generate_config.py --type react --output ./config
    python generate_config.py --type static --name static-site --output ./deploy
"""

import argparse
import json
import os
from pathlib import Path

# Framework presets with their default configurations
PRESETS = {
    "nextjs": {
        "preset": "next",
        "build_command": "npm run build",
        "output_dir": ".next",
        "runtime": "edge"
    },
    "react": {
        "preset": "react",
        "build_command": "npm run build",
        "output_dir": "dist",
        "runtime": "edge"
    },
    "vue": {
        "preset": "vue",
        "build_command": "npm run build",
        "output_dir": "dist",
        "runtime": "edge"
    },
    "angular": {
        "preset": "angular",
        "build_command": "npm run build",
        "output_dir": "dist",
        "runtime": "edge"
    },
    "astro": {
        "preset": "astro",
        "build_command": "npm run build",
        "output_dir": "dist",
        "runtime": "edge"
    },
    "hexo": {
        "preset": "hexo",
        "build_command": "npm run build",
        "output_dir": "public",
        "runtime": "edge"
    },
    "static": {
        "preset": "static",
        "build_command": None,
        "output_dir": "public",
        "runtime": "edge"
    },
    "javascript": {
        "preset": "javascript",
        "build_command": None,
        "output_dir": "dist",
        "runtime": "edge"
    },
    "typescript": {
        "preset": "typescript",
        "build_command": None,
        "output_dir": "dist",
        "runtime": "edge"
    }
}


def generate_azion_json(name: str, preset: str, output_dir: str) -> dict:
    """Generate azion.json configuration."""
    preset_config = PRESETS.get(preset, PRESETS["static"])

    config = {
        "id": None,
        "name": name,
        "preset": preset_config["preset"],
        "runtime": preset_config["runtime"],
        "active": True,
        "application": {
            "delivery_protocol": "https",
            "http3": True,
            "http_port": 80,
            "https_port": 443,
            "origin": []
        },
        "functions": [],
        "cache_settings": {
            "browser_cache_settings": "honor",
            "browser_cache_settings_maximum_ttl": 0,
            "cdn_cache_settings": "honor",
            "cdn_cache_settings_maximum_ttl": 3600
        },
        "rules_engine": {
            "request": [],
            "response": []
        }
    }

    return config


def generate_args_json(preset: str) -> dict:
    """Generate args.json configuration."""
    preset_config = PRESETS.get(preset, PRESETS["static"])

    config = {
        "preset": preset_config["preset"],
        "mode": "deliver",
        "build_command": preset_config["build_command"],
        "output_dir": preset_config["output_dir"],
        "env_vars": [],
        "secrets": []
    }

    return config


def main():
    parser = argparse.ArgumentParser(
        description="Generate Azion CLI configuration files"
    )
    parser.add_argument(
        "--type",
        choices=list(PRESETS.keys()),
        required=True,
        help="Framework preset type"
    )
    parser.add_argument(
        "--name",
        default="my-application",
        help="Application name (default: my-application)"
    )
    parser.add_argument(
        "--output",
        default=".",
        help="Output directory (default: current directory)"
    )

    args = parser.parse_args()

    # Create output directory if needed
    output_path = Path(args.output)
    output_path.mkdir(parents=True, exist_ok=True)

    # Generate configurations
    azion_config = generate_azion_json(args.name, args.type, args.output)
    args_config = generate_args_json(args.type)

    # Write azion.json
    azion_path = output_path / "azion.json"
    with open(azion_path, "w") as f:
        json.dump(azion_config, f, indent=2)
    print(f"Created: {azion_path}")

    # Write args.json
    args_path = output_path / "args.json"
    with open(args_path, "w") as f:
        json.dump(args_config, f, indent=2)
    print(f"Created: {args_path}")

    # Create azion directory structure
    azion_dir = output_path / "azion"
    azion_dir.mkdir(exist_ok=True)
    print(f"Created: {azion_dir}")

    print(f"\nConfiguration generated for {args.type} project '{args.name}'")
    print(f"\nNext steps:")
    print(f"  1. Edit azion.json to configure your application")
    print(f"  2. Run: azion config apply")
    print(f"  3. Build and deploy: azion build && azion deploy")


if __name__ == "__main__":
    main()