import React, { Component } from "react";
import PropTypes from "prop-types";
import styles from "./VidPlayer.scss";

class VidPlayer extends Component {
  constructor(props) {
    super(props);
  }

  render() {
    return (
      <div class="col-md-12">
        <h2 style={{ textAlign: "center" }}>Juke Highlights Vid</h2>

        <hr />

        {this.props.vidExists === false ? null : (
          <video
            src={`http://localhost:6969/api/juke_highlights_vid.mp4`}
            style={{ width: "100%", height: "auto" }}
            controls
          />
        )}
      </div>
    );
  }
}

export default VidPlayer;
