# GoBooking.com

GoBooking.com is a booking and reservation project built with Go 1.23.4. This project leverages several Go libraries to provide a robust and secure booking system.

## Features

- **CSRF Protection**: Implements CSRF protection using nosurf.
- **Chi Router**: A lightweight, idiomatic and composable router for building Go HTTP services.
- **Session Management**: Utilizes Alex Edwards' SCS session management for secure and efficient session handling.

## Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/psanodiya94/gobooking.com.git
    ```
2. Navigate to the project directory:
    ```sh
    cd gobooking.com
    ```
3. Install dependencies:
    ```sh
    go mod tidy
    ```

## Usage

1. Run the application:
    ```sh
    go run main.go
    ```
2. Open your browser and navigate to `http://localhost:8080` to access the application.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any changes.

## License


## Acknowledgements

- [nosurf](https://github.com/justinas/nosurf)
- [Chi Router](https://github.com/go-chi/chi)
- [SCS Session Management](https://github.com/alexedwards/scs)