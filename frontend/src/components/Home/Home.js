import React, { Component } from "react";
import { Link } from "react-router-dom";
import "./Home.css";
import home_pic from "./Home.png";
import axios from "axios";
import Swal from "sweetalert2";

export class Home extends Component {
  constructor(props) {
    //Call the constrictor of Super class i.e The Component
    super(props);
    //maintain the state required for this component
    this.state = {
      uri: "",
      shortLink: "",
      isShortened: false,
      authFlag: false,
    };
  }

  componentDidMount() {
    axios
      .get(process.env.REACT_APP_BACKEND_URI + "/lrs/all/trends", {})
      .then((response) => {
        console.log("Received response");
        //update the state with the response data
        console.log(response.data);
      });
  }

  uriChangeHandler = (e) => {
    this.setState({
      uri: e.target.value,
    });
  };

  copyText = (e) => {
    e.preventDefault();
    var copyText = document.getElementById("uri");
    // console.log(copyText);
    copyText.select();
    copyText.setSelectionRange(0, 99999);
    document.execCommand("copy");

    Swal.fire("The link has been copied!", "Try it out!", "success").then(
      () => {
        window.location.reload();
      }
    );
  };

  //submit Shorten button to send a request to the go backend
  submitShorten = (e) => {
    //make a post request with the uri
    axios
      .post(
        process.env.REACT_APP_BACKEND_URI + "/cps/createlink",
        { Uri: this.state.uri },
        {
          headers: {
            "content-type": "application/x-www-form-urlencoded",
          },
        }
      )
      .then((response) => {
        console.log("Status Code : ", response);
        if (response.status === 200) {
          this.setState({
            isShortened: true,
            uri: response.data.ShortLink,
          });
        }
      })
      .catch((e) => {
        console.log("in catch " + e);
      });
  };

  render() {
    return (
      <div>
        <div className="content-wrapper">
          <div className="row">
            <div className="col-1"></div>
            <div
              className="col-6 content-style"
              style={{ paddingLeft: "2rem" }}
            >
              <h1 className="header-xl"> Short links, big results</h1>
              <h2 className="suubhead-xl ">
                A URL shortener built with powerful tools to help you grow and
                protect your brand.
              </h2>
              <Link to="/trends">
                <button className="btn trending-button">Trending links</button>
              </Link>
            </div>
            <div className="col-5">
              <img src={home_pic} className="pic-style"></img>
            </div>
          </div>
        </div>
        <div>
          <div className="footer-wrapper">
            <div className="footer">
              <form>
                <div class="form-row">
                  <div className="col-1"></div>
                  <div class="col-7" style={{ marginLeft: "2rem" }}>
                    <input
                      type="text"
                      name="uri"
                      id="uri"
                      class={
                        this.state.isShortened
                          ? "form-control footer-input short-link"
                          : "form-control footer-input"
                      }
                      placeholder="Shorten your link"
                      value={this.state.uri}
                      onChange={this.uriChangeHandler}
                    ></input>
                  </div>
                  <div class="col-2">
                    {this.state.isShortened ? (
                      <button
                        type="button"
                        name="copy"
                        class="form-control btn footer-input footer-button"
                        onClick={this.copyText}
                      >
                        Copy Link
                      </button>
                    ) : (
                      <button
                        type="button"
                        name="shorten"
                        class="form-control btn footer-input footer-button"
                        onClick={this.submitShorten}
                      >
                        Shorten
                      </button>
                    )}
                  </div>
                </div>
              </form>
            </div>
          </div>
        </div>
      </div>
    );
  }
}

export default Home;
