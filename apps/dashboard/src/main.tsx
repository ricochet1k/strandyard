import { render } from "solid-js/web"
import App from "./App"
import TransitionTest from "./components/TransitionTest"
import "./styles.css"

const params = new URLSearchParams(window.location.search)
const transitionTestMode = params.get("e2eTransition") === "1"

render(() => (transitionTestMode ? <TransitionTest /> : <App />), document.getElementById("root")!)
