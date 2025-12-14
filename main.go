package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
)

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorGray   = "\033[90m"
	ColorBold   = "\033[1m"
)

var (
	// Input/Output flags
	inputFile  = flag.String("i", "", "Input JSON file")
	outputFile = flag.String("o", "", "Output file")

	// Formatting flags
	indent   = flag.String("indent", "  ", "Indentation string (spaces/tabs)")
	prefix   = flag.String("prefix", "", "Prefix for each line")
	compact  = flag.Bool("c", false, "Compact JSON (minify)")
	sortKeys = flag.Bool("s", false, "Sort keys alphabetically")
	colorize = flag.Bool("color", true, "Colorize output (if writing to terminal)")

	// Utility flags
	validate = flag.Bool("v", false, "Validate JSON only (no output)")
	help     = flag.Bool("h", false, "Show help menu")
	version  = flag.Bool("version", false, "Show version")
)

const VERSION = "1.1.0"

func main() {
	flag.Usage = printHelp
	flag.Parse()

	if *version {
		printVersion()
		return
	}
	if *help {
		printHelp()
		return
	}

	// 1. Determine Input Source
	var reader io.Reader
	if *inputFile != "" {
		f, err := os.Open(*inputFile)
		if err != nil {
			printError(fmt.Sprintf("Error opening file: %v", err))
			os.Exit(1)
		}
		defer f.Close()
		reader = f
		printInfo(fmt.Sprintf("Processing: %s", *inputFile))
	} else {
		// Check stdin
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			printError("No input provided. Use -i or pipe data.")
			os.Exit(1)
		}
		reader = os.Stdin
	}

	// 2. Prepare Output
	var writer io.Writer = os.Stdout
	if *outputFile != "" {
		f, err := os.Create(*outputFile)
		if err != nil {
			printError(fmt.Sprintf("Error creating output file: %v", err))
			os.Exit(1)
		}
		defer f.Close()
		writer = f

		*colorize = false
	} else {
		// If piping to another command (not a TTY), usually disable color
		// We simulate isatty check here safely without external deps
		fi, _ := os.Stdout.Stat()
		if (fi.Mode() & os.ModeCharDevice) == 0 {
			*colorize = false
		}
	}

	// 3. Process the Stream
	decoder := json.NewDecoder(reader)
	encoder := json.NewEncoder(writer)

	// Configure Encoder
	encoder.SetEscapeHTML(false) // Don't escape < > &
	if !*compact {
		encoder.SetIndent(*prefix, *indent)
	}

	// We use a loop to handle multiple JSON objects in one stream (ndjson support)
	count := 0
	for {
		var raw json.RawMessage
		var data interface{}

		// Decode logic
		if *sortKeys {
			// To sort, we must decode into interface{}
			if err := decoder.Decode(&data); err != nil {
				if err == io.EOF {
					break
				}
				printError(fmt.Sprintf("Invalid JSON: %v", err))
				os.Exit(1)
			}
		} else {
			// To preserve order (default), decode into RawMessage
			if err := decoder.Decode(&raw); err != nil {
				if err == io.EOF {
					break
				}
				printError(fmt.Sprintf("Invalid JSON: %v", err))
				os.Exit(1)
			}
		}

		if *validate {
			count++
			continue
		}

		// Encode logic (Output)
		var err error
		var outputBytes []byte

		if *sortKeys {
			if *compact {
				outputBytes, err = json.Marshal(data)
			} else {
				outputBytes, err = json.MarshalIndent(data, *prefix, *indent)
			}
		} else {
			// If we have RawMessage
			if *compact {
				buffer := new(bytes.Buffer)
				if err := json.Compact(buffer, raw); err != nil {
					printError(fmt.Sprintf("Compact error: %v", err))
					os.Exit(1)
				}
				outputBytes = buffer.Bytes()
			} else {
				// Re-indenting RawMessage requires unmarshal/marshal cycle usually,
				// or use the encoder directly. But for highlighting, we need bytes.
				// Let's use the standard indentation on the RawMessage
				var buf bytes.Buffer
				err = json.Indent(&buf, raw, *prefix, *indent)
				outputBytes = buf.Bytes()
			}
		}

		if err != nil {
			printError(fmt.Sprintf("Encoding error: %v", err))
			os.Exit(1)
		}

		// Apply Syntax Highlighting if enabled
		if *colorize && *outputFile == "" {
			outputBytes = syntaxHighlight(outputBytes)
		}

		// Write to output
		writer.Write(outputBytes)
		writer.Write([]byte("\n"))
		count++
	}

	if *validate {
		printSuccess(fmt.Sprintf("✓ Validated %d JSON object(s)", count))
	} else if *outputFile != "" {
		printSuccess(fmt.Sprintf("Saved to: %s", *outputFile))
	}
}

