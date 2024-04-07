import { type Component } from "solid-js";
import { Router, Route } from "@solidjs/router";
import FormPage from "./pages/FormPage/FormPage";

const App: Component = () => {
  return (
    <Router>
      <Route path="form/:id" component={FormPage} />
    </Router>
  );
};

export default App;
