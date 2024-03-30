import { createSignal, type Component } from "solid-js";
import { createStore } from "solid-js/store";

import { createReconnectingWS, WSMessage } from "@solid-primitives/websocket";

import logo from "./logo.svg";
import styles from "./App.module.css";

const App: Component = () => {
  const user = Math.random().toString().substr(2, 6);
  const ws = createReconnectingWS("ws://localhost:8080/ws");
  const [data, setData] = createSignal("");

  ws.addEventListener("message", (ev) => {
    console.log(ev.data);
    setData(JSON.parse(ev.data).value);
  });

  return (
    <div class={styles.App}>
      <label for="testInput">Input:</label>
      <input
        id="testInput"
        type="text"
        value={data()}
        onInput={(ev) =>
          ws.send(JSON.stringify({ user, value: ev.currentTarget.value }))
        }
      />
    </div>
  );
};

export default App;
