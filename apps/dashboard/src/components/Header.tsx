import { For, Show } from "solid-js"
import "./Header.css"

type ProjectInfo = {
  name: string
  storage_root: string
  tasks_root: string
  roles_root: string
  templates_root: string
  git_root: string
  storage: string
}

type HeaderProps = {
  projects: ProjectInfo[]
  currentProject: string
  connected: boolean
  onProjectChange: (project: string) => void
}

export default function Header(props: HeaderProps) {
  return (
    <header class="top-bar">
      <p class="kicker">StrandYard</p>
      <div style={{ display: "flex", gap: "1rem", "align-items": "center" }}>
        <div class="status-panel">
          <div class={`status-dot ${props.connected ? "online" : "offline"}`} />
          <div>
            <p class="status-detail">{props.connected ? "Connected" : "Reconnecting"}</p>
          </div>
        </div>
      </div>
    </header>
  )
}
