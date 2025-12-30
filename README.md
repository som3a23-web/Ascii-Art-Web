# ASCII Art Web

## Description

**ASCII Art Web** is a web-based application written in Go that transforms standard text into stylized ASCII art. The application provides a user-friendly interface where users can input text, choose from different banner styles (Standard, Shadow, and Thinkertoy), and instantly see the graphical representation of their text.

The project handles various edge cases, including multi-line inputs, special characters, and comprehensive HTTP error handling (400, 404, 405, and 500 status codes) to ensure a robust user experience.

## Authors

* **[Ismail](https://github.com/som3a23-web)**
* **[Anas](https://github.com/Anasmoner2022)**
* **[Omnia](https://github.com/OmniaAbdoun)**
* **[Amr](https://github.com/YOUR_ACTUAL_GITHUB)**

## Usage: How to Run

### Prerequisites

* [Go](https://go.dev/doc/install) (version 1.16 or higher recommended) installed on your machine.

### Installation & Execution

1. **Clone the repository:**
```bash
git clone https://github.com/som3a23-web/Ascii-Art-Web.git
cd Ascii-Art-Web
```

2. **Run the server:**
```bash
go run main.go
```

3. **Access the application:**
Open your web browser and navigate to:
```
http://localhost:8080
```

### How to Use the UI

1. Enter the text you wish to convert in the text area.
2. Select a banner style (Standard, Shadow, or Thinkertoy) using the radio buttons.
3. Click the **GENERATE ASCII ART** button.
4. The generated ASCII art will appear in the output section below.

---

## Implementation Details: Algorithm

The application follows a modular architecture, separating the web server logic from the ASCII generation engine.

### 1. Web Routing and Handling

The server uses the `net/http` package to handle requests:

* **GET `/`**: Serves the main landing page with the HTML form.
* **POST `/ascii-art`**: Receives the form data (text and banner choice), processes the ASCII art generation, and returns the rendered template with the result.
* **Static Files**: Serves CSS files from the `/static` directory to style the frontend.
* **Error Pages**: Custom error templates for 400, 404, 405, and 500 HTTP status codes.

### 2. Input Validation Logic

Before processing, the server validates the input:

* **Banner Validation**: Ensures the selected banner is one of the three allowed types (`standard`, `shadow`, `thinkertoy`).
* **Character Validation**: Ensures all input characters are within the printable ASCII range (32–126) or are valid newlines (`\n`, `\r`). If invalid characters are found, a `400 Bad Request` is returned.

### 3. The ASCII Generation Algorithm

The core logic resides in the `features` package. The algorithm follows these steps:

1. **Banner Loading**: The program reads the selected `.txt` banner file (located in the `banner/` directory). Each character in these files is represented by **8 lines** of ASCII art text.

2. **Input Parsing**: The user's input string is split by newline characters (`\n`) to handle multi-line requests.

3. **Character Mapping**:
   * The algorithm iterates through each character of the input string.
   * It calculates the starting line in the banner file using the formula: 
```
     StartingLine = (ASCII_Value - 32) * 9 + 1
```
   * It extracts the 8 corresponding lines for that character from the banner file.

4. **Line-by-Line Construction**: 
   * To print a single line of input text, the algorithm must print the **first line** of every character in that sequence, then the **second line**, and so on (8 iterations per line of input).
   * Empty lines in the input are handled by adding a newline to the output.

5. **Output Assembly**: The resulting ASCII art lines are joined back into a single string and passed to the HTML template for display within a `<pre>` tag to preserve formatting.

### 4. HTTP Error Handling

The project implements custom templates for various HTTP errors:

* **400 Bad Request**: Triggered by invalid input (unsupported characters or invalid banner selection).
* **404 Not Found**: Triggered by accessing undefined routes.
* **405 Method Not Allowed**: Triggered by using the wrong HTTP verb (e.g., GET on `/ascii-art` or POST on `/`).
* **500 Internal Server Error**: Triggered if banner files are missing or the server encounters an unexpected filesystem issue.

### 5. Template Rendering

Go's `html/template` package is used to:
* Serve the main page with dynamic content
* Display the generated ASCII art in the output section
* Show appropriate error pages with custom designs

---

## Project Structure
```
Ascii-Art-Web/
├── banner/
│   ├── standard.txt
│   ├── shadow.txt
│   └── thinkertoy.txt
├── features/
│   └── asciidrawing.go
├── static/
│   ├── style.css
│   └── [error SVG images]
├── templates/
│   ├── index.html
│   ├── 400.html
│   ├── 404.html
│   ├── 405.html
│   └── 500.html
├── main.go
├── main_test.go
├── go.mod
└── README.md
```

---

## Testing

Run the test suite:
```bash
go test -v
```

The tests cover:
- Template existence
- Banner functionality with special characters
- HTTP status code handling
- Input validation
- Error handling scenarios

---

## Technologies Used

* **Backend**: Go (Golang) - HTTP server, routing, and ASCII generation
* **Frontend**: HTML5, CSS3 (custom styling with CSS variables)
* **Templates**: Go's `html/template` package
* **Font**: VT323 (monospace, terminal-style)

---

## Features

✅ Three banner styles (Standard, Shadow, Thinkertoy)  
✅ Multi-line text support  
✅ Special character handling  
✅ Real-time ASCII art generation  
✅ Responsive design  
✅ Comprehensive error handling  
✅ Custom error pages with illustrations  

---

## License

This project is part of the Zone01 Kisumu curriculum.