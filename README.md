# Confluence Exporter

Confluence Exporter is a CLI tool designed to export content from Confluence into Markdown format. This project provides a simple way to retrieve pages and their content from Confluence and convert them into a more portable format.

## Features

- Fetch pages from Confluence
- Convert Confluence content to Markdown
- Easy configuration through environment variables or config files

## Project Structure

```
confluence-exporter
├── cmd
│   └── exporter
│       └── main.go          # Entry point of the CLI application
├── internal
│   ├── api
│   │   └── confluence.go    # Functions to interact with the Confluence API
│   ├── converter
│   │   └── markdown.go      # Convert Confluence content to Markdown
│   ├── config
│   │   └── config.go        # Configuration settings for the application
│   └── models
│       └── page.go          # Data structures for Confluence pages
├── pkg
│   └── utils
│       ├── auth.go          # Utility functions for authentication
│       └── logger.go        # Logging functionality
├── go.mod                    # Module definition and dependencies
├── go.sum                    # Checksums for module dependencies
└── README.md                 # Project documentation
```

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/Ilhicas/confluence-exporter.git
   cd confluence-exporter
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

## Usage

To run the application, use the following command:

```
go run cmd/exporter/main.go --config config.json
```

Replace `<path_to_config_file>` with the path to your configuration file containing the necessary API credentials.

## License

This project is licensed under the MIT License. See the LICENSE file for details.# confluence-exporter

## Disclaimer

The project is still in development and may not be fully functional. Use at your own risk.
This project was generated with AI assistance from Claude 3.7 Sonnet Thinking model.
