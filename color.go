package ilibgo

import (
	"fmt"
	"image/color"
)

type rgb struct {
	red   int
	green int
	blue  int
}

var colorMap = map[string]rgb{
	// X11 rgb.txt
	"aliceblue":            {240, 248, 255},
	"antiquewhite":         {250, 235, 215},
	"aquamarine":           {127, 255, 212},
	"azure":                {240, 255, 255},
	"beige":                {245, 245, 220},
	"bisque":               {255, 228, 196},
	"black":                {0, 0, 0},
	"blanchedalmond":       {255, 235, 205},
	"blue":                 {0, 0, 255},
	"blueviolet":           {138, 43, 226},
	"brown":                {165, 42, 42},
	"burlywood":            {222, 184, 135},
	"cadetblue":            {95, 158, 160},
	"chartreuse":           {127, 255, 0},
	"chocolate":            {210, 105, 30},
	"coral":                {255, 127, 80},
	"cornflowerblue":       {100, 149, 237},
	"cornsilk":             {255, 248, 220},
	"cyan":                 {0, 255, 255},
	"darkblue":             {0, 0, 139},
	"darkcyan":             {0, 139, 139},
	"darkgoldenrod":        {184, 134, 11},
	"darkgray":             {169, 169, 169},
	"darkgreen":            {0, 100, 0},
	"darkgrey":             {169, 169, 169},
	"darkkhaki":            {189, 183, 107},
	"darkmagenta":          {139, 0, 139},
	"darkolivegreen":       {85, 107, 47},
	"darkorange":           {255, 140, 0},
	"darkorchid":           {153, 50, 204},
	"darkred":              {139, 0, 0},
	"darksalmon":           {233, 150, 122},
	"darkseagreen":         {143, 188, 143},
	"darkslateblue":        {72, 61, 139},
	"darkslategray":        {47, 79, 79},
	"darkslategrey":        {47, 79, 79},
	"darkturquoise":        {0, 206, 209},
	"darkviolet":           {148, 0, 211},
	"deeppink":             {255, 20, 147},
	"deepskyblue":          {0, 191, 255},
	"dimgray":              {105, 105, 105},
	"dimgrey":              {105, 105, 105},
	"dodgerblue":           {30, 144, 255},
	"firebrick":            {178, 34, 34},
	"floralwhite":          {255, 250, 240},
	"forestgreen":          {34, 139, 34},
	"gainsboro":            {220, 220, 220},
	"ghostwhite":           {248, 248, 255},
	"gold":                 {255, 215, 0},
	"goldenrod":            {218, 165, 32},
	"gray":                 {190, 190, 190},
	"green":                {0, 255, 0},
	"greenyellow":          {173, 255, 47},
	"grey":                 {190, 190, 190},
	"honeydew":             {240, 255, 240},
	"hotpink":              {255, 105, 180},
	"indianred":            {205, 92, 92},
	"ivory":                {255, 255, 240},
	"khaki":                {240, 230, 140},
	"lavender":             {230, 230, 250},
	"lavenderblush":        {255, 240, 245},
	"lawngreen":            {124, 252, 0},
	"lemonchiffon":         {255, 250, 205},
	"lightblue":            {173, 216, 230},
	"lightcoral":           {240, 128, 128},
	"lightcyan":            {224, 255, 255},
	"lightgoldenrod":       {238, 221, 130},
	"lightgoldenrodyellow": {250, 250, 210},
	"lightgray":            {211, 211, 211},
	"lightgreen":           {144, 238, 144},
	"lightgrey":            {211, 211, 211},
	"lightpink":            {255, 182, 193},
	"lightsalmon":          {255, 160, 122},
	"lightseagreen":        {32, 178, 170},
	"lightskyblue":         {135, 206, 250},
	"lightslateblue":       {132, 112, 255},
	"lightslategray":       {119, 136, 153},
	"lightslategrey":       {119, 136, 153},
	"lightsteelblue":       {176, 196, 222},
	"lightyellow":          {255, 255, 224},
	"limegreen":            {50, 205, 50},
	"linen":                {250, 240, 230},
	"magenta":              {255, 0, 255},
	"maroon":               {176, 48, 96},
	"mediumaquamarine":     {102, 205, 170},
	"mediumblue":           {0, 0, 205},
	"mediumorchid":         {186, 85, 211},
	"mediumpurple":         {147, 112, 219},
	"mediumseagreen":       {60, 179, 113},
	"mediumslateblue":      {123, 104, 238},
	"mediumspringgreen":    {0, 250, 154},
	"mediumturquoise":      {72, 209, 204},
	"mediumvioletred":      {199, 21, 133},
	"midnightblue":         {25, 25, 112},
	"mintcream":            {245, 255, 250},
	"mistyrose":            {255, 228, 225},
	"moccasin":             {255, 228, 181},
	"navajowhite":          {255, 222, 173},
	"navy":                 {0, 0, 128},
	"navyblue":             {0, 0, 128},
	"oldlace":              {253, 245, 230},
	"olivedrab":            {107, 142, 35},
	"orange":               {255, 165, 0},
	"orangered":            {255, 69, 0},
	"orchid":               {218, 112, 214},
	"palegoldenrod":        {238, 232, 170},
	"palegreen":            {152, 251, 152},
	"paleturquoise":        {175, 238, 238},
	"palevioletred":        {219, 112, 147},
	"papayawhip":           {255, 239, 213},
	"peachpuff":            {255, 218, 185},
	"peru":                 {205, 133, 63},
	"pink":                 {255, 192, 203},
	"plum":                 {221, 160, 221},
	"powderblue":           {176, 224, 230},
	"purple":               {160, 32, 240},
	"red":                  {255, 0, 0},
	"rosybrown":            {188, 143, 143},
	"royalblue":            {65, 105, 225},
	"saddlebrown":          {139, 69, 19},
	"salmon":               {250, 128, 114},
	"sandybrown":           {244, 164, 96},
	"seagreen":             {46, 139, 87},
	"seashell":             {255, 245, 238},
	"sienna":               {160, 82, 45},
	"skyblue":              {135, 206, 235},
	"slateblue":            {106, 90, 205},
	"slategray":            {112, 128, 144},
	"slategrey":            {112, 128, 144},
	"snow":                 {255, 250, 250},
	"springgreen":          {0, 255, 127},
	"steelblue":            {70, 130, 180},
	"tan":                  {210, 180, 140},
	"thistle":              {216, 191, 216},
	"tomato":               {255, 99, 71},
	"turquoise":            {64, 224, 208},
	"violet":               {238, 130, 238},
	"violetred":            {208, 32, 144},
	"wheat":                {245, 222, 179},
	"white":                {255, 255, 255},
	"whitesmoke":           {245, 245, 245},
	"yellow":               {255, 255, 0},
	"yellowgreen":          {154, 205, 50},
	// https://cloford.com/resources/colours/namedcol.htm
	"crimson": {220, 20, 60},
	"fuchsia": {255, 0, 255},
	"indigo":  {75, 0, 130},
	"aqua":    {0, 255, 255},
	"teal":    {0, 128, 128},
	"lime":    {0, 255, 0},
	"olive":   {128, 128, 0},
	"silver":  {192, 192, 192},
	// https://en-academic.com/dic.nsf/enwiki/1589972
	"almond":                 {239, 222, 205},
	"antiquebrass":           {205, 149, 117},
	"apricot":                {253, 217, 181},
	"asparagus":              {135, 169, 107},
	"atomictangerine":        {255, 164, 116},
	"bananamania":            {250, 231, 181},
	"beaver":                 {159, 129, 112},
	"bittersweet":            {253, 124, 110},
	"blizzardblue":           {172, 229, 238},
	"bluebell":               {162, 162, 208},
	"bluegray":               {102, 153, 204},
	"bluegreen":              {13, 152, 186},
	"blush":                  {222, 93, 131},
	"brickred":               {203, 65, 84},
	"burntorange":            {255, 127, 73},
	"burntsienna":            {234, 126, 93},
	"canary":                 {255, 255, 153},
	"caribbeangreen":         {28, 211, 162},
	"carnationpink":          {255, 170, 204},
	"cerise":                 {221, 68, 146},
	"cerulean":               {29, 172, 214},
	"chestnut":               {188, 93, 88},
	"copper":                 {221, 148, 117},
	"cornflower":             {154, 206, 235},
	"cottoncandy":            {255, 188, 217},
	"dandelion":              {253, 219, 109},
	"denim":                  {43, 108, 196},
	"desertsand":             {239, 205, 184},
	"eggplant":               {110, 81, 96},
	"electriclime":           {206, 255, 29},
	"fern":                   {113, 188, 120},
	"fuzzywuzzy":             {204, 102, 102},
	"grannysmithapple":       {168, 228, 160},
	"greenblue":              {17, 100, 180},
	"hotmagenta":             {255, 29, 206},
	"inchworm":               {178, 236, 93},
	"jazzberryjam":           {202, 55, 103},
	"junglegreen":            {59, 176, 143},
	"laserlemon":             {254, 254, 34},
	"lemonyellow":            {255, 244, 79},
	"macaroniandcheese":      {255, 189, 136},
	"magicmint":              {170, 240, 209},
	"mahogany":               {205, 74, 76},
	"maize":                  {237, 209, 156},
	"manatee":                {151, 154, 170},
	"mangotango":             {255, 130, 67},
	"mauvelous":              {239, 152, 170},
	"melon":                  {253, 188, 180},
	"mountainmeadow":         {48, 186, 143},
	"mulberry":               {197, 75, 140},
	"neoncarrot":             {255, 163, 67},
	"olivegreen":             {186, 184, 108},
	"orangeyellow":           {248, 213, 104},
	"outerspace":             {65, 74, 76},
	"outrageousorange":       {255, 110, 74},
	"pacificblue":            {28, 169, 201},
	"peach":                  {255, 207, 171},
	"periwinkle":             {197, 208, 230},
	"piggypink":              {253, 221, 230},
	"pinegreen":              {21, 128, 120},
	"pinkflamingo":           {252, 116, 253},
	"pinksherbert":           {247, 143, 167},
	"purpleheart":            {116, 66, 200},
	"purplemountainsmajesty": {157, 129, 186},
	"purplepizzazz":          {254, 78, 218},
	"radicalred":             {255, 73, 108},
	"rawsienna":              {214, 138, 89},
	"rawumber":               {113, 75, 35},
	"razzledazzlerose":       {255, 72, 208},
	"razzmatazz":             {227, 37, 107},
	"redorange":              {255, 83, 73},
	"redviolet":              {192, 68, 143},
	"robinseggblue":          {31, 206, 203},
	"royalpurple":            {120, 81, 169},
	"scarlet":                {252, 40, 71},
	"screamingreen":          {118, 255, 122},
	"sepia":                  {165, 105, 79},
	"shadow":                 {138, 121, 93},
	"shamrock":               {69, 206, 162},
	"shockingpink":           {251, 126, 253},
	"sunglow":                {255, 207, 72},
	"sunsetorange":           {253, 94, 83},
	"tealblue":               {24, 167, 181},
	"ticklemepink":           {252, 137, 172},
	"timberwolf":             {219, 215, 210},
	"tropicalrainforest":     {23, 128, 109},
	"tumbleweed":             {222, 170, 136},
	"turquoiseblue":          {119, 221, 231},
	"unmellowyellow":         {255, 255, 102},
	"violetpurple":           {146, 110, 174},
	"violetblue":             {50, 74, 178},
	"vividtangerine":         {255, 160, 137},
	"vividviolet":            {143, 80, 157},
	"wildblueyonder":         {162, 173, 208},
	"wildstrawberry":         {255, 67, 164},
	"wildwatermelon":         {252, 108, 133},
	"wisteria":               {205, 164, 222},
	"yelloworange":           {255, 174, 66},
}

func newIColor(red uint8, green uint8, blue uint8, alpha uint8) Color {
	return Color{color: color.RGBA{red, green, blue, alpha}}
}

// Allocate a color (with no alpha) by its name (e.g. "red", "blue", etc)
func AllocNamedColor(name string) (Color, error) {
	if val, ok := colorMap[name]; ok {
		c := newIColor(uint8(val.red), uint8(val.green), uint8(val.blue), 255)
		return c, nil
	} else {
		c := newIColor(255, 255, 255, 255)
		return c, fmt.Errorf("no such color'%s'", name)
	}
}

// Allocate a color (with no alpha) from RGB values 0-255.
func AllocColor(red, green, blue uint8) (Color, error) {
	c := newIColor(red, green, blue, 255)
	return c, nil
}
