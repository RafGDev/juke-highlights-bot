import React, { Component } from "react";
import logo from "./logo.svg";
import "./App.css";
import {
  TextField,
  RaisedButton,
  CircularProgress,
  Dialog,
  FlatButton
} from "material-ui";
import { MuiThemeProvider } from "material-ui/styles";

import LeftPanel from "./LeftPanel/LeftPanel";
import ClipPlayer from "./ClipPlayer/ClipPlayer";
import VidPlayer from "./VidPlayer/VidPlayer";

class App extends Component {
  constructor(props) {
    super(props);

    this.state = {
      firstClipIndex: 0,
      downloadInProgress: false,
      concatInProgress: false,
      uploadInProgress: false,
      durationInProgress: false,
      downloadButtonDisabled: false,
      concatButtonDisabled: false,
      uploadButtonDisabled: false,
      durationButtonDisabled: false,
      username: "",
      password: "",
      videos: null,
      currentIndex: 0,
      open: false,
      numOfVideos: "",
      currentGame: "Fortnite",
      vidDuration: "",
      customTitle: "",
      vidExists: false,
      customName: ""
    };
  }

  checkVidExists = async () => {
    try {
      const res = await fetch("http://localhost:6969/api/checkVidExists");

      const json = await res.json();

      if (json.response === "success") {
        if (json.data === true) {
          this.setState({
            vidExists: true
          });
        } else {
          throw new Error("fileNotExists");
        }
      } else {
        this.setState({
          vidExists: false
        });
      }
    } catch (e) {
      this.setState({
        vidExists: false
      });
    }
  };

  getClips = async () => {
    try {
      const res = await fetch("http://localhost:6969/api/getClips");

      const json = await res.json();

      if (json.response === "success") {
        this.setState({
          videos: json.data.clips
        });

        return;
      }

      throw new Error(json.errorType);
    } catch (e) {
      this.setState({
        videos: []
      });
    }
  };

  componentDidMount = async () => {
    // Fetch list of mp4's and list them in the video player

    await Promise.all([this.checkVidExists(), this.getClips()]);
  };

  setVideos = videos => {
    this.setState({
      videos
    });
  };

  changeFirstClipIndex = (event, value) => {
    if (value === false) {
      return;
    }

    this.setState({
      firstClipIndex: this.state.currentIndex
    });
  };

  goToNextVideo = () => {
    this.setState({
      currentIndex: this.state.currentIndex + 1
    });
  };

  goToPreviousVideo = () => {
    this.setState({
      currentIndex: this.state.currentIndex - 1
    });
  };

  render() {
    return (
      <MuiThemeProvider>
        <div className="container">
          <div className="col-md-12">
            <div className="col-md-6">
              <LeftPanel
                setVideos={this.setVideos}
                firstClipIndex={this.state.firstClipIndex}
              />
            </div>
            <div className="col-md-6">
              <ClipPlayer
                videos={this.state.videos}
                setVideos={this.setVideos}
                currentIndex={this.state.currentIndex}
                firstClipIndex={this.state.firstClipIndex}
                changeFirstClipIndex={this.changeFirstClipIndex}
                goToNextVideo={this.goToNextVideo}
                goToPreviousVideo={this.goToPreviousVideo}
              />
            </div>
          </div>
          <div class="col-md-12">
            <VidPlayer vidExists={this.state.vidExists} />
          </div>
        </div>
      </MuiThemeProvider>
    );
  }
}

export default App;
