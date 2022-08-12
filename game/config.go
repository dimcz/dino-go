package game

const (
	Height = 720.0
	Width  = 1280.0

	gameFPS = 60

	infoPadding = 10.0

	animationPause = 5
	runPosition    = 100.0
	textPadding    = 5.0
	dinoPadding    = 30.0
	jumpPower      = 10.0

	maxCountOfCacti = 3
	cactusPosition  = 110

	roadImage    = "assets/images/road.png"
	roadPosition = 100

	fontPath   = "assets/fonts/intuitive.ttf"
	normalSize = 26.0
	smallSize  = 16.0
	bigSize    = 100.0
)

var dinoNames = [...]string{
	"Fluffy", "Scruffy", "Ruffles", "Taz", "Ziggy", "Shaggy", "Frizzle", "Buster", "Wookiee",
	"Fizgig", "Bigfoot", "Beowulf", "Gordo", "Bear", "Woolly", "Oso", "Chewbacca", "Curly",
	"Gus", "Velvet",
}

var dinoColors = [...]string{
	"default", "aqua", "black", "bloody", "cobalt", "gold", "insta", "lime", "magenta",
	"magma", "navy", "neon", "orange", "pinky", "purple", "rgb", "silver", "subaru",
	"sunny", "toxic",
}