// syntaxHighlight adds ANSI color codes to JSON
func syntaxHighlight(js []byte) []byte {
	str := string(js)

	keyColor := ColorBlue + ColorBold
	stringColor := ColorGreen
	numberColor := ColorYellow
	boolColor := ColorPurple
	nullColor := ColorRed
	reset := ColorReset

	// Highlight Keys (captured by "key": )
	reKey := regexp.MustCompile(`"([^"]+)"\s*:`)
	str = reKey.ReplaceAllString(str, keyColor+`"$1"`+reset+`:`)

	// Highlight Strings (values that are strings, look for "text" not followed by colon)
	// This regex is tricky to not overlap with keys.
	// We cheat slightly by doing keys first, which adds ANSI codes,
	// so the next regex won't match keys because they now contain \033.
	reString := regexp.MustCompile(`:(\s*)"([^"]*)"`)
	str = reString.ReplaceAllString(str, `:`+`$1`+stringColor+`"$2"`+reset)

	// Highlight Numbers
	reNum := regexp.MustCompile(`:(\s*)([0-9]+(?:\.[0-9]+)?(?:[eE][+-]?[0-9]+)?)`)
	str = reNum.ReplaceAllString(str, `:`+`$1`+numberColor+`$2`+reset)

	// Highlight Booleans
	reBool := regexp.MustCompile(`:(\s*)(true|false)`)
	str = reBool.ReplaceAllString(str, `:`+`$1`+boolColor+`$2`+reset)

	// Highlight Null
	reNull := regexp.MustCompile(`:(\s*)(null)`)
	str = reNull.ReplaceAllString(str, `:`+`$1`+nullColor+`$2`+reset)

	return []byte(str)
}

// ---------------- Helper Functions ----------------

func printHelp() {
	banner := `
╔═══════════════════════════════════════════════════════════╗
║                                                           ║
║            JSON BEAUTIFIER PRO v` + VERSION + `                     ║
║                                                           ║
╚═══════════════════════════════════════════════════════════╝
`
	fmt.Print(ColorCyan + ColorBold + banner + ColorReset)
	fmt.Println(ColorYellow + ColorBold + "\n📖 USAGE:" + ColorReset)
	fmt.Println("  " + ColorGreen + "go run main.go [OPTIONS]" + ColorReset)

	fmt.Println(ColorYellow + ColorBold + "\n⚙️  OPTIONS:" + ColorReset)

	printFlag("-i <file>", "Input JSON file (or stdin)")
	printFlag("-o <file>", "Output file (stdout if empty)")
	printFlag("-c", "Compact/minify JSON")
	printFlag("-s", "Sort keys alphabetically")
	printFlag("-color", "Force color output (default: auto)")
	printFlag("-indent", "Custom indentation (default: 2 spaces)")
	printFlag("-v", "Validate only")
	printFlag("-h", "Show help")

	fmt.Println(ColorYellow + ColorBold + "\n💡 FEATURES:" + ColorReset)
	fmt.Println("  • Syntax Highlighting")
	fmt.Println("  • Stream Processing (NDJSON support)")
	fmt.Println("  • Key Sorting")
}

func printFlag(flag, desc string) {
	fmt.Printf("  "+ColorCyan+"%-15s"+ColorReset+" %s\n", flag, desc)
}

func printVersion() {
	fmt.Printf(ColorCyan+ColorBold+"JSON Beautifier Pro v%s\n"+ColorReset, VERSION)
}

func printError(msg string) {
	fmt.Fprintf(os.Stderr, ColorRed+ColorBold+"[ERROR] "+ColorReset+ColorRed+"%s\n"+ColorReset, msg)
}

func printSuccess(msg string) {
	fmt.Fprintf(os.Stderr, ColorGreen+ColorBold+"[✓] "+ColorReset+ColorGreen+"%s\n"+ColorReset, msg)
}

func printInfo(msg string) {
	fmt.Fprintf(os.Stderr, ColorBlue+"[→] "+ColorReset+"%s\n", msg)
}
