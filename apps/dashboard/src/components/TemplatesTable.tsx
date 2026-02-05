import { For } from "solid-js"
import type { TemplateItem } from "../App"
import "./TaskTable.css"

type TemplatesTableProps = {
  templates: TemplateItem[]
  activePath: string
  onSelect: (template: TemplateItem) => void
}

export default function TemplatesTable(props: TemplatesTableProps) {
  return (
    <table class="task-table">
      <thead>
        <tr>
          <th class="col-name">Name</th>
          <th class="col-role">Role</th>
          <th class="col-priority">Priority</th>
          <th class="col-description">Description</th>
        </tr>
      </thead>
      <tbody>
        <For each={props.templates}>
          {(template) => (
            <tr
              class={`task-row ${props.activePath === template.path ? "active" : ""}`}
              onClick={() => props.onSelect(template)}
            >
              <td class="col-name">
                <button class="task-link">{template.name}</button>
              </td>
              <td class="col-role">
                <span class="pill role">{template.role}</span>
              </td>
              <td class="col-priority">
                <span class={`pill priority-${template.priority}`}>{template.priority}</span>
              </td>
              <td class="col-description">{template.description}</td>
            </tr>
          )}
        </For>
      </tbody>
    </table>
  )
}
