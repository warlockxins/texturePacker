package space

type FrameSize struct {
	W int `json:"w"`
	H int `json:"h"`
}
type Meta struct {
	App     string `json:"app"`     // https://www.codeandweb.com/texturepacker
	Version string `json:"version"` //example had "3"
}

type Frame struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

/*
@example
filename: "blueJellyfish0001"
frame: {x: 484, y: 836, w: 63, h: 65}
rotated: false
sourceSize: {w: 66, h: 66}
spriteSourceSize: {x: 1, y: 0, w: 66, h: 66}
trimmed: true
*/
type ImageFrame struct {
	FileName         string    `json:"fileName"` // name of a file frame
	Frame            Frame     `json:"frame"`
	Rotated          bool      `json:"rotated"`
	SourceSize       FrameSize `json:"sourceSize"`
	SpriteSourceSize Frame     `json:"spriteSourceSize"`
	Trimmed          bool      `json:"trimmed"`
}

type Texture struct {
	Format string       `json:"format"` //  RGBA8888
	Image  string       `json:"image"`  // name of png file associated
	Scale  string       `json:"scale"`  // "1"
	Size   FrameSize    `json:"size"`
	Frames []ImageFrame `json:"frames"`
}

type SpriteAtlas struct {
	Textures []Texture `json:"textures"`
	Meta     Meta      `json:"meta"`
}

// example https://labs.phaser.io/view.html?src=src\animation\create%20animation%20from%20texture%20atlas.js
// newer example https://labs.phaser.io/view.html?src=src\animation\chained%20animation.js
