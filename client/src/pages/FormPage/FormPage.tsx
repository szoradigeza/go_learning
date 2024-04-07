import { createSignal, type Component } from "solid-js";

import { createReconnectingWS } from "@solid-primitives/websocket";
import { useParams } from "@solidjs/router";

const FormPage: Component = () => {
  const user = Math.random().toString().substr(2, 6);
  const params = useParams();
  const ws = createReconnectingWS("ws://localhost:8080/ws");
  const [data, setData] = createSignal({ value: "", id: Number(params.id) });

  ws.send(JSON.stringify({ action: "getForm", id: Number(params.id) }));

  ws.addEventListener("message", (ev) => {
    console.log(ev.data);
    setData(JSON.parse(ev.data));
  });

  return (
    <>
      <h1>{data().id}</h1>
      <label for="testInput">Input:</label>
      <input
        id="testInput"
        type="text"
        value={data()?.value || ""}
        onInput={(ev) => {
          console.log(ev.currentTarget.value);
          try {
            ws.send(
              JSON.stringify({
                ...data(),
                user,
                value: ev.currentTarget.value,
              }),
            );
          } catch (err) {
            console.log(err);
          }
        }}
      />
    </>
  );
};

export default FormPage;
