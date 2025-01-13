# compSudoku

`compSudoku` is a web-based Sudoku game with real-time chat and customizable difficulty settings. It is built using Go, Echo framework, and templ for efficient rendering and real-time interactions.

## Features

- **Interactive Sudoku Game**: Solve Sudoku puzzles of varying difficulty.
- **Real-Time Chat**: Engage with other users via a chat interface.
- **Responsive Design**: Optimized for various screen sizes.
- **Customizable Game Settings**: Generate new games with different difficulty levels.
- **WebSocket Integration**: Real-time updates and interactions.
- **Hot Reloading**: Streamline development with live updates using Golang Air.

## Technologies and How It Works

### Core Technologies
- **Go**: The primary programming language used for backend development.
- **Echo Framework**: A fast and extensible web framework for routing, middleware, and static file serving.
- **templ**: A component-based template engine for generating HTML efficiently.
- **WebSocket**: Enables real-time communication for chat and game updates.
- **htmx**: A lightweight JavaScript library used on the client side to handle HTML responses efficiently.
- **Golang Air**: Hot reloading for Go applications, enabling rapid development by automatically restarting the server on code changes.
- **RapidAPI Sudoku API**: Fetches Sudoku puzzles and solutions based on selected difficulty.

### How It Works

1. **WebSocket Communication**:
   - The server uses WebSockets to manage real-time interactions.
   - Players interact with the game board and chat system without needing to reload the page.
   - WebSocket events trigger updates that are broadcast to all connected clients.

2. **Dynamic HTML Rendering**:
   - The `templ` library is used to create reusable components for UI elements like the Sudoku board, chatbox, and settings.
   - Templates are rendered server-side and sent to the client for efficient updates via `htmx`.

3. **htmx for Efficient Client-Side Updates**:
   - `htmx` enables declarative, event-driven interactions on the client side.
   - It handles server responses that return small HTML fragments, dynamically updating only the affected parts of the page.
   - Combined with WebSockets, `htmx` ensures seamless real-time updates to the game board and chat.

4. **Game State Management**:
   - The game state, including the board layout, active square, and mistakes, is maintained on the server.
   - Server-side logic validates game actions and updates the state accordingly.

5. **API Integration**:
   - The application integrates with RapidAPI's Sudoku API to fetch puzzles dynamically.
   - Difficulty levels are configurable, and users can generate new games via the settings panel.

6. **Hot Reloading with Golang Air**:
   - Golang Air is used during development to monitor changes in the codebase.
   - When a change is detected, the application server restarts automatically, reducing the need for manual restarts and improving development efficiency.

7. **Responsive Design**:
   - The CSS files in the `web/public/css` directory ensure the application adapts to various screen sizes, providing a seamless user experience.

## Project Structure

```
compSudoku/
├── cmd/compSudoku/           # Main entry point
│   └── main.go
├── internal/
│   ├── helper/               # Utility functions
│   ├── models/               # Data models for game and chat
│   └── routes/               # WebSocket routes
├── web/
│   ├── components/           # Templ components for UI
│   └── public/               # Static files (CSS, assets)
├── .gitignore
└── Dockerfile
```

## Getting Started

### Prerequisites

- Go 1.23 or later
- Docker (optional for containerized deployment)
- [Golang Air](https://github.com/cosmtrek/air) for hot reloading
- [Golang Templ](https://templ.guide/) for templating

### Setup

0. Get API Key for RapidSudoku at https://rapidapi.com/andrewarochukwu/api/sudoku-board (free)

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/compSudoku.git
   cd compSudoku
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Install Golang Air && Templ (if not already installed). Ensure they are in your path before continuing:
   ```bash
   go install github.com/cosmtrek/air@latest
   go install github.com/a-h/templ/cmd/templ@latest
   ```

4. Set the toml file for air:
    ```bash
    air -c .air.toml
    ```

5. Create `Makefile` in root folder and add the following:
    ```
    run:
	    templ generate && go build -o ./tmp/main cmd/compSudoku/main.go
    ```

6. Set up environment variables:
   Create a `.env` file in the root directory and define:
   ```
   SUDOKU_API_KEY=<your-api-key>
   ```

6. Run the application with hot reloading:
   ```bash
   air
   ```

   By default, the server runs on `http://localhost:8080`.

### Using Docker

Build and run the application using Docker:

```bash
docker build -t comp-sudoku .
docker run -p 8080:8080 -e SUDOKU_API_KEY=<your-api-key> comp-sudoku
```

## Usage

- Open `http://localhost:8080` in your browser.
- Start solving Sudoku puzzles and interact with the chat feature.
- Use the settings panel to generate new games or change difficulty levels.

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository.
2. Create a new branch (`git checkout -b feature-branch`).
3. Commit your changes (`git commit -m "Add new feature"`).
4. Push to the branch (`git push origin feature-branch`).
5. Open a Pull Request.

## License

This project is licensed under the [MIT License](LICENSE).

## Acknowledgements

- [Echo](https://echo.labstack.com/) - Web framework for Go
- [templ](https://github.com/a-h/templ) - Component-based template engine for Go
- [htmx](https://htmx.org/) - Client-side library for HTML-based interactivity
- [Golang Air](https://github.com/cosmtrek/air) - Hot reloading for Go applications
- [RapidAPI](https://rapidapi.com/) - Sudoku board generation API

---

Would you like additional instructions or details about integrating Air into your workflow?