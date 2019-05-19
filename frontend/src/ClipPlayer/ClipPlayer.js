import React, { Component } from "react";
import PropTypes from "prop-types";
import styles from "./ClipPlayer.css";

import {
  CircularProgress,
  RaisedButton,
  Dialog,
  FlatButton,
  Checkbox
} from "material-ui";

class ClipPlayer extends Component {
  constructor(props) {
    super(props);

    this.state = {
      open: false
    };
  }

  deleteVideo = async () => {
    const res = await fetch(
      `http://localhost:6969/api/deleteVideo?tracking_id=${
        this.props.videos[this.props.currentIndex].tracking_id
      }`,
      {
        method: "GET"
      }
    );

    if (this.props.videos.length === 1) {
    } else {
      const videos = this.props.videos.slice();
      videos.splice(this.state.currentIndex, 1);

      this.props.setVideos(videos);
    }

    this.handleClose();
  };

  handleOpen = () => {
    this.setState({
      open: true
    });
  };

  handleClose = () => {
    this.setState({
      open: false
    });
  };

  render() {
    const actions = [
      <FlatButton label="Cancel" primary={true} onClick={this.handleClose} />,
      <FlatButton
        label="Submit"
        primary={true}
        keyboardFocused={true}
        onClick={this.deleteVideo}
      />
    ];

    const styles = {
      defaultPanel: {
        width: "49%",
        display: "inline-block",
        cursor: "pointer",
        textAlign: "center"
      },
      leftGame: {
        borderTopRightRadius: "0px",
        borderBottomRightRadius: "0px"
      },
      rightGame: {
        borderTopLeftRadius: "0px",
        borderBottomLeftRadius: "0px",
        cursor: "pointer"
      },
      activeGame: {
        borderColor: "rgb(0, 180, 206)",
        color: "rgb(0, 180, 206)"
      }
    };

    console.log(this.props.videos);

    return (
      <div>
        {this.props.videos === null ? (
          <CircularProgress
            size={60}
            thickness={5}
            style={{
              display: "block",
              margin: "0 auto",
              marginTop: "100px"
            }}
          />
        ) : this.props.videos.length === 0 ? null : (
          <div>
            <h3>{this.props.videos[this.props.currentIndex].title}</h3>

            <hr />
            <h4>
              {this.props.videos[this.props.currentIndex].tracking_id}.mp4
            </h4>

            <Checkbox
              checked={this.props.currentIndex === this.props.firstClipIndex}
              onCheck={this.props.changeFirstClipIndex}
            />

            <video
              src={`http://localhost:6969/api/clip_files/${
                this.props.videos[this.props.currentIndex].tracking_id
              }.mp4`}
              style={{ width: "100%", height: "auto" }}
              controls
            />

            <div
              style={{
                width: "33%",
                display: "inline-block",
                textAlign: "center",
                fontSize: "30px"
              }}
            >
              {this.props.currentIndex > 0 ? (
                <i
                  className="fa fa-arrow-left arrows"
                  aria-hidden="true"
                  style={{
                    color: "rgb(31, 188, 211)",
                    cursor: "pointer"
                  }}
                  onClick={this.props.goToPreviousVideo}
                />
              ) : null}
            </div>

            <div style={{ width: "33%", display: "inline-block" }}>
              <RaisedButton
                label="Delete"
                secondary={true}
                style={{ margin: "0 auto", display: "block" }}
                onClick={this.handleOpen}
              />

              <Dialog
                title={`Delete ${
                  this.props.videos[this.props.currentIndex].tracking_id
                }.mp4`}
                actions={actions}
                modal={false}
                open={this.state.open}
                onRequestClose={this.handleClose}
              >
                Are you sure you want to delete{" "}
                {this.props.videos[this.props.currentIndex].tracking_id}.mp4 ?
              </Dialog>
            </div>

            <div
              style={{
                width: "33%",
                display: "inline-block",
                textAlign: "center",
                fontSize: "30px"
              }}
            >
              {this.props.currentIndex < this.props.videos.length - 1 ? (
                <i
                  className="fa fa-arrow-right arrows"
                  aria-hidden="true"
                  style={{
                    color: "rgb(31, 188, 211)",
                    cursor: "pointer"
                  }}
                  onClick={this.props.goToNextVideo}
                />
              ) : null}
            </div>
          </div>
        )}
      </div>
    );
  }
}

export default ClipPlayer;
