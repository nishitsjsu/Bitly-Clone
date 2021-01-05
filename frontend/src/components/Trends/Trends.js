import React, { Component } from "react";
import "./Trends.css";
import axios from "axios";
import Pagination from "../Pagination/Pagination";

export class Trends extends Component {
  constructor(props) {
    super(props);
    this.state = {
      uri: "",
      hits: "",
      isShortened: false,
      urls: [],
      currenPage: 1,
      sectionsPerPage: 4,
    };
  }

  componentDidMount() {
    axios
      .get(process.env.REACT_APP_BACKEND_URI + "/lrs/all/trends", {})
      .then((response) => {
        console.log("Received response ", response.data);
        this.setState({
          urls: this.state.urls.concat(response.data),
        });
      });
  }

  uriChangeHandler = (e) => {
    this.setState({
      uri: e.target.value,
    });
  };

  copyText = (e) => {
    window.location.reload();
  };

  submitShorten = (e) => {
    let shortUriArray = this.state.uri.split("/");
    let shortUri = shortUriArray[shortUriArray.length - 1];
    axios
      .get(process.env.REACT_APP_BACKEND_URI + "/lrs/trends/" + shortUri)
      .then((response) => {
        console.log("Status Code : ", response);
        if (response.status === 200) {
          this.setState({
            isShortened: true,
            uri:
              "Requested URI has been hit " + response.data.Count + " times!",
          });
        }
      })
      .catch((e) => {
        console.log("in catch " + e);
      });
  };

  render() {
    const indexOfLastSection =
      this.state.currenPage * this.state.sectionsPerPage;
    const indexOfFirstSection = indexOfLastSection - this.state.sectionsPerPage;
    const currenSections = this.state.urls
      .sort((a, b) => parseFloat(b.Count) - parseFloat(a.Count))
      .slice(indexOfFirstSection, indexOfLastSection);

    const paginate = (pageNumber) => {
      this.setState({
        currenPage: pageNumber,
      });
    };

    let details = currenSections.map((url) => {
      return (
        <tr key="index" className="table-row-content">
          <td className="table-row-content">{url.Uri}</td>
          <td className="table-row-content">{url.ShortLink}</td>
          <td className="table-row-content">{url.Count}</td>
        </tr>
      );
    });
    return (
      <div>
        <div className="content-wrapper">
          <div className="row">
            <div className="col-4"></div>
            <div className="col-4 content-style">
              <h1 className="header-xl"> Trending Links!</h1>
            </div>
            <div className="col-4"></div>
          </div>
          <div className="row">
            <div className="col-2"></div>
            <div className="col-8">
              <table className="table table-bordered table-hover">
                <thead>
                  <tr>
                    <th className="table-head-content">Original URI</th>
                    <th className="table-head-content">Shortened URI</th>
                    <th className="table-head-content">Hits</th>
                  </tr>
                </thead>
                <tbody>{details}</tbody>
              </table>
              <Pagination
                postsPerPage={this.state.sectionsPerPage}
                totalPosts={this.state.urls.length}
                paginate={paginate}
              />
            </div>
            <div className="col-2"></div>
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
                      placeholder="Enter the URI to get Hit count"
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
                        Okay!
                      </button>
                    ) : (
                      <button
                        type="button"
                        name="shorten"
                        class="form-control btn footer-input footer-button"
                        onClick={this.submitShorten}
                      >
                        Get Count
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

export default Trends;
