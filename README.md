# go-drive

tui for Google Drive in `Go`

## usage

1. you will first need to configure a Google Cloud project. follow instructions [from here](https://developers.google.com/workspace/drive/api/quickstart/go) to set it up
    - make a Cloud project
    - enable Google Drive API for the project
    - create a OAuth client picking "Desktop app" as the application type
    - you will need the `credentials.json` file (the actual name is different, so you'll need to rename it)
2. put the `credentials.json` file to `$HOME/.config/godrive/` (alternatively, you may pass a different path in the command line arguments to the app)
3. authenticate once using `go-drive auth`; this produces a `token.json` file in the same `$HOME/.config/godrive/` directory which is then used each time you make API calls
4. the filesystem can be accessed interactively via `go-drive fs`

*installation instructions to be posted*

for the best looks, you'll need to use [one of the nerdfonts](https://www.nerdfonts.com).

### To-Do

- [ ] selective sync functionality
- [ ] seamless merge of multiple accounts (combining multiple drive filesystems into one via tui)

### disclaimer

even though this project is aimed to make working with Google Drive easier, in no way do i endorse using Google Drive in general. i can think of myriads of reasons on why not to use Google Drive, some of which are outlined [in this article](https://proton.me/blog/is-google-drive-secure). there are plenty of alternatives ranging from fully open-source and even self-hosted to proprietary. the reasons i targeted Google Drive is:
- they provide a free API support (something, say, Proton Drive lacks);
- their free tier has 15 GB of space, which i don't believe any of other services match.

in other words: unless you absolutely have to, do NOT use Google Drive. and for gods sake, definitely don't buy their paid tier.
