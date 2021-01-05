import React from "react";
import { HashRouter as Router, Route, BrowserRouter } from "react-router-dom";
import "./App.css";
import Header from "./components/layout/Header";
import Home from "./components/Home/Home";
import Trends from "./components/Trends/Trends";

function App() {
  return (
    <Router>
      <BrowserRouter>
        <Header />
        <div className="App">
          <Route exact path="/" component={Home} />
          <Route exact path="/trends" component={Trends} />
        </div>
      </BrowserRouter>
    </Router>
  );
}

export default App;
