1. Copy the `list-deps.sh` (Linux/OSX) or `list-deps.bat` (Windows) file from the `android` dir to the root of your Android project.
2. Change it to match your configuration; your app module may not be named `app`
3. Execute it, command line preferably
4. The result will be in `gradle-dep-list.txt` file.
5. Place `gradle-dep-list.txt` and `gradle-dep-ack(.exe)` files together
6. Execute `gradle-dep-ack(.exe)`