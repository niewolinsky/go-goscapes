# GoScapes - Ambient sounds from terminal
TUI application allowing you to play relaxing, ambient sounds. Created while learning [Bubbletea](https://github.com/charmbracelet/bubbletea) framework. I am now certified front-end terminal developer 😎. The code needs a lot of improvement but the application works.

![GoScapes Screenshot](https://i.imgur.com/N0Q0XvO.png)

[GIF Version](https://i.imgur.com/WGMEbmr.mp4)

## Features:
- Play up to 9 different soundscapes using user interface in the terminal.
- Control each soundscape volume, stop and start sounds on demand.
- All audio files embedded using Go's FS filesystem package (single binary).

List of available ambient sounds:
- rain
- thunder
- waves
- wind
- fire
- birds
- crickets
- singing bowls
- white noise

## Running:
Simply execute binary OR build yourself, you can create "soundscapes/" folder and upload up to 9 ".mp3" files if you want to customize the sounds.

Tested on Linux.

## Todo:
- 2D navigation, instead of tab navigation
- visual indicator for active (audible) soundscape
- visual indicator for current soundscape volume
- easier custom soundscapes uploading

## Stack:
- Go 1.20 + `charmbracelet/bubbletea` + `hajimehoshi/oto`

## Credits:
Thanks to [Hajime Hoshi](https://github.com/hajimehoshi) for creating simple and accessible Go audio libraries and [Charm](https://github.com/charmbracelet) for excellent command-line tools and libraries.
