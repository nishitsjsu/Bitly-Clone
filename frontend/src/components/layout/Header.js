import React, { Component } from "react";
import { Link } from "react-router-dom";

export class Header extends Component {
  render() {
    return (
      <div>
        <nav
          className="navbar navbar-expand-sm navbar-light bg-light"
          style={{ height: "70px" }}
        >
          <button
            className="navbar-toggler"
            type="button"
            data-toggle="collapse"
            data-target="#navbarTogglerDemo01"
            aria-controls="navbarTogglerDemo01"
            aria-expanded="false"
            aria-label="Toggle navigation"
          >
            <span className="navbar-toggler-icon"></span>
          </button>
          <div className="collapse navbar-collapse" id="navbarTogglerDemo01">
            <a className="navbar-brand" href="/">
              <img
                src="https://docrdsfx76ssb.cloudfront.net/static/1606773372/pages/wp-content/uploads/2019/02/bitly.png"
                width="100"
                height="auto"
                style={{ marginLeft: "8rem" }}
              ></img>
            </a>
            <ul className="navbar-nav mr-auto mt-2 mt-lg-0">
              <li className="nav-item active">
                <a
                  className="navbar-brand"
                  href="#"
                  style={{ marginLeft: "8rem" }}
                >
                  Why Bitly? <span className="sr-only">(current)</span>
                </a>
              </li>
              <li className="nav-item active">
                <a
                  className="navbar-brand"
                  href="#"
                  style={{ marginLeft: "1rem" }}
                >
                  Solutions <span className="sr-only">(current)</span>
                </a>
              </li>
              <li className="nav-item active">
                <a
                  className="navbar-brand"
                  href="#"
                  style={{ marginLeft: "1rem" }}
                >
                  Features <span className="sr-only">(current)</span>
                </a>
              </li>
              <li className="nav-item active">
                <a
                  className="navbar-brand"
                  href="#"
                  style={{ marginLeft: "1rem" }}
                >
                  Pricing <span className="sr-only">(current)</span>
                </a>
              </li>
              <li className="nav-item active">
                <a
                  className="navbar-brand"
                  href="#"
                  style={{ marginLeft: "1rem" }}
                >
                  Resources <span className="sr-only">(current)</span>
                </a>
              </li>
              <li className="nav-item active">
                <a
                  className="navbar-brand"
                  href="#"
                  style={{ marginLeft: "10rem" }}
                >
                  Login <span className="sr-only">(current)</span>
                </a>
              </li>
              <li className="nav-item active">
                <a
                  className="navbar-brand"
                  href="#"
                  style={{ marginLeft: "1rem", color: "#2a5bd7" }}
                >
                  Sign Up <span className="sr-only">(current)</span>
                </a>
              </li>
              <li className="nav-item active">
                <Link to="/trends">
                  <button
                    className="btn "
                    style={{
                      marginTop: "0.3rem",
                      backgroundColor: "#2a5bd7",
                      color: "white",
                      height: "40px",
                      width: "120px",
                      fontSize: "15px",
                    }}
                  >
                    Get Trends
                  </button>
                </Link>
              </li>
            </ul>
          </div>
        </nav>
      </div>
    );
  }
}

export default Header;
