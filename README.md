# 🎛 Kalterwind MIDI Mapper

**Kalterwind Mapper** is a cross-platform MIDI routing and transformation tool written in Go. It allows you to map and process MIDI messages using flexible YAML-based presets.

---

### What's the need?
At the very first it was [written using C#](https://github.com/morganchester/midi_kalterwindmapper) to map usual MIDI-keyboard to the **Korg Kaossilator** and **Korg Kaoss Pad 3** on the stage

When switched to MacOS, I decided to make it cross-platform

You can find a presets for Kaossilator in the [Presets](presets) directory

## ✅ Features:

* **Custom MIDI Mapping via YAML**
  Define how incoming messages are transformed before reaching the output device:

  ```yaml
  inputDevice: "GenericMidi"
  outputDevice: "Korg Kaossilator"
  octaves: 2
  description: "Chromatic scale, 2 octaves"

  mapping:
    noteon:
      60: [ cc, 12, 5 ]
      61: [ cc, 12, 9 ]
    noteoff:
      action: "send"
      message: [ cc, 92, 0 ]

  controls:
    13: [ cc, 13 ]
    90: [ cc, 90 ]
  ```

* **Works with Any Device**
  Route from any MIDI input to any MIDI output (hardware or software).

* **Virtual MIDI Ports**
  On macOS and Linux, virtual ports are created automatically for DAW integration (DAW → Mapper → Device or Device → Mapper → DAW).
  On Windows, use **loopMIDI** or similar tools for virtual routing (native support planned).

* **Cross-Platform GUI**
  Built with [Fyne](https://fyne.io), providing a simple interface for:

    * Selecting input/output devices
    * Choosing presets
    * Starting/stopping routing

---

### Why?

It solves scenarios like:

* Translating standard MIDI notes into Kaossilator X/Y CC controls.
* Custom mappings for unique hardware setups.
* Acting as a live bridge between DAWs and hardware devices.

---

### What is *Kalterwind*?

It's the German name for the music band "Холодный ветер", where I'm a co-founder. So, it's just a name)