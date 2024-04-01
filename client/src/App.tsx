import { createResource, createSignal, type Component } from "solid-js";

import { createReconnectingWS } from "@solid-primitives/websocket";

import styles from "./App.module.css";

const fetchFormData = async () => {
  const res = await fetch("http://localhost:8080/form/1");
  return res.json();
};

const App: Component = () => {
  const user = Math.random().toString().substr(2, 6);
  const ws = createReconnectingWS("ws://localhost:8080/ws");
  const [data, { mutate }] = createResource(fetchFormData);

  ws.addEventListener("message", (ev) => {
    console.log(data());
    console.log(ev.data);
    mutate(JSON.parse(ev.data));
  });

  return (
    <div class={styles.App}>
      <label for="testInput">Input:</label>
      <input
        id="testInput"
        type="text"
        value={data()?.value || ""}
        onInput={(ev) =>
          ws.send(JSON.stringify({ user, value: ev.currentTarget.value }))
        }
      />
    </div>
  );
};

export default App;
