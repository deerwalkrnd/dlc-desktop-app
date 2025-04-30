---
sidebar_position: 1
---

# DLC Desktop Application

This guide will take you though the process of setting up the Deerwalk Learning Center Desktop Application on your device.
We will begin things by assuming a few things about your system.

1. Your System is running on Windows 7 or higher.
2. Your System has 20MB Disc Space.
3. Your System has a web-browser installed.

## Getting Started

Get started by **Downloading the DLC Destop App** from **[release](https://docusaurus.new)**.

### Verifying

After downloading the zip file, ensure that the zip follows the given directory structure.

```
dlc-desktop
├── dlc-desktop-app.exe
└── web
```
Once this is done, we are ready to add videos/lectures to our application for further use.

## Add Video/Lecture

To add videos or lectures to our application, we need to create a new folder in the same folder we have the exe and the web files. In that directory we create a new folder with the name of `DLC`. Our app automatically globs over this folder and finds all the lecture files. For our lecture files to be available.

The videos must follow a very specific structure.


The video format must be: 

"TEACHER NAME - LESSION NO - LESSION TITLE - LECTURE TITLE - SUBJECT NAME - SUBJECT TYPE - CLASS .mp4"

Example:

"Arjun Khatiwada - 1.2 - Factorization of Polynomials - Polynomials in Maths - Mathematics - OLD - 8 .mp4"

An example folder structure can be something like this:

```
.
├── DLC
│   └── Arjun Khatiwada - 1.2 - Factorization of Polynomials - Polynomials in Maths - Mathematics - OLD - 8 .mp4
├── dlc-desktop-app.exe
└── web
```

The above example just shows one file but we can have as much file as possible in any folder as long as the name is correct, the system should pick it up.


## Running The Application
To run our application, we run the `dlc-desktop-app.exe` file in our folder which should open up our application and if the video structure is correct, the app should be running at `http://localhost:3000` 


## Adding new Videos
If you have added new videos to the `DLC` folder then you would need delete the sqlite database `dlc.sqlite` and re-run the application again.