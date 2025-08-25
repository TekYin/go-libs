package colors

var red = "\033[31m"
var yellow = "\033[33m"
var green = "\033[32m"
var blue = "\033[34m"
var cyan = "\033[36m"
var reset = "\033[0m"

func RedString(s string) string {
	return red + s + reset
}

func YellowString(s string) string {
	return yellow + s + reset
}

func GreenString(s string) string {
	return green + s + reset
}

func BlueString(s string) string {
	return blue + s + reset
}

func CyanString(s string) string {
	return cyan + s + reset
}
