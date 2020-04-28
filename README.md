Go mockup of windowless/widgetless rendering using SDL2

# Demo

![Demo Image](https://github.com/ahmedaly113/roadbook/blob/master/demo.png)

# Prerequisites
## Go

Standard macOS install

## SDL2

macOS:
```brew install sdl2{,_image,_mixer,_ttf,_gfx} pkg-config```

# Instructions
```make run```

# Keymap
a=roadbook scroll forward
z=roadbook scroll reverse

space=menu

up=increment odometer (or menu up)
down=decrement odometer  (or menu down)
right=mode/reset (or menu accept)

# Recommendations
* use Visual Studio Code
* install Go extension + Go tools
* use CMD-Click to navigate to symbol definitions

# Safety and Reliability Tips

## Reserve panic for unexpected situations

Calling `panic(...)` should be reserved only for conditions that suggest an error in the code, or possibly in the hardware/runtime. Timeouts, invalid user inputs, empty files, EOFs: all of these should be detected and returned from functions as errors. https://stackoverflow.com/questions/44504354/should-i-use-panic-or-return-error

## Goroutines with the safe package

Any goroutine longer than a line or two should probably be called with `safe.Do` to avoid unmanaged crashes.

A function can recover from a panic with recover:

```
func a() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered in f", r)
        }
    }()
}
```

however, if a panic occurs in a "naked" goroutine, the panic cannot be caught by the function which originally launched the goroutine. It has "detached" from the parent. Calling `safe.Do`
consistently, even for simple functions, will make naked goroutines easy to see.

Because the panic'd goroutine will terminate, it may be necessary to gracefully shut down the
application anyway after a panic is detected. Call `safe.DoOrQuit` with a `chan struct{}` to
propagate the quit signal throughout the application. See https://medium.com/@matryer/golang-advent-calendar-day-two-starting-and-stopping-things-with-a-signal-channel-f5048161018 for more on quit/signal channels.

## errors.Wrap your errors

Some errors will end up in log file at some point. Use `errors.Wrap` to provide context as you return the error:

```
func resizePNG(src []byte, imgHeight uint) ([]byte, error) {
	i, _, err := image.Decode(bytes.NewReader(src))
	if err != nil {
		return nil, errors.Wrap(err, "problem decoding PNG during resize")
	}
...
```