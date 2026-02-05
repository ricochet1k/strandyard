import { For } from "solid-js"
import type { RoleItem } from "../App"
import "./TaskTable.css"

type RolesTableProps = {
  roles: RoleItem[]
  activePath: string
  onSelect: (role: RoleItem) => void
}

export default function RolesTable(props: RolesTableProps) {
  return (
    <table class="task-table">
      <thead>
        <tr>
          <th class="col-name">Name</th>
          <th class="col-description">Description</th>
        </tr>
      </thead>
      <tbody>
        <For each={props.roles}>
          {(role) => (
            <tr
              class={`task-row ${props.activePath === role.path ? "active" : ""}`}
              onClick={() => props.onSelect(role)}
            >
              <td class="col-name">
                <button class="task-link">{role.name}</button>
              </td>
              <td class="col-description">{role.description}</td>
            </tr>
          )}
        </For>
      </tbody>
    </table>
  )
}
