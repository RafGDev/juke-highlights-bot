import React, { Component } from "react";
import "./LeftPanel.css";

import { RaisedButton, CircularProgress, TextField } from "material-ui";

class LeftPanel extends Component {
  constructor(props) {
    super(props);
    this.state = {
      downloadInProgress: false,
      concatInProgress: false,
      uploadInProgress: false,
      durationInProgress: false,
      downloadButtonDisabled: false,
      concatButtonDisabled: false,
      uploadButtonDisabled: false,
      durationButtonDisabled: false,
      currentGame: "Fortnite",
      vidDuration: "",
      customTitle: ""
    };
  }

  changeCurrentGame = (game, event) => {
    this.setState({
      currentGame: game
    });
  };

  changeCustomTitle = event => {
    this.setState({
      customTitle: event.target.value
    });
  };

  downloadButtonPressed = async () => {
    this.setState({
      downloadInProgress: true,
      downloadButtonDisabled: true,
      concatButtonDisabled: true,
      uploadButtonDisabled: true,
      durationButtonDisabled: true
    });

    const headers = new Headers();

    headers.append(
      "Authorization",
      `Basic ${btoa(`${this.state.username}:${this.state.password}`)}`
    );

    try {
      const numOfClips = this.state.numOfVideos || 25;

      const res = await fetch(
        `http://localhost:6969/api/downloadClips?numOfClips=${numOfClips}&gameType=${
          this.state.currentGame
        }`,
        {
          method: "POST",
          headers
        }
      );

      const json = await res.json();

      if (json.response === "success") {
        this.setState({
          downloadInProgress: false,
          downloadButtonDisabled: false,
          concatButtonDisabled: false,
          uploadButtonDisabled: false,
          durationButtonDisabled: false
        });

        this.props.setVideos(json.data.clips);
      }
    } catch (e) {
      this.setState({
        downloadInProgress: false,
        downloadButtonDisabled: false,
        concatButtonDisabled: false,
        uploadButtonDisabled: false,
        durationButtonDisabled: false
      });
    }
  };

  concatButtonPressed = async () => {
    this.setState({
      concatInProgress: true,
      concatButtonDisabled: true,
      downloadButtonDisabled: true,
      uploadButtonDisabled: true,
      durationButtonDisabled: true
    });

    const headers = new Headers();

    headers.append(
      "Authorization",
      `Basic ${btoa(`${this.state.username}:${this.state.password}`)}`
    );

    try {
      const res = await fetch(
        `http://localhost:6969/api/concatClips?gameType=${
          this.state.currentGame
        }&customTitle=${this.state.customTitle}&firstClipIndex=${
          this.props.firstClipIndex
        }`,
        {
          method: "POST",
          headers
        }
      );

      const json = await res.json();

      if (json.response === "success") {
        this.setState({
          concatInProgress: false,
          concatButtonDisabled: false,
          downloadButtonDisabled: false,
          uploadButtonDisabled: false,
          durationButtonDisabled: false
        });
      }
    } catch (e) {
      this.setState({
        concatInProgress: false,
        concatButtonDisabled: false,
        downloadButtonDisabled: false,
        uploadButtonDisabled: false,
        durationButtonDisabled: false
      });
    }
  };

  uploadButtonPressed = async () => {
    this.setState({
      uploadInProgress: true,
      downloadButtonDisabled: true,
      concatButtonDisabled: true,
      uploadButtonDisabled: true,
      durationButtonDisabled: true
    });

    try {
      const res = await fetch(
        `http://localhost:6969/api/uploadVid?gameType=${
          this.state.currentGame
        }&customTitle=${this.state.customTitle}`
      );

      const json = await res.json();

      if (json.response !== "success") {
        throw new Error(json.ErrorType);
      }

      this.setState({
        uploadInProgress: false,
        downloadButtonDisabled: false,
        concatButtonDisabled: false,
        uploadButtonDisabled: false,
        durationButtonDisabled: false
      });
    } catch (e) {
      this.setState({
        uploadInProgress: false,
        downloadButtonDisabled: false,
        concatButtonDisabled: false,
        uploadButtonDisabled: false,
        durationButtonDisabled: false
      });
    }
  };

  durationButtonPressed = async () => {
    this.setState({
      durationInProgress: true,
      downloadButtonDisabled: true,
      concatButtonDisabled: true,
      uploadButtonDisabled: true,
      durationButtonDisabled: true
    });

    try {
      const res = await fetch(
        `http://localhost:6969/api/getDuration?gameType=${
          this.state.currentGame
        }`
      );

      console.log("But did I make it here?");

      const json = await res.json();

      if (json.response !== "success") {
        throw new Error(json.errorType);
      }

      const roundedSeconds = Math.round(json.data);

      // Get number of seconds
      const numberOfSeconds = roundedSeconds % 60;

      // Get number of minutes
      const numberOfMinutes = (roundedSeconds - numberOfSeconds) / 60;

      this.setState({
        vidDuration: `${numberOfMinutes} mins ${numberOfSeconds} secs`,
        durationInProgress: false,
        downloadButtonDisabled: false,
        concatButtonDisabled: false,
        uploadButtonDisabled: false,
        durationButtonDisabled: false
      });
    } catch (e) {
      this.setState({
        durationInProgress: false,
        downloadButtonDisabled: false,
        concatButtonDisabled: false,
        uploadButtonDisabled: false,
        durationButtonDisabled: false
      });
    }
  };

