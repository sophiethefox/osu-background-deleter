# osu! Background Deleter

This app allows you to quickly delete the background of the current map! All you need to do is start the app, and whenever you select a map, the background will be shown in the app along with a `Delete` button.

By default, this app will minimize to the system tray when closed. This is so that you don't need to have an extra application open when not playing osu!, while still allowing easy access. This can be changed in `config.ini`, and is the only setting you may need to change.

Deleted backgrounds are stores in a `deleted_backgrounds` directory next to the executable. This is done so that the previously deleted background can be restored, and so that you can manually restore select backgrounds. You may delete this folder at any time.

# Download

Please go to the [Releases](https://github.com/sophiethefox/osu-background-deleter/releases) page to download the latest release. This app is only available for 64-bit Windows currently.

# Planned features

- Delete beatmap thumbnail (cache stored in `osu/Data/bt`)
- Ability to restore background of current beatmap (if saved)

# Credits

Much of this project was based upon [l3lackShark](https://github.com/l3lackShark)'s and others [gosumemory](https://github.com/l3lackShark/gosumemory). Specifically, I am using a stripped down version of his memory reader implementation to see what map the user has selected in osu!. This project is licensed under GNU GPL-3 to match gosumemory.
