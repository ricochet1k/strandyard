import { render } from "solid-js/web"
import App from "./App"
import KanbanTransitionTest from "./components/KanbanTransitionTest"
import TransitionTest from "./components/TransitionTest"
import "./styles.css"

const params = new URLSearchParams(window.location.search)
const transitionTestMode = params.get("e2eTransition") === "1"
const kanbanTransitionTestMode = params.get("e2eKanban") === "1"

render(() => {
  if (transitionTestMode) return <TransitionTest />
  if (kanbanTransitionTestMode) return <KanbanTransitionTest />
  return <App />
}, document.getElementById("root")!)
