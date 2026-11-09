---
sidebar_position: 4
title: Troubleshooting
---

# Troubleshooting

### A video I added doesn't show up

1. Open the **Library** page and press **Rescan library**.
2. Still missing? Check **Skipped files** on the same page. It says exactly
   what is wrong with the name. The most common mistakes are:
   - a missing part, often the class number at the end,
   - the parts separated by `-` with no spaces where a name also has a
     hyphen,
   - the lesson number written as text (`one`) instead of a number (`1`).
3. Check that the file is inside the folder shown at the top of the Library
   page.

### A video is in the wrong syllabus

The syllabus part must be exactly `OLD` (in capitals) for the old syllabus.
Anything else, including `old`, counts as the new syllabus. The Library page
shows a warning for these files.

### A subject or chapter appears twice

Subjects and chapters are matched by their exact spelling, including capital
letters and spaces. `Mathematics` and `Maths` are two different subjects, and
so are two chapters with the same number but different titles. Rename the
files so they match, then rescan.

### "This video could not be played"

- The file was moved, renamed or deleted since the last scan: press
  **Rescan library**.
- The file is damaged or isn't an H.264 MP4: try opening it in a media player.
  If it plays there, convert it to H.264 MP4 with any video converter.

### The browser shows "This site can't be reached"

The app isn't running. Start `dlc-desktop-app.exe` and keep its console
window open.

### "port 127.0.0.1:3000 is used by another program"

Another program is using port 3000. Close it, or start the app on a different
port from a Command Prompt:

```
dlc-desktop-app.exe -addr 127.0.0.1:3001
```

Then open `http://localhost:3001`.

### Using a video folder somewhere else

By default the app uses the `DLC` folder next to `dlc-desktop-app.exe`. To use
another folder, for example on a USB drive:

```
dlc-desktop-app.exe -library "E:\DLC"
```

### All command-line options

```
dlc-desktop-app.exe -help
```

| Option              | Default                     | Meaning                              |
| ------------------- | --------------------------- | ------------------------------------ |
| `-library <folder>` | `DLC` next to the app       | Where the videos are                 |
| `-addr <host:port>` | `127.0.0.1:3000`            | Address to listen on                 |
| `-open=false`       | open                        | Don't open the browser automatically |
| `-verbose`          | off                         | More detailed log output             |
| `-version`          |                             | Print the version and exit           |