  render() {
    return (
      <div>
        <h1
          style={{
            fontFamily: "'Roboto', sans-serif",
            textAlign: "center",
            opacity: "0.90"
          }}
        >
          Juke Highlights Panel
        </h1>

        <br />

        <div
          className={`panel panel-default ${
            this.state.currentGame === "Fortnite"
              ? "activeGame defaultPanel leftGame"
              : "defaultPanel leftGame"
          }`}
          onClick={this.changeCurrentGame.bind(null, "Fortnite")}
        >
          <div className="panel-body">Fortnite</div>
        </div>

        <div
          className={`panel panel-default ${
            this.state.currentGame == "VRChat"
              ? "activeGame defaultPanel rightGame"
              : "defaultPanel rightGame"
          }`}
          onClick={this.changeCurrentGame.bind(null, "VRChat")}
        >
          <div className="panel-body">VRChat</div>
        </div>

        <div
          style={{
            display: "inline-block",
            width: "49%"
          }}
        >
          <br />

          <RaisedButton
            label={
              this.state.downloadInProgress ? (
                <CircularProgress
                  size={25}
                  thickness={3}
                  style={{ marginTop: "5px" }}
                />
              ) : (
                "Download"
              )
            }
            disabled={this.state.downloadButtonDisabled}
            style={{ margin: "0 auto", display: "block" }}
            onClick={this.downloadButtonPressed}
            primary={true}
          />
        </div>

        <div
          style={{
            display: "inline-block",
            width: "49%",
            marginLeft: "10px"
          }}
        >
          <TextField
            type="number"
            onChange={this.changeNumOfVids}
            fullWidth={true}
          />
        </div>
        <div
          style={{
            width: "49%",
            display: "inline-block"
          }}
        >
          <RaisedButton
            label={
              this.state.concatInProgress ? (
                <CircularProgress
                  size={25}
                  thickness={3}
                  style={{ marginTop: "5px" }}
                />
              ) : (
                "Concat"
              )
            }
            disabled={this.state.concatButtonDisabled}
            style={{ margin: "0 auto", display: "block" }}
            onClick={this.concatButtonPressed}
            primary={true}
          />
        </div>

        <div
          style={{
            width: "49%",
            display: "inline-block",
            marginLeft: "1%"
          }}
        >
          <TextField
            value={this.state.customTitle}
            onChange={this.changeCustomTitle}
            fullWidth={true}
            label={this.customTitle}
          />
        </div>

        <div
          style={{
            width: "49%",
            display: "inline-block"
          }}
        >
          <RaisedButton
            label={
              this.state.durationInProgress ? (
                <CircularProgress
                  size={25}
                  thickness={3}
                  style={{ marginTop: "5px" }}
                />
              ) : (
                "Find Video Duration"
              )
            }
            disabled={this.state.durationButtonDisabled}
            style={{
              margin: "0 auto",
              display: "block",
              marginTop: "5px"
            }}
            onClick={this.durationButtonPressed}
            primary={true}
          />
        </div>

        <div
          style={{
            width: "49%",
            display: "inline-block",
            textAlign: "center",
            marginLeft: "1%"
          }}
        >
          <TextField
            value={this.state.vidDuration}
            fullWidth={true}
            inputStyle={{
              textAlign: "center",
              pointerEvents: "none",
              cursor: "default"
            }}
          />
        </div>

        <div style={{ width: "49%", display: "inline-block" }}>
          <RaisedButton
            label={
              this.state.uploadInProgress ? (
                <CircularProgress
                  size={25}
                  thickness={3}
                  style={{ marginTop: "5px" }}
                />
              ) : (
                "Upload"
              )
            }
            disabled={this.state.uploadButtonDisabled}
            style={{
              margin: "0 auto",
              display: "block",
              marginTop: "10px"
            }}
            onClick={this.uploadButtonPressed}
            primary={true}
          />
        </div>

        <div
          style={{ width: "49%", display: "inline-block", marginLeft: "1%" }}
        >
          <TextField
            onChange={this.changeCustomTitle}
            fullWidth={true}
            label={this.customTitle}
            value={this.state.customTitle}
          />
        </div>
      </div>
    );
  }
}

export default LeftPanel;
