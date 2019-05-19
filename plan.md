# Modular backend

Files:

main.go - Get rid of command line scripts, just runs server

    - main() Will start the server

server.go - Will contain the server routes and will start the server

- StartServer() Will start the server

bot.go - Contains the main application functions that will be exposed to server.js

- Download()
- Compile()
- FindDuration()
- Upload()


download.go 
- GetTwitchData - will get the data from twitch
- DownloadClips - Will download all the clips

compile.go
- createThumbnail() - Creates the thumbnail
- convertToTS 
- concatFiles()


frontend:

<App /> - Container for the app, will run all components

    <LeftPanel /> - Contains the elements in the left part of the admin panel
        
        <Download /> - Contain all code to Download clips (including amount to download)
        <Concat /> - Contain all code to Concat clips (including custom title name)
        <VidDuration / > - Contain all code for finding video duration
        <Upload /> - Contains all code to upload to yt

    <ClipPlayer /> - Container component that shows the clips (including moving and deleting
    through vids)

    <VidPlayer> - Container component that shows the final vid
        
