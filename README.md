# statii - a status tool

Statii is a terminal application for displaying notifications from various services. The project arose because I noticed that users (myself included) increasingly disable notifications, due to notification flooding. At the same time, there are specific notifications a user wants, based on their specific needs. 

In some cases, notifcations aren't made available for the systems in question - and statii becomes your way to write custom plugins to surface them.

Statii will draw a self-updating table, in the terminal. As of this writing, it updates every 30 seconds, but will properly support async updates in future versions. Clicking on any table row will open the associated notification link, in your browser.

## Building

Statii requires [go version 1.19](https://go.dev/) or higher and [gnu make](https://www.gnu.org/software/make/), though the latter is a wrapper around ```go build```. But, once installed you need to run ```make```.

## Usage

Configure your statii.conf - as per the [example file](statii.conf.example), and run the application via ```./statii```

*./statii --help* provides useful output.
